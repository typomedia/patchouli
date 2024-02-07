package structs

type Systems []System

type System struct {
	Id   string
	Name string
	LTS  bool
	EOL  string
}
