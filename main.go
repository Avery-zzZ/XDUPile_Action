package main

import (
	"log"
	"os"
	"time"
)

const (
	MAX_LOGIN_FAIL_TIMES = 3
	MAX_POST_FAIL_TIMES  = 3
)

func main() {
	logged_in := false

	try_login_time := 0
	for !logged_in {

		var err error
		logged_in, err = login()

		if logged_in {
			break
		}

		if err != nil {
			log.Println("Login fail: " + err.Error())
			if err.Error() != "connection error" {
				log.Println("Unable to login. Exiting")
				os.Exit(1)
			}
		}

		try_login_time++

		//if cannot login after 3 trys
		if try_login_time >= MAX_LOGIN_FAIL_TIMES {
			log.Println("Unable to login")
			// cannot_login_handle()
			os.Exit(2)
		}

		time.Sleep(5 * time.Minute)
	}
	log.Println("Login success")


	log.Println("Posting...")
	var state int
	var msg string
	_, state, msg = post()
	/*
		state:
			0-OK
			1-Already done
			2-Not logged in
			128-Unknow response
			256-Connection error
	*/
	switch state {
	case TASK_OK:
		log.Println("Done")
		push_email("XDUPILE HealthCard","Success")
	case TASK_ALREADT_DONE:
		log.Println("Already done")
	case TASK_NOT_LOGGED_IN:
		log.Println("Fail: Not logged in. Exiting")
		push_email("XDUPILE HealthCard","[Error]Unable to login. Exited")
		os.Exit(3)
	case TASK_UNKNOWN_RESPONSE:
		log.Println("Unknown response: " + msg + ". Exiting")
		push_email("XDUPILE HealthCard","[Error]Meet unknow response. Exited")
		os.Exit(3)
	case TASK_CONNECTION_ERROR:
		try_time := 1
		log.Printf("Fail: Connection error. Retrying(%d)...", try_time)
		for {
			try_time++
			_, state, _ = post()
			if state == TASK_OK {
				break
			}
			if try_time >= MAX_POST_FAIL_TIMES {
				log.Println("All trys failed. Exiting")
				push_email("XDUPILE HealthCard","[Error]connection error. Exited")
				os.Exit(4)
			}
			log.Printf("Fail. Retrying(%d)...", try_time)
		}
	default:
		log.Println("[ERROR]Unknow response status. Check your code. Exiting")
		push_email("XDUPILE HealthCard","[Error]System error: Unknow Status. Exited")
		os.Exit(5)
	}
}
