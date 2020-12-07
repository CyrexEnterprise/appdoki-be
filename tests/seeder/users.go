package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// SeedUsers generates and inserts a batch of 5 users
func (s *Seeder) SeedUsers() {
	gofakeit.Seed(0)
	sqlStr := "INSERT INTO users (id, name, email, picture) VALUES "

	for u := 1; u <= 5; u++ {
		usr := generateUserData(strconv.Itoa(u))
		separator := ", "
		if u == 5 {
			separator = ""
		}
		sqlStr = sqlStr + fmt.Sprintf("('%s', '%s', '%s', '%s')%s", usr[0], usr[1], usr[2], usr[3], separator)
	}

	_, err := s.db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("SeedUsers failed: %+v", err)
	}
}

func generateUserData(id string) []string {
	return []string{id, gofakeit.Name(), gofakeit.Email(), gofakeit.URL()}
}

// SeedUsers truncates the users table and all tables that have foreign-key references
func (s *Seeder) TruncateUsers() {
	_, err := s.db.Exec("TRUNCATE users RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("TruncateUsers failed: %+v", err)
	}
}
