package machine

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Edit(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	// set bucket
	err := db.SetBucket("machine")
	if err != nil {
		log.Error(err)
	}

	var machine structs.Machine
	err = db.Get(id, &machine, "machine")
	if err != nil {
		log.Error(err)
	}

	var config structs.Config
	config, err = db.GetConfig()
	if err != nil {
		return err
	}
	if machine.Interval == 0 {
		machine.Interval = config.General.Interval
	}
	defer db.Close()

	return c.Render("app/views/machine/edit", fiber.Map{
		"Machine": machine,
	})
}
