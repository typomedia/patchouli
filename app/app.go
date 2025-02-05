package app

var app *Application

type Application struct {
	Name         string
	Version      string
	Author       string
	Description  string
	MailTemplate string
}

func new() *Application {
	app = &Application{}
	app.Name = "Patchouli"
	app.Version = "0.2.0"
	app.Author = "Philipp Speck <philipp@typo.media>"
	app.Description = "Patch Management Planner"

	return app
}
func GetApp() *Application {
	if app == nil {
		app = new()
	}
	return app
}
