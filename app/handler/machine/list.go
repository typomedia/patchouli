package machine

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	machines, _ := db.GetAllMachines(false)

	var Machines, active, inactive structs.Machines

	for _, machine := range machines {
		if machine.Inactive {
			inactive = append(inactive, machine)
		} else {
			active = append(active, machine)

		}
	}

	Machines = append(active, inactive...)

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": Machines,
	})
}
