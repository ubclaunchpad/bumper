package firebasedb

import (
	"context"
	"log"

	"github.com/ubclaunchpad/bumper/server/models"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// LeaderboardEntry datatype for interacting with the Leaderboard DB
type LeaderboardEntry struct {
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
func (DBC *DBClient) ConnectDB(credentialsPath string) {
	// Initialize default DB App
	opt := option.WithCredentialsFile(credentialsPath)

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

// UpdatePlayerScore updates the score for the given player in the database
// Use as a goroutine to make it non-blocking. (Not fully thread safe yet)
func UpdatePlayerScore(p *models.Player) {

	scoreData := LeaderboardEntry{
		Name:  p.Name,
		Score: p.Points,
	}

	// Send score data to DB
	err := DBC.DBCon.NewRef("leaderboard/Testers/"+p.ID).Set(context.Background(), scoreData)
	if err != nil {
		log.Fatalf("Couldn't set data: %v", err)
	}
}

// FetchPlayerScore retreives the score for the given player
func FetchPlayerScore(p *models.Player) LeaderboardEntry {

	var scoreData LeaderboardEntry
	err := DBC.DBCon.NewRef("leaderboard/Testers/"+p.ID).Get(context.Background(), &scoreData)
	if err != nil {
		log.Fatalf("Couldn't set data: %v", err)
	}
	return scoreData
}

// FetchTop5Players prints the top 5 player - for debugging, no test written.
func FetchTop5Players() {

	query := DBC.DBCon.NewRef("leaderboard/Testers").OrderByChild("Score").LimitToFirst(5)
	result, err := query.GetOrdered(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Results will be logged in the increasing order of balance.
	for _, r := range result {
		var playerScore LeaderboardEntry
		err = r.Unmarshal(&playerScore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s => %v\n", r.Key(), playerScore)
	}
}
