package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"timeMachine/pkg/logger"
	"timeMachine/pkg/util"

	_ "github.com/lib/pq"
)

var log = logger.Get()

type Config struct {
	Host   string `koanf:"host"`
	Port   int    `koanf:"port"`
	User   string `koanf:"user"`
	Pass   string `koanf:"pass"`
	DBName string `koanf:"dbname"`
}

type DB struct {
	config Config
	db     *sql.DB
}

func New(config Config) *DB {
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

func (db *DB) GetTrafficLength(zone int8) (int32, error) {
	var length int32
	var nowTehran = util.Now("Asia/Tehran")
	if err := db.db.QueryRow(`select length::int from traffic.traffic_length where date_time>$1::timestamp at time zone 'asia/tehran'-interval'20min' and zone_id=$2 order by date_time desc limit 1;`, nowTehran, zone).Scan(&length); err != nil {
		log.Warn().Msgf("get traffic length from db err is: ", err)
		return length, err
	}
	log.Info().Msgf("Traffic length for %s is %d", nowTehran, length)
	return length, nil
}

func (db *DB) GetOnlineConfig(city string) (float64, error) {
	var coef float64
	if err := db.db.QueryRow(`select coefficient from taha_temp.time_machine_config where city=$1 limit 1;`, city).Scan(&coef); err != nil {
		log.Warn().Msgf("get model coefficient  from db err is: ", err)
		return coef, err
	}
	log.Info().Msgf("model coefficient is %d", coef)
	return coef, nil
}
