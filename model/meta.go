package model

type Meta struct {
	Name    string
	Styles  []Style
	Version string
}

type Style struct {
	Name string
	Id   int
	Type string
}
