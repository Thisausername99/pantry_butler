package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
)

type DBLogger struct{}

func (d DBLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d DBLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}

func StartDB() (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)

	//check if we are in prod
	//then use the db url from the env
	if os.Getenv("ENV") == "PROD" {
		opts, err = pg.ParseURL(os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}
	} else {
		opts = &pg.Options{
			//default port
			//depends on the db service from docker compose
			Addr:     "localhost:5432",
			User:     "postgres",
			Password: "postgres",
			Database: "pantry_butler_dev",
		}
	}

	//connect db
	db := pg.Connect(opts)

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
	return db, err
}
