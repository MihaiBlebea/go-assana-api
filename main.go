package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

		projects := client.Projects()

		fs := (*projects)[6]

		tasks := client.ProjectTasks(fs.Id)

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

		var events []Event
		err = mapstructure.Decode(data["events"], &events)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("New webhook payload received")

		for _, event := range events {
			if event.Resource.Resource_type == "task" && event.Action == "added" {

				gid, err := strconv.Atoi(event.Resource.Gid)
				if err != nil {
					log.Panic(err)
				}
				task := client.Task(gid)

				for _, field := range task.Custom_Fields {
					fmt.Println("Looping")
					if field.Name == "Unique ID" && field.Number_Value == 0 {

						// Check if the task exists in the database
						found := checkTaskExists(task.Id)

						fmt.Println("FOUND", found)

						if found == false {
							continue
						}
						taskUID := addTask(task.Id)

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

						fmt.Println(data)
						client.UpdateTask(task.Id, reader)
					}
				}
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
