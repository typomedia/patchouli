package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app"
)

func Edit(c *fiber.Ctx) error {
	config := app.GetApp().Config
	if config.General.Hostname == "" {
		config.General.Hostname = fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	}
	return c.Render("app/views/config/edit", fiber.Map{
		"Config": config,
	})
}
