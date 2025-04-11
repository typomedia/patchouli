package structs

import "time"

type Systems []System

type System struct {
	Id           string
	Name         string
	LTS          bool
	EOL          string
	MachineCount int
}

func (s System) IsEOL() bool {
	eolTime, err := time.Parse(time.DateOnly, s.EOL)
	if err != nil {
		return false
	}
	return eolTime.UTC().Before(time.Now().UTC())
}
