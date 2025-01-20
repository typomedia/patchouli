package filter

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Operator(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	// set bucket
	err := db.SetBucket("machine")
	if err != nil {
		return err
	}

	machines, _ := db.GetAllByOperatorId(id, "machine")

	Machines := structs.Machines{}

	for _, v := range machines {
		machine := structs.Machine{}
		err = json.Unmarshal(v, &machine)
		if err != nil {
			return err
		}
		Machines = append(Machines, machine)

	}

	defer db.Close()

	return c.Render("app/views/machine/list", fiber.Map{
		"Machines": Machines,
	})
}
