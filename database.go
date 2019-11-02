package main

import (
	"database/sql"
	"fmt"
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
			gid VARCHAR(250) NOT NULL,
			created timestamp default NULL
		)
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = stm.Exec()
}

func checkTaskExists(gid string) bool {
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
			gid = $1
	`, gid)

	err = row.Scan(&id, &gid)

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

func addTask(task Task) int {
	db, err := connectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	lastInsertId := 0
	err = db.QueryRow(`
		INSERT INTO tasks 
			(gid, created) 
		VALUES 
			($1, $2)
		RETURNING id`, task.Gid, task.Created_At).Scan(&lastInsertId)

	if err != nil {
		log.Panic(err)
	}

	if lastInsertId == 0 {
		log.Panic("Id not found2")
	}

	return lastInsertId
}

func getTasks() []DatabaseTask {
	db, err := connectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	row, err := db.Query(`
		SELECT
			*
		FROM
			tasks
	`)

	var (
		id      int
		gid     string
		created string
	)

	var tasks []DatabaseTask
	for row.Next() {
		err = row.Scan(&id, &gid, &created)
		if err != nil {
			fmt.Println(err)
			continue
		}

		task := DatabaseTask{id, gid, created}
		tasks = append(tasks, task)
	}

	return tasks
}
