package main

type Event struct {
	User       User
	Created_at string
	Action     string
	Resource   Resource
	Parent     Resource
}

type Resource struct {
	Gid           string
	Resource_type string
	Created_at    string
}

func (e *Event) wasTaskAdded() bool {
	if e.Resource.Resource_type == "task" && e.Action == "added" {
		return true
	}
	return false
}

func (e *Event) wasTaskDeleted() bool {
	if e.Resource.Resource_type == "task" && e.Action == "deleted" {
		return true
	}
	return false
}
