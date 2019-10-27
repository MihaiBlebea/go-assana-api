package main

type Event struct {
	User       User
	Created_At string
	Action     string
	Resource   Resource
	Parent     Resource
}

type Resource struct {
	Gid           string
	Resource_type string
}
