package structs

type Updates []Update

type Update struct {
	Id          string   `json:"id"`
	Machine     string   `json:"machine"`
	Date        string   `json:"date"`
	Operator    Operator `json:"operator"`
	Description string   `json:"description"`
}
