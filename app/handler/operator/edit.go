package operator

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Edit(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	err := db.SetBucket("operator")
	if err != nil {
		return err
	}

	var operator structs.Operator
	err = db.Get(id, &operator, "operator")
	if err != nil {
		log.Error(err)
	}

	machinesOfOperator, err := db.GetMachinesByOperator(operator.Id)
	if err != nil {
		return err
	}
	operator.MachineCount = len(machinesOfOperator)

	defer db.Close()

	return c.Render("app/views/operator/edit", fiber.Map{
		"Operator": operator,
	})
}
