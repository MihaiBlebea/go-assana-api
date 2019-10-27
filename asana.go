package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Client struct {
	client    *http.Client
	baseUrl   string
	token     string
	workspace string
}

func (c *Client) Project(projectID int) *Project {
	req, err := http.NewRequest("GET", c.baseUrl+"/projects/"+convertIntToString(projectID)+"?opt_expand=team,id,owner&opt_fields=name,id,owner,workspace", nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := c.Call(req)
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

func (c *Client) Projects(workspace, limit int) *[]Project {
	req, err := http.NewRequest("GET", c.baseUrl+"/projects?limit="+convertIntToString(limit)+"&workspace="+c.workspace, nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := c.Call(req)
	if err != nil {
		log.Panic(err)
	}

	projects := new([]Project)
	err = mapstructure.Decode(data["data"], &projects)
	if err != nil {
		log.Panic(err)
	}

	return projects
}

func (c *Client) ProjectTasks(projectID, limit int) *[]Task {
	req, err := http.NewRequest("GET", c.baseUrl+"/projects/"+convertIntToString(projectID)+"/tasks?opt_expand=name&opt_fields=name,id,owner,team&limit="+convertIntToString(limit), nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := c.Call(req)
	if err != nil {
		log.Panic(err)
	}

	tasks := new([]Task)
	err = mapstructure.Decode(data["data"], &tasks)
	if err != nil {
		log.Panic(err)
	}

	return tasks
}

func (c *Client) Task(taskID int) *Task {
	req, err := http.NewRequest("GET", c.baseUrl+"/tasks/"+convertIntToString(taskID), nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := c.Call(req)
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

func (c *Client) UpdateTask(taskID int, body io.Reader) bool {
	req, err := http.NewRequest("PUT", c.baseUrl+"/tasks/"+convertIntToString(taskID), body)
	if err != nil {
		log.Panic(err)
	}

	_, err = c.Call(req)
	if err != nil {
		log.Panic(err)
	}

	return true
}

func (c *Client) TaskStories(taskID, limit int) *[]Story {
	req, err := http.NewRequest("GET", c.baseUrl+"/tasks/"+convertIntToString(taskID)+"/stories?limit="+convertIntToString(limit), nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := c.Call(req)
	if err != nil {
		log.Panic(err)
	}

	stories := new([]Story)
	err = mapstructure.Decode(data["data"], &stories)
	if err != nil {
		log.Panic(err)
	}

	return stories
}

func (c *Client) NextPage(uri string) map[string]interface{} {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Panic(err)
	}

	result, err := c.Call(req)
	if err != nil {
		log.Panic(err)
	}
	return result
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

func NewAsanaSDK(token, workspace string) *Client {
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
