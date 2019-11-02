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

	router.GET("/sprint/review/:number", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		param := ps.ByName("number")

		sprintNumber, err := strconv.Atoi(param)
		if err != nil {
			log.Panic(err)
		}

		tasks := client.SprintTasks(sprintNumber)

		var (
			fullTasks       []Task
			totalPoints     int
			completedPoints int
		)

		for _, task := range *tasks {
			fullTask := client.Task(task.Gid)

			taskPoint, err := fullTask.GetCustomFieldValue("Points")
			if err != nil {
				log.Panic(err)
			}
			totalPoints += taskPoint.(int)

			status, err := fullTask.GetCustomFieldValue("Status")
			if err != nil {
				log.Panic(err)
			}
			if status == "Complete" {
				completedPoints += taskPoint.(int)
			}

			fullTasks = append(fullTasks, *fullTask)
		}

		type Response struct {
			Tasks           int
			Points          int
			PointsPerTask   int
			PointsPerDay    int
			PointsPerPerson int
			CompletedPoints int
			PointsLeft      int
		}

		resp, err := json.Marshal(Response{
			Tasks:           len(fullTasks),
			Points:          totalPoints,
			PointsPerTask:   totalPoints / len(fullTasks),
			PointsPerDay:    totalPoints / 10,
			PointsPerPerson: totalPoints / 4,
			CompletedPoints: completedPoints,
			PointsLeft:      totalPoints - completedPoints,
		})
		if err != nil {
			log.Panic(err)
		}

		w.Write(resp)
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
			if event.wasTaskAdded() == true {

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

			if event.wasTaskDeleted() == true {
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
