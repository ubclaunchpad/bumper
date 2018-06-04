package firebasedb

import (
	"log"
	"testing"
)

func TestConnectDB(t *testing.T) {
	DBC.ConnectDB()

	if DBC.DBCon == nil {
		log.Fatal("DBClient not initialized correctly")
	}
}
