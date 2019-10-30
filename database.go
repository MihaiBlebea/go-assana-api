package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func connectDB(url string) (*sql.DB, error) {
	return sql.Open("postgres", url)
}

func createTable() {
	db, err := connectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	stm, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			task_id integer NOT NULL
		)
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = stm.Exec()
}

func checkTaskExists(taskID int) bool {
	db, err := connectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	var id int

	row := db.QueryRow(`
		SELECT
			*
		FROM
			tasks
		WHERE
			id = $1
	`, taskID)

	err = row.Scan(&id, &taskID)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		log.Panic(err)
	}

	if id != 0 {
		return true
	}
	return false
}

func addTask(taskID int) int {
	db, err := connectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	lastInsertId := 0
	err = db.QueryRow(`INSERT INTO tasks (task_id) VALUES ($1) RETURNING id`, taskID).Scan(&lastInsertId)

	if err != nil {
		log.Panic(err)
	}

	if lastInsertId == 0 {
		log.Panic("Id not found")
	}

	return lastInsertId
}
