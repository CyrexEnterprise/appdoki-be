package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	log "github.com/sirupsen/logrus"
)

func (s *seeder) seedUsers() {
	gofakeit.Seed(0)
	users := [][]string{
		{"1", gofakeit.Name(), gofakeit.Email(), gofakeit.URL()},
		{"2", gofakeit.Name(), gofakeit.Email(), gofakeit.URL()},
		{"3", gofakeit.Name(), gofakeit.Email(), gofakeit.URL()},
		{"4", gofakeit.Name(), gofakeit.Email(), gofakeit.URL()},
		{"5", gofakeit.Name(), gofakeit.Email(), gofakeit.URL()},
	}

	sqlStr := "INSERT INTO users (id, name, email, picture) VALUES "

	for i, usr := range users {
		separator := ", "
		if i == len(users)-1 {
			separator = ""
		}
		sqlStr = sqlStr + fmt.Sprintf("('%s', '%s', '%s', '%s')%s", usr[0], usr[1], usr[2], usr[3], separator)
	}

	_, err := s.db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("seedUsers failed: %+v", err)
	}
	log.Info("seeded users")
}
