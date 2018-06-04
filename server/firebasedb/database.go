package firebasedb

import (
	"context"
	"flag"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// Score datatype for interacting with the Leaderboard DB
type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// DBClient contains connection to Firebase DB
type DBClient struct {
	DBCon *db.Client
}

// DBC is a connection handle to the firebase database
var DBC DBClient

// ConnectDB connects the DB handle to firebase db.
func (DBC *DBClient) ConnectDB() {
	// Initialize default DB App
	if flag.Lookup("test.v") == nil {
		opt := option.WithCredentialsFile("firebasedb/BumperDB-3b7d790985b1.json")
	} else {
		opt := option.WithCredentialsFile("BumperDB-3b7d790985b1.json")
	}

	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://bumperdb-d7f48.firebaseio.com",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	// Connect access to the DB Client
	DBC.DBCon, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("error getting DB client: %v", err)
	}
}