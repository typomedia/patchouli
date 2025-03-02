package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"golang.org/x/crypto/bcrypt"
)

func Edit(c *fiber.Ctx) error {
	db := boltdb.New()

	err := db.SetBucket("config")
	if err != nil {
		return err
	}

	var config structs.Config
	db.Get("main", &config, "config")
	if config.Smtp.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(config.Smtp.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		config.Smtp.Password = string(passwordHash)
	}
	defer db.Close()
	return c.Render("app/views/config/edit", fiber.Map{
		"Config": config,
	})
}
