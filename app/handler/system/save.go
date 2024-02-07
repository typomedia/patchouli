package system

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Save(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "new" {
		id = helper.GenerateId()
	}

	db := boltdb.New()
	err := db.SetBucket("systems")
	if err != nil {
		log.Error(err)
	}

	system := structs.System{}
	err = c.BodyParser(&system)
	if err != nil {
		log.Error(err)
	}

	system.Id = id

	db.Set(id, system, "systems")

	defer db.Close()

	return c.Redirect("/system")
}
