package database

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

// DBC is a connection handle to the firebase database
var DBC *db.Client

// ConnectDB connects the DB handle to firebase db.
func ConnectDB(credentialsPath string) {
	// Initialize default DB App
	opt := option.WithCredentialsFile(credentialsPath)

	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://bumperdb-d7f48.firebaseio.com",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
	}

	// Connect access to the DB Client
	DBC, err = app.Database(ctx)
	if err != nil {
		log.Printf("error getting DB client: %v", err)
	}
}

// UpdatePlayerScore updates the score for the given player in the database
// Use as a goroutine to make it non-blocking. (Not fully thread safe yet)
func UpdatePlayerScore(p *models.Player) {
	if DBC == nil {
		return
	}

	scoreData := LeaderboardEntry{
		Name:  p.Name,
		Score: p.Points,
	}

	err := DBC.NewRef("leaderboard/"+p.ID).Set(context.Background(), scoreData)
	if err != nil {
		log.Printf("Couldn't set data: %v", err)
	}
}

// FetchPlayerScore retreives the score for the given player
func FetchPlayerScore(p *models.Player) *LeaderboardEntry {
	if DBC == nil {
		return &LeaderboardEntry{}
	}

	var scoreData LeaderboardEntry
	err := DBC.NewRef("leaderboard/"+p.ID).Get(context.Background(), &scoreData)
	if err != nil {
		log.Printf("Couldn't set data: %v", err)
	}
	return &scoreData
}

// FetchTop5Players prints the top 5 player - for debugging, no test written.
func FetchTop5Players() {
	if DBC == nil {
		return
	}

	query := DBC.NewRef("leaderboard/").OrderByChild("Score").LimitToFirst(5)
	result, err := query.GetOrdered(context.Background())
	if err != nil {
		log.Print(err)
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
