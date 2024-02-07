package structs

import (
	"github.com/typomedia/patchouli/app/helper"
)

type Machines []Machine
type ByDate []Machine

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	timeA := helper.UnixToDate(a[i].Update.Date)
	timeB := helper.UnixToDate(a[j].Update.Date)
	return timeA.Before(timeB)
}

type Machine struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	System   System   `json:"system"`
	Location string   `json:"location"`
	Ip       string   `json:"ip"`
	Fqdn     string   `json:"fqdn"`
	Service  string   `json:"service"`
	Comment  string   `json:"comment"`
	Backup   string   `json:"backup"`
	Operator Operator `json:"operator"`
	Update   Update   `json:"update"` // last update
	Days     int      `json:"days"`   // remaining days
	Status   string   `json:"status"` // status color
}
