package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

const (
	limit = 20
)

type Client struct {
	client    *http.Client
	baseUrl   string
	token     string
	workspace string
}

func (c *Client) Project(projectID int) *Project {
	data, err := c.GET(c.baseUrl + "/projects/" + convertIntToString(projectID) + "?opt_expand=team,id,owner&opt_fields=name,id,owner,workspace")
	if err != nil {
		log.Panic(err)
	}

	project := new(Project)
	err = mapstructure.Decode(data["data"], &project)
	if err != nil {
		log.Panic(err)
	}

	return project
}

func (c *Client) Projects() *[]Project {
	data, err := c.GET(c.baseUrl + "/projects?limit=" + convertIntToString(limit) + "&workspace=" + c.workspace)
	if err != nil {
		log.Panic(err)
	}

	var projects []Project
	var collection ProjectCollection
	err = mapstructure.Decode(data, &collection)
	if err != nil {
		log.Panic(err)
	}

	projects = append(projects, *collection.Projects()...)

	for collection.NextPageUrl() != "" {

		data, err := c.GET(collection.NextPageUrl())
		if err != nil {
			log.Panic(err)
		}

		var col ProjectCollection
		err = mapstructure.Decode(data, &col)
		if err != nil {
			log.Panic(err)
		}

		collection = col
		fmt.Println(col.NextPageUrl())
		projects = append(projects, *collection.Projects()...)
	}
	return &projects
}

func (c *Client) ProjectTasks(projectID int) *[]Task {
	data, err := c.GET(c.baseUrl + "/projects/" + convertIntToString(projectID) + "/tasks?opt_expand=name&opt_fields=name,id,owner,team&limit=" + convertIntToString(limit))
	if err != nil {
		log.Panic(err)
	}

	var tasks []Task
	var collection TaskCollection
	err = mapstructure.Decode(data, &collection)
	if err != nil {
		log.Panic(err)
	}

	tasks = append(tasks, *collection.Tasks()...)

	for collection.NextPageUrl() != "" {
		data, err := c.GET(collection.NextPageUrl())
		if err != nil {
			log.Panic(err)
		}

		var col TaskCollection
		err = mapstructure.Decode(data, &col)
		if err != nil {
			log.Panic(err)
		}

		collection = col
		fmt.Println(col.NextPageUrl())
		tasks = append(tasks, *collection.Tasks()...)
	}

	return &tasks
}

func (c *Client) Task(taskID int) *Task {
	data, err := c.GET(c.baseUrl + "/tasks/" + convertIntToString(taskID))
	if err != nil {
		log.Panic(err)
	}

	task := new(Task)
	err = mapstructure.Decode(data["data"], &task)
	if err != nil {
		log.Panic(err)
	}

	return task
}

func (c *Client) UpdateTask(taskID int, body io.Reader) *Task {
	data, err := c.PUT(c.baseUrl+"/tasks/"+convertIntToString(taskID), body)
	if err != nil {
		log.Panic(err)
	}

	task := new(Task)
	err = mapstructure.Decode(data["data"], &task)
	if err != nil {
		log.Panic(err)
	}

	return task
}

func (c *Client) TaskStories(taskID int) *[]Story {
	data, err := c.GET(c.baseUrl + "/tasks/" + convertIntToString(taskID) + "/stories?limit=" + convertIntToString(limit))
	if err != nil {
		log.Panic(err)
	}

	var stories []Story
	var collection StoryCollection
	err = mapstructure.Decode(data, &collection)
	if err != nil {
		log.Panic(err)
	}

	stories = append(stories, *collection.Stories()...)

	for collection.NextPageUrl() != "" {

		data, err := c.GET(collection.NextPageUrl())
		if err != nil {
			log.Panic(err)
		}

		var col StoryCollection
		err = mapstructure.Decode(data, &col)
		if err != nil {
			log.Panic(err)
		}

		collection = col
		stories = append(stories, *collection.Stories()...)
	}

	return &stories
}

func (c *Client) Call(req *http.Request) (map[string]interface{}, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	result, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return jsonParse(result.Body)
}

func (c *Client) GET(url string) (map[string]interface{}, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	result, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return jsonParse(result.Body)
}

func (c *Client) PUT(url string, body io.Reader) (map[string]interface{}, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	result, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return jsonParse(result.Body)
}

func NewClient(token, workspace string) *Client {
	client := &http.Client{}
	return &Client{client, "https://app.asana.com/api/1.0", token, workspace}
}

func jsonParse(data io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}

	err := json.NewDecoder(data).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func convertIntToString(value int) string {
	return strconv.Itoa(value)
}
