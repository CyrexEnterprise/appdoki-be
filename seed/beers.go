package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (s *seeder) seedBeerTransfers() {
	gofakeit.Seed(0)
	transfers := [][]string{
		{"1", "2", strconv.Itoa(gofakeit.Number(1, 50)), "2020-01-01"},
		{"2", "1", strconv.Itoa(gofakeit.Number(1, 50)), "2020-01-05"},
		{"3", "2", strconv.Itoa(gofakeit.Number(1, 50)), "2020-03-08"},
		{"4", "5", strconv.Itoa(gofakeit.Number(1, 50)), "2020-06-09"},
		{"4", "5", strconv.Itoa(gofakeit.Number(1, 50)), "2020-06-22"},
		{"5", "3", strconv.Itoa(gofakeit.Number(1, 50)), "2020-07-07 10:11:11"},
		{"5", "1", strconv.Itoa(gofakeit.Number(1, 50)), "2020-07-07 10:11:12"},
		{"5", "4", strconv.Itoa(gofakeit.Number(1, 50)), "2020-07-07 10:11:13"},
	}

	sqlStr := "INSERT INTO beer_transfers (giver_id, taker_id, beers, given_at) VALUES "

	for i, usr := range transfers {
		separator := ", "
		if i == len(transfers) - 1 {
			separator = ""
		}
		sqlStr = sqlStr + fmt.Sprintf("('%s', '%s', '%s', '%s')%s", usr[0], usr[1], usr[2], usr[3], separator)
	}

	_, err := s.db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("seedBeerTransfers failed: %+v", err)
	}
	log.Info("beers users")
}
