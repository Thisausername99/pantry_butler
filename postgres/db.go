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
	fmt.Println(q.FormattedQuery())
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
		//TODO: Need to wrap it in .env file
		opts = &pg.Options{
			//default port
			//depends on the db service from docker compose
			Addr:     "postgres://postgres:postgres@localhost:5432/pantry_butler_dev?sslmode=disable",
			User:     "postgres",
			Password: "postgres",
		}
	}

	//connect db
	db := pg.Connect(opts)
	// //run migrations
	// collection := migrations.NewCollection()
	// err = collection.DiscoverSQLMigrations("migrations")
	// if err != nil {
	// 	return nil, err
	// }

	// //start the migrations
	// _, _, err = collection.Run(db, "init")
	// if err != nil {
	// 	return nil, err
	// }

	// oldVersion, newVersion, err := collection.Run(db, "up")
	// if err != nil {
	// 	return nil, err
	// }
	// if newVersion != oldVersion {
	// 	log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	// } else {
	// 	log.Printf("version is %d\n", oldVersion)
	// }

	//return the db connection
	return db, err
}
