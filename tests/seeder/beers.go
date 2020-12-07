package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func (s *Seeder) SeedBeerTransfers() {
	gofakeit.Seed(0)

	var transfers [][]string
	for i := 1; i <= 50; i++ {
		transfers = append(transfers, generateRandomBeerTransfer())
	}

	sqlStr := "INSERT INTO beer_transfers (giver_id, taker_id, beers, given_at) VALUES "

	for i, usr := range transfers {
		separator := ", "
		if i == len(transfers)-1 {
			separator = ""
		}
		sqlStr = sqlStr + fmt.Sprintf("('%s', '%s', '%s', '%s')%s", usr[0], usr[1], usr[2], usr[3], separator)
	}

	_, err := s.db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("SeedBeerTransfers failed: %+v", err)
	}
}

func generateRandomBeerTransfer() []string {
	giverID := gofakeit.Number(1, 5)
	receiverID := giverID
	for receiverID == giverID {
		receiverID = gofakeit.Number(1, 5)
	}

	return []string{
		strconv.Itoa(giverID),
		strconv.Itoa(receiverID),
		strconv.Itoa(gofakeit.Number(1, 50)),
		gofakeit.DateRange(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format(time.RFC3339),
	}
}

func (s *Seeder) TruncateBeerTransfers() {
	_, err := s.db.Exec("TRUNCATE beer_transfers RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("TruncateBeerTransfers failed: %+v", err)
	}
}