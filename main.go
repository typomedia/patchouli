package main

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/robfig/cron/v3"
	patchouli "github.com/typomedia/patchouli/app"
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
	"github.com/typomedia/patchouli/app/notifier"
	"net/http"
	"time"
)

//go:embed app/views
var views embed.FS

//go:embed public
var public embed.FS

//go:embed public/html/mail/update.html
var mailTemplate string

//go:embed public/html/mail/notify.html
var notifyTemplate string

func main() {
	App := patchouli.GetApp()
	App.MailTemplate = mailTemplate
	App.NotifyTemplate = notifyTemplate
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
	app.Get("/machine/filter/system/:id", machineFilter.System)
	app.Post("/machine/save/:id", machine.Save)

	app.Get("/machine/update/list/:id", update.List)
	app.Get("/machine/update/new/:machine", update.New)
	app.Get("/machine/update/edit/:id", update.Edit)
	app.Get("/machine/update/mail/send/:id", update.Mail)
	app.Post("/machine/update/save/:id", update.Save)

	app.Get("/machine/deactivate/:id", machine.Deactivate)
	app.Get("/machine/activate/:id", machine.Activate)

	app.Get("/operator", operator.List)
	app.Get("/operator/new", operator.New)
	app.Get("/operator/edit/:id", operator.Edit)
	app.Post("/operator/save/:id", operator.Save)

	app.Get("/operator/deactivate/:id", operator.Deactivate)
	app.Get("/operator/activate/:id", operator.Activate)

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

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		n := notifier.Notifier{}
		c := cron.New()
		_, err := c.AddFunc("0 5 * * 1", n.Run)
		if err != nil {
			return err
		}
		c.Start()
		return nil
	})

	log.Fatal(app.Listen(":5000"))
}
