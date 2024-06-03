package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func NewPoolPostgres(pgUrl string) (*pgxpool.Pool, error) {
	logger := log.New(os.Stdout, "pgx: ", log.LstdFlags)

	queryLoggerTracer := &QueryLoggerTracer{Logger: logger}

	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}
	config.ConnConfig.Tracer = queryLoggerTracer

	DB, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err = DB.Ping(context.Background()); err != nil {
		fmt.Println("can't ping DataBase Link: ", err)
		return nil, err
	}
	fmt.Println("connected to postgres")
	return DB, nil
}
