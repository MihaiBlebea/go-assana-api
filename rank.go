package main

type Rank struct {
	Top map[string]int
}

type Developer struct {
	Name   string
	Points int
}

func (r *Rank) AddDeveloper(name string, points int) {
	r.Top[name] += points
}

func newRank() *Rank {
	top := make(map[string]int)
	return &Rank{top}
}
