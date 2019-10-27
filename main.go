package main

import (
	"log"
	"time"

	"encoding/json"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

const baseUrl = "https://app.asana.com/api/1.0"

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

	http.ListenAndServe(":"+getPort(), router)
}

func getPort() string {
	if os.Getenv("PORT") != "" {
		return os.Getenv("PORT")
	}
	return "8083"
}

// client id 1146727810660992

// client secret 6ab659430e63172a5ef7096ef632a0f9

// Sprint review token 0/140204d60a962b96ee4cfcb4ec9dff87
