package main

import (
	"errors"
	"time"
)

type Task struct {
	Gid           string
	UniqueId      int
	Name          string
	Assignee      User
	Completed     bool
	Completed_At  string
	Created_At    string
	Custom_Fields []CustomField
}

type Tasks []Task

type CustomField struct {
	Gid          string
	Enum_Options []EnumOption
	Enum_Value   EnumOption
	Name         string
	Type         string
	Number_Value int
}

type EnumOption struct {
	Gid   string
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

func (t *Task) GetCreatedTime() (time.Time, error) {
	ret, err := time.Parse(time.RFC3339, t.Created_At)
	if err != nil {
		return time.Now(), err
	}
	return ret, nil
}

func (t *Task) GetCompletedTime() (time.Time, error) {
	if t.Completed_At == "" {
		return time.Now(), errors.New("Task was not completed yet")
	}

	ret, err := time.Parse(time.RFC3339, t.Completed_At)
	if err != nil {
		return time.Now(), err
	}

	return ret, nil
}

func (t *Task) GetDuration() (time.Duration, error) {
	created, err := t.GetCreatedTime()
	if err != nil {
		return time.Microsecond, err
	}

	completed, err := t.GetCompletedTime()
	if err != nil {
		return time.Microsecond, err
	}

	return completed.Sub(created), nil
}

func (t *Task) GetCustomFieldValue(name string) (interface{}, error) {
	for _, field := range t.Custom_Fields {
		if field.Name == name {
			if field.Type == "number" {
				return field.Number_Value, nil
			} else if field.Type == "enum" {
				return field.Enum_Value.Name, nil
			}

			return nil, errors.New("Type not found")
		}
	}
	return nil, errors.New("Custom field not found")
}

func (t *Task) AddUniqueId(uniqueId int) {
	t.UniqueId = uniqueId
}
