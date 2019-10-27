package main

type NextPage struct {
	path   string
	offset string
	Uri    string
}

type TaskCollection struct {
	Data      *[]Task
	Next_page NextPage
}

func (tc *TaskCollection) Tasks() *[]Task {
	return tc.Data
}

func (tc *TaskCollection) NextPageUrl() string {
	return tc.Next_page.Uri
}

type ProjectCollection struct {
	Data      *[]Project
	Next_page NextPage
}

func (pc *ProjectCollection) Projects() *[]Project {
	return pc.Data
}

func (pc *ProjectCollection) NextPageUrl() string {
	return pc.Next_page.Uri
}

type StoryCollection struct {
	Data      *[]Story
	Next_page NextPage
}

func (sc *StoryCollection) Stories() *[]Story {
	return sc.Data
}

func (sc *StoryCollection) NextPageUrl() string {
	return sc.Next_page.Uri
}
