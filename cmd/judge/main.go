package main

import (
	"contestive/judge"
	"contestive/judgeconnection"
	"flag"
	"log"
	"time"
)

var judgeID = flag.Int("id", 0, "(Required, integer) Unique judge id")
var secret = flag.String("secret", "", "(Required, string) should be set in contest sytem too. Is used to authentiate judge")
var address = flag.String("add", "127.0.0.1:4545", "(string) network address where contest system is listening")
var problemDir = flag.String("dir", "/var/local/judge0", "Directory where judge keeps problem packages")

var retryTime = flag.Int("retry", 2, "Timeout before redial")

func main() {
	flag.Parse()

	if *judgeID < 0 || *secret == "" {
		flag.Usage()
		return
	}

	for {
		log.Printf("Connecting too %s", *address)
		conn, err := judgeconnection.NewClient(*judgeID, *secret, *address)
		if err != nil {
			log.Printf("%v\n", err)
			log.Printf("Connection failed. Restarting in %ds\n", *retryTime)
			time.Sleep(time.Duration(*retryTime) * time.Second)
			continue
		}

		err = judge.RunJudge(conn, *problemDir, *judgeID)
		if err != nil {
			conn.Close()
			log.Printf("Judge failed with error %v", err)
			time.Sleep(time.Duration(*retryTime) * time.Second)
		}
		conn.Close()
	}
}
