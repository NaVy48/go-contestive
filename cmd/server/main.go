package main

import (
	"contestive/config"
	"contestive/repository/postgresql"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

func main() {
	run()
}

func run() {
	defaultOp := func() {
		runServer()
	}

	if len(os.Args) == 1 {
		defaultOp()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "serve":
		runServer()
	case "dropdb":
		dropDB()
	case "resetdb":
		ResetDB()
	case "genkey":
		genKeyHex()
	default:
		defaultOp()
	}
}

func dropDB() {
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}
	config := cfg.Database.PostgreSQL
	repository, err := postgresql.Connect(config.Address, config.Username, config.Password, config.Database)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	repository.Drop()
}

func ResetDB() {
	success := false
	fmt.Println("Reseting database...")
	defer func() {
		if success {
			fmt.Println("Database reset successfully")
		} else {
			fmt.Println("Database reset failed")
		}
	}()
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}
	config := cfg.Database.PostgreSQL
	repository, err := postgresql.Connect(config.Address, config.Username, config.Password, config.Database)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = repository.Reset()
	if err != nil {
		logger.Fatalf(err.Error())
	}
	success = true
}

func MigrateDB() {
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}
	config := cfg.Database.PostgreSQL
	repository, err := postgresql.Connect(config.Address, config.Username, config.Password, config.Database)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	repository.Drop()
}

func genKeyHex() {
	byteLen := flag.Int("l", 32, "Key length in bytes")
	flag.Parse()
	bytes := make([]byte, *byteLen)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("This is random key:\n%s\n", hex.EncodeToString(bytes))
	fmt.Printf("This is random key:\n%s\n", base64.StdEncoding.EncodeToString(bytes))
}
