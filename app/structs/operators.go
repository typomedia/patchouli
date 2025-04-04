package structs

import "strings"

type Operators []Operator

type Operator struct {
	Id         string `json:"id"`
	Name       string `json:"operator"`
	Department string `json:"department"`
	Email      string `json:"email"`
}

func (o Operator) Firstname() string {
	return strings.Split(o.Name, " ")[0]
}

func (o Operator) Lastname() string {
	return strings.Split(o.Name, " ")[1]
}
