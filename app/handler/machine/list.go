package machine

import (
	"cmp"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"slices"
	"strings"
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

	slices.SortFunc(active, func(a, b structs.Machine) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	slices.SortFunc(inactive, func(a, b structs.Machine) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	Machines = append(active, inactive...)

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": Machines,
	})
}
