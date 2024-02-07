package structs

type Operators []Operator

type Operator struct {
	Id         string `json:"id"`
	Name       string `json:"operator"`
	Department string `json:"department"`
	Email      string `json:"email"`
}
