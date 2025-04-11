package system

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		log.Error(errors.New("id is required"))
	}

	db := boltdb.New()
	err := db.SetBucket("systems")
	if err != nil {
		log.Error(err)
	}

	system := structs.System{}
	err = db.Get(id, &system, "systems")
	if err != nil {
		log.Error(err)
	}

	err = db.Delete(system.Id, "systems")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	return c.Redirect("/system")
}
