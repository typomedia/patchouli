package app

import (
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

var app *Application

type Application struct {
	Name           string
	Version        string
	Author         string
	Description    string
	MailTemplate   string
	NotifyTemplate string
	Config         structs.Config
}

func new() *Application {
	app = &Application{}
	app.Name = "Patchouli"
	app.Version = "0.3.0"
	app.Author = "Philipp Speck <philipp@typo.media>"
	app.Description = "Patch Management Planner"
	err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	if app.Config.Security.Generated {
		err = app.StoreConfig()
		if err != nil {
			panic(err)
		}
	}

	return app
}

func (app *Application) StoreConfig() error {
	db := boltdb.New()
	defer db.Close()

	err := db.SetBucket("config")
	if err != nil {
		return err
	}

	err = db.Set("main", app.Config, "config")
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) LoadConfig() error {
	db := boltdb.New()
	defer db.Close()

	config, err := db.GetConfig()
	if err != nil {
		return err
	}

	if !config.Security.Generated {
		err = config.GenerateCipherKey()
		if err != nil {
			return err
		}
	}

	app.Config = config
	return nil
}

func GetApp() *Application {
	if app == nil {
		app = new()
	}
	return app
}
