package machine

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Activate(c *fiber.Ctx) error {
	machineId := c.Params("id")

	db := boltdb.New()

	err := db.SetBucket("machine")
	if err != nil {
		return err
	}
	defer db.Close()

	var machine structs.Machine
	err = db.Get(machineId, &machine, "machine")
	if err != nil {
		return err
	}

	machine.Inactive = false

	err = db.Set(machineId, machine, "machine")

	if err != nil {
		return err
	}

	return c.Redirect("/machine")
}
