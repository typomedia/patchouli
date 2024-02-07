package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

	db.Set("main", config, "config")

	defer db.Close()

	return c.Redirect("/config")
}
