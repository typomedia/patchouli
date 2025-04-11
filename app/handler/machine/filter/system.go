package filter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func System(c *fiber.Ctx) error {
	systemId := c.Params("id")
	db := boltdb.New()

	Machines, err := db.GetAllMachinesBySystemId(systemId)
	if err != nil {
		return err
	}

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": Machines,
	})
}
