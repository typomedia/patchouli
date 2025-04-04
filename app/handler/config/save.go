package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app"
	"github.com/typomedia/patchouli/app/encryption"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Save(c *fiber.Ctx) error {
	db := boltdb.New()
	err := db.SetBucket("config")
	if err != nil {
		log.Error(err)
	}

	config := structs.Config{}
	err = c.BodyParser(&config)
	if err != nil {
		log.Error(err)
	}

	// update config from body with security params from app
	config.Security = app.GetApp().Config.Security

	_, err = encryption.DecryptString(config.Smtp.Password)
	if err != nil { // config.Smtp.Password is not encrypted, encrypt it
		config.Smtp.Password, err = encryption.EncryptString(config.Smtp.Password)
		if err != nil {
			log.Error(err)
		}
	}

	db.Set("main", config, "config")

	app.GetApp().Config = config

	defer db.Close()

	return c.Redirect("/config")
}
