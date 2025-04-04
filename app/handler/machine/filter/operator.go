package filter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func Operator(c *fiber.Ctx) error {
	operatorId := c.Params("id")
	db := boltdb.New()

	machines, err := db.GetAllMachinesByOperatorId(operatorId)
	if err != nil {
		return err
	}

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": machines,
	})
}
