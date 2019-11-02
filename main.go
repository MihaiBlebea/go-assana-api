package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"

	_ "github.com/lib/pq"
)

const baseUrl = "https://app.asana.com/api/1.0"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
	createTable()
}

func main() {
	if os.Getenv("ASANA_TOKEN") == "" {
		log.Panic("No Asana token found")
	}

	if os.Getenv("ASANA_WORKSPACE") == "" {
		log.Panic("No Asana workspace found")
	}

	client := NewClient(os.Getenv("ASANA_TOKEN"), os.Getenv("ASANA_WORKSPACE"))

	router := httprouter.New()

	router.GET("/projects", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		startTime := time.Now()

		projects := client.Projects()

		resp, err := json.Marshal(projects)
		if err != nil {
			log.Panic(err)
		}

		w.Write(resp)

		completeTime := time.Now().Sub(startTime)
		log.Println(completeTime)
	})

	router.GET("/tasks", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		startTime := time.Now()

		dTasks := getTasks()

		var tasks []Task
		for _, dTask := range dTasks {
			task := client.Task(dTask.Gid)

			task.AddUniqueId(dTask.Id)

			tasks = append(tasks, *task)
		}
		resp, err := json.Marshal(tasks)
		if err != nil {
			log.Panic(err)
		}

		w.Write(resp)

		completeTime := time.Now().Sub(startTime)
		log.Println(completeTime)
	})

	router.POST("/webhook/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := jsonParse(r.Body)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println(data)

		var events []Event
		err = mapstructure.Decode(data["events"], &events)
		if err != nil {
			log.Panic(err)
		}

		for _, event := range events {
			if event.Resource.Resource_type == "task" && event.Action == "added" {

				found := checkTaskExists(event.Resource.Gid)

				fmt.Println(found)

				if found == true {
					continue
				}

				taskUID := addTask(DatabaseTask{
					Gid:     event.Resource.Gid,
					Created: event.Created_at,
				})

				data := map[string]map[string]map[string]int{
					"data": {
						"custom_fields": {
							"1146745094235892": taskUID,
						},
					},
				}

				encoded, err := json.Marshal(data)
				if err != nil {
					log.Panic(err)
				}

				reader := bytes.NewReader(encoded)

				client.UpdateTask(event.Resource.Gid, reader)
			}

			if event.Resource.Resource_type == "task" && event.Action == "deleted" {
				// Check if the task already exists in the database
				found := checkTaskExists(event.Resource.Gid)

				if found == false {
					continue
				}
				deleteTask(event.Resource.Gid)
			}
		}

		w.Header().Add("X-Hook-Secret", r.Header.Get("X-Hook-Secret"))
		w.WriteHeader(200)
	})

	http.ListenAndServe(":"+getPort(), router)
}

func getPort() string {
	if os.Getenv("PORT") != "" {
		return os.Getenv("PORT")
	}
	return "8083"
}

func getRandomId(max int) int {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return random.Intn(max)
}
