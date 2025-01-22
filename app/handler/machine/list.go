package machine

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	// set bucket
	err := db.SetBucket("machine")
	if err != nil {
		return err
	}

	machines, _ := db.GetAll("machine")

	Machines := structs.Machines{}
	inactiveMachines := structs.Machines{}

	for _, v := range machines {
		machine := structs.Machine{}
		err = json.Unmarshal(v, &machine)
		if err != nil {
			return err
		}
		if machine.Inactive {
			inactiveMachines = append(inactiveMachines, machine)
		} else {
			Machines = append(Machines, machine)

		}
	}

	Machines = append(Machines, inactiveMachines...)

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": Machines,
	})
}
