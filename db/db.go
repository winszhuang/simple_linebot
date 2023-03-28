package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() error {
	userName := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	source := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", userName, password, dbName)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return err
	}
	defer db.Close()

	// 檢查連接是否正常
	err = db.Ping()
	if err != nil {
		return err
	}

	// 插入數據
	_, err = db.Exec("INSERT INTO users (line_id, name, language, picture) VALUES ($1, $2, $3, $4)", "regregijil", "Tom", "en", "http://example.com/pic.jpg")
	if err != nil {
		return err
	}

	fmt.Println("Insert successful.")
	return nil
}
