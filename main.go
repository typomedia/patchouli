package main

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/typomedia/patchouli/app/handler/api/csv"
	"github.com/typomedia/patchouli/app/handler/api/htmx"
	"github.com/typomedia/patchouli/app/handler/api/json"
	"github.com/typomedia/patchouli/app/handler/config"
	"github.com/typomedia/patchouli/app/handler/dashboard"
	"github.com/typomedia/patchouli/app/handler/machine"
	"github.com/typomedia/patchouli/app/handler/operator"
	"github.com/typomedia/patchouli/app/handler/system"
	"github.com/typomedia/patchouli/app/handler/update"
	"log"
	"net/http"
)

//go:embed app/views
var views embed.FS

//go:embed public
var public embed.FS

func main() {
	app := fiber.New(fiber.Config{
		Views: html.NewFileSystem(http.FS(views), ".html"),
	})

	app.Get("/", dashboard.List)

	app.Get("/machine", machine.List)
	app.Get("/machine/new", machine.New)
	app.Get("/machine/edit/:id", machine.Edit)
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

	app.Use("/", filesystem.New(filesystem.Config{
		//app.Static("/", "public")
		// publish static embedded files like css, js, images
		Root:       http.FS(public),
		PathPrefix: "public",
	}))

	log.Fatal(app.Listen(":5000"))
}
