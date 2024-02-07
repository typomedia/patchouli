package system

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Edit(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	err := db.SetBucket("systems")
	if err != nil {
		return err
	}

	var system structs.System
	err = db.Get(id, &system, "systems")
	if err != nil {
		log.Error(err)
	}

	defer db.Close()

	return c.Render("app/views/system/edit", fiber.Map{
		"System": system,
	})
}
