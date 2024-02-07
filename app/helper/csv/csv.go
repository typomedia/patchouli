package csv

import (
	"github.com/gocarina/gocsv"
	"io"
	"strings"
)

func Marshal(data interface{}) (string, error) {
	var csvString strings.Builder

	writer := io.Writer(&csvString)

	err := gocsv.Marshal(data, writer)
	if err != nil {
		return "", err
	}

	return csvString.String(), nil
}
