package main

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

//only support qq email
func push_email(subject, body string) {
	address := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	m := gomail.NewMessage()
	m.SetHeader("From", address)
	m.SetHeader("To", address) //主送
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer("smtp.qq.com", 587, address, password)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Oops, some errors occured when sending email: "+err.Error())
	}
}