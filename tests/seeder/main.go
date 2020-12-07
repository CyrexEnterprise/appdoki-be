package seeder

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Seeder struct {
	db *sqlx.DB
}

func NewSeeder(dbURI string) *Seeder {
	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatalln(err)
	}

	return &Seeder{db: db}
}

func (s *Seeder) TruncateAll() {
	_, err := s.db.Exec("TRUNCATE beer_transfers, users RESTART IDENTITY;")
	if err != nil {
		log.Fatalf("truncateAll failed: %+v", err)
	}
	log.Info("seeder: truncated all tables")
}
