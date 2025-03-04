package filter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Operator(c *fiber.Ctx) error {
	operatorId := c.Params("id")
	db := boltdb.New()
	defer db.Close()

	machines, err := db.GetActiveMachines()
	if err != nil {
		return err
	}

	var operatorMachines structs.Machines
	for _, m := range machines {
		if m.Operator.Id == operatorId {
			operatorMachines = append(operatorMachines, m)
		}
	}

	return c.Render("app/views/dashboard/list", fiber.Map{
		"Machines": operatorMachines,
	})
}
