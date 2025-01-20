package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/typomedia/patchouli/app/handler/api/csv"
	"github.com/typomedia/patchouli/app/handler/api/htmx"
	"github.com/typomedia/patchouli/app/handler/api/json"
	"github.com/typomedia/patchouli/app/handler/config"
	"github.com/typomedia/patchouli/app/handler/dashboard"
	dashboardFilter "github.com/typomedia/patchouli/app/handler/dashboard/filter"
	"github.com/typomedia/patchouli/app/handler/machine"
	machineFilter "github.com/typomedia/patchouli/app/handler/machine/filter"
	"github.com/typomedia/patchouli/app/handler/operator"
	"github.com/typomedia/patchouli/app/handler/system"
	"github.com/typomedia/patchouli/app/handler/update"
)

//go:embed app/views
var views embed.FS

//go:embed public
var public embed.FS

type Application struct {
	Name        string
	Version     string
	Author      string
	Description string
}

var App = Application{
	Name:        "Patchouli",
	Version:     "0.1.1",
	Author:      "Philipp Speck <philipp@typo.media>",
	Description: "Patch Management Planner",
}

func main() {
	engine := html.NewFileSystem(http.FS(views), ".html")
	engine.AddFunc("Name", func() string {
		return App.Name
	})
	engine.AddFunc("Version", func() string {
		return App.Version
	})
	engine.AddFunc("Year", func() string {
		return fmt.Sprintf("%d", time.Now().Year())
	})

	app := fiber.New(fiber.Config{
		AppName: App.Name,
		Views:   engine,
	})

	app.Get("/", dashboard.List)
	app.Get("/filter/operator/:id", dashboardFilter.Operator)

	app.Get("/machine", machine.List)
	app.Get("/machine/new", machine.New)
	app.Get("/machine/edit/:id", machine.Edit)
	app.Get("/machine/filter/operator/:id", machineFilter.Operator)
	app.Post("/machine/save/:id", machine.Save)

	app.Get("/machine/update/list/:id", update.List)
	app.Get("/machine/update/new/:machine", update.New)
	app.Get("/machine/update/edit/:id", update.Edit)
	app.Post("/machine/update/save/:id", update.Save)

	app.Get("/operator", operator.List)
	app.Get("/operator/new", operator.New)
	app.Get("/operator/edit/:id", operator.Edit)
	app.Post("/operator/save/:id", operator.Save)

	app.Get("/system", system.List)
	app.Get("/system/new", system.New)
	app.Get("/system/edit/:id", system.Edit)
	app.Post("/system/save/:id", system.Save)

	app.Get("/config", config.Edit)
	app.Post("/config/save", config.Save)

	app.Get("/api/htmx/systems", htmx.Systems)
	app.Get("/api/htmx/operators", htmx.Operators)
	app.Get("/api/htmx/operator/:machine", htmx.Operator)
	app.Get("/api/json/export", json.Export)
	app.Get("/api/csv/export", csv.Export)

	// publish static embedded files like css, js, images
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(public),
		PathPrefix: "public",
	}))

	log.Fatal(app.Listen(":5000"))
}
