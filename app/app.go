package app

import "fmt"

type Application struct {
	Name        string
	Version     string
	Author      string
	Description string
	Explanation string
}

var App = Application{
	Name:        "Patchouli",
	Version:     "0.1.0",
	Author:      "Philipp Speck <philipp@typo.media>",
	Description: "Patch Management Tool",
	Explanation: "Git HTTP Daemon for managing multiple repositories via web hooks.",
}

func Logo() string {
	banner := fmt.Sprintf("%s %s\n", App.Name, App.Version)
	banner += App.Author

	return banner
}
