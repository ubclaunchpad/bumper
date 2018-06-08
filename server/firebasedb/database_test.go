package firebasedb

import (
	"log"
	"testing"

	"github.com/ubclaunchpad/bumper/server/models"
)

func TestConnectDB(t *testing.T) {
	// Connect to DB
	DBC.ConnectDB("BumperDB-3b7d790985b1.json")

	if DBC.DBCon == nil {
		log.Fatal("DBClient not initialized correctly")
	}
}

func TestUpdateFetchPlayerScore(t *testing.T) {
	// Connect to DB
	DBC.ConnectDB("BumperDB-3b7d790985b1.json")

	// Create a Player
	p := new(models.Player)
	p.AddPoints(100)
	UpdatePlayerScore(p)

	returnScore := FetchPlayerScore(p)

	if returnScore.Name != p.Name || returnScore.Score != p.Points {
		log.Fatal("DBC did not store or retreve score correctly")
	}
}
