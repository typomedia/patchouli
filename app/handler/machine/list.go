package machine

import (
	"encoding/json"
	"github.com/typomedia/patchouli/app"

	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	config := app.GetApp().Config

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
		if machine.Interval == 0 {
			machine.Interval = config.General.Interval
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
