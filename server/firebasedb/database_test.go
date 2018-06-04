package firebasedb

import (
	"log"
	"testing"
)

func TestConnectDB(t *testing.T) {
	DBC.ConnectDB("BumperDB-3b7d790985b1.json")

	if DBC.DBCon == nil {
		log.Fatal("DBClient not initialized correctly")
	}
}
