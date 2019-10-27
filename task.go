package main

import "errors"

type Task struct {
	Id            int
	Name          string
	Assignee      User
	Completed     bool
	Completed_At  string
	Created_At    string
	Custom_Fields []CustomField
}

type Tasks []Task

type CustomField struct {
	Id           int
	Enum_Options []CustomFieldOption
	Enum_Value   CustomFieldOption
	Name         string
	Type         string
	Number_Value int
}

type CustomFieldOption struct {
	Id    int
	Color string
	Name  string
}

// Task methods

func (t *Task) GetSprint() (int, error) {
	for _, field := range t.Custom_Fields {
		if field.Type == "number" && field.Name == "Points" {
			return field.Number_Value, nil
		}
	}

	return 0, errors.New("Could not find sprint number")
}
