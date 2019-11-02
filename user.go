package main

type User struct {
	Gid             string
	Name            string
	CompletedPoints int
}

func (u *User) AddCompletedPoints(points int) {
	u.CompletedPoints = points
}
