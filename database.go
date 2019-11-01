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
			task_gid VARCHAR(250) NOT NULL,
			name VARCHAR(2000) NOT NULL,
			created timestamp default NULL,
			in_progress timestamp default NULL,
			completed timestamp default NULL
		)
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = stm.Exec()
}

func checkTaskExists(taskGid string) bool {
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
			task_gid = $1
	`, taskGid)

	err = row.Scan(&id, &taskGid)

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
			(task_gid, name, created, in_progress, completed) 
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING id`, task.Gid, task.Name, task.Created_At, task.Created_At, task.Created_At).Scan(&lastInsertId)

	if err != nil {
		log.Panic(err)
	}

	if lastInsertId == 0 {
		log.Panic("Id not found2")
	}

	return lastInsertId
}

func getTasks() []Task {
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
		id          int
		taskGid     string
		name        string
		created     string
		in_progress string
		completed   string
	)

	var tasks []Task
	for row.Next() {
		err = row.Scan(&id, &taskGid, &name, &created, &in_progress, &completed)
		if err != nil {
			fmt.Println(err)
			continue
		}

		task := Task{
			Gid:          taskGid,
			Name:         name,
			Completed_At: completed,
			Created_At:   created,
		}

		tasks = append(tasks, task)
	}

	return tasks
}
