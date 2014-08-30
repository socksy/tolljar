package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sfreiberg/gotwilio"
	"log"
	"os"
	"time"
)

func connect() *sql.DB {
	password := os.Getenv("POSTGRESPASSWORD")
	if password == "" {
		fmt.Println("Failed to get $POSTGRESPASSWORD. Have you sourced authdetails.sh?")
	}

	db, err := sql.Open("postgres", "user=postgres password='"+password+
		"' dbname=tolljardates sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func setup_twilio() (*gotwilio.Twilio, string) {
	accountSid := os.Getenv("ACCOUNTSID")
	authToken := os.Getenv("AUTHTOKEN")
	if accountSid == "" || authToken == "" {
		fmt.Println("couldn't find $ACCOUNTSID or $AUTHTOKEN. Have you sourced authdetails.sh?")
	}

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	return twilio, "+441315104795"
}

func process_today() {
	db := connect()
	rows, err := db.Query("SELECT * FROM dates WHERE due = current_date;")
	if err != nil {
		log.Fatal(err)
	}
	twilio, from_number := setup_twilio()

	for rows.Next() {
		var prediction *Prediction = new(Prediction)
		err = rows.Scan(&prediction.due, &prediction.when,
			&prediction.to, &prediction.from, &prediction.what)

		if err != nil {
			log.Fatal(err)
		}

		prediction.process(twilio, from_number)
	}
}

type Prediction struct {
	due  time.Time
	when time.Time
	to   string
	from string
	what string
}

func (prediction *Prediction) process(twilio *gotwilio.Twilio, from_number string) {
	fmt.Println(prediction.due.Date())
	const layout = "Mon 2 Jan, 2006"
	message := "On " + prediction.when.Format(layout) + ", " +
		prediction.from + " predicted that: " +
		prediction.what + "\nWere they right? Reply w/ yes/no! \n- Courtesy of Tolljar"
	twilio.SendSMS(from_number, prediction.to, message, "", "")
}

func main() {
	process_today()
}
