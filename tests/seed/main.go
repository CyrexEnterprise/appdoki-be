package seed

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
)

type seeder struct {
	db *sqlx.DB
}

func main() {
	log.Infoln("seeding...")
	dbURI := os.Args[1]
	if len(dbURI) == 0 {
		log.Fatal("missing DB_URI")
	}

	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatalln(err)
	}

	Seed(db)
}

func Seed(db *sqlx.DB) {
	s := seeder{db}
	s.truncateAll()
	s.seedUsers()
	s.seedBeerTransfers()
}

func (s *seeder) truncateAll() {
	_, err := s.db.Exec("TRUNCATE beer_transfers, users RESTART IDENTITY;")
	if err != nil {
		log.Fatalf("seedBeerTransfers failed: %+v", err)
	}
	log.Debug("seeder: truncated all tables")
}