package main

type Project struct {
	Gid      string
	Name     string
	Archived bool
	Members  []User
}
