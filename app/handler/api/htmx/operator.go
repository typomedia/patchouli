package htmx

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Operator(c *fiber.Ctx) error {
	id := c.Params("machine")
	selected := c.Query("selected")

	db := boltdb.New()

	err := db.SetBucket("machine")
	if err != nil {
		return err
	}

	machine := structs.Machine{}
	err = db.Get(id, &machine, "machine")
	if err != nil {
		log.Error(err)
	}

	err = db.SetBucket("operator")
	if err != nil {
		return err
	}

	operators, _ := db.GetAll("operator")

	Operators := structs.Operators{}
	for _, v := range operators {
		operator := structs.Operator{}
		err = json.Unmarshal(v, &operator)
		if err != nil {
			return err
		}
		if operator.Inactive {
			continue
		}
		Operators = append(Operators, operator)
	}

	defer db.Close()

	if selected == "" {
		selected = machine.Operator.Id
	}

	return c.Render("app/views/api/htmx/operators", fiber.Map{
		"Operators": Operators,
		"Selected":  selected,
	})
}
