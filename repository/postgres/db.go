package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"time"
	_ "time/tzdata"
	"timeMachine/pkg/logger"

	_ "github.com/lib/pq"
)

type Config struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string
}

type DB struct {
	config Config
	db     *sql.DB
}

func New(config Config) *DB {
	log := logger.Get()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, os.Getenv("SECRETS_DBUSER"), os.Getenv("SECRETS_DBPASS"), config.DBName))
	if err != nil {
		log.Fatal().Msgf("can't connect to db %s", err)
		panic(err)
	}
	log.Info().Msg("connection to db is OK")

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &DB{
		config: config,
		db:     db,
	}
}

func (db *DB) GetTrafficLength() (int32, error) {
	var length int32
	var nowTehran = nowTehran()
	if err := db.db.QueryRow(`select length::int from traffic.traffic_length where date_time>$1::timestamp at time zone 'asia/tehran'-interval'20min' and zone_id=1 order by date_time desc limit 1;`, nowTehran).Scan(&length); err != nil {
		log.Warn().Msgf("get traffic length from db err is: ", err)
		return length, err
	}
	log.Info().Msgf("Traffic length for %s is %d", nowTehran, length)
	return length, nil
}

func nowTehran() string {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Fatal().Msgf("can't find location timezone ", err)
	}
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}
