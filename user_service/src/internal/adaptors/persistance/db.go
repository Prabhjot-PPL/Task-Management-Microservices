package persistance

import (
	"database/sql"
	"fmt"
	"user_service/src/internal/config"
	logger "user_service/src/pkg/logger"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectToDatabase() (*Database, error) {

	config := config.LoadConfig()

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.DB_User, config.DB_Password, config.DB_Host, config.DB_Port, config.DB_Name)

	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, err
	}

	fmt.Print("\n")
	logger.Log.Info("Connected to Database")

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// goose.SetDialect("postgres")

	// err = goose.Up(db, "src/migrations")

	// if err != nil {
	// 	logger.Log.Fatal(err)
	// }

	return &Database{db: db}, nil
}
