package operator

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Activate(c *fiber.Ctx) error {
	operatorId := c.Params("id")

	db := boltdb.New()

	err := db.SetBucket("operator")
	if err != nil {
		return err
	}
	defer db.Close()

	var operator structs.Operator
	err = db.Get(operatorId, &operator, "operator")
	if err != nil {
		return err
	}

	operator.Inactive = false

	err = db.Set(operatorId, operator, "operator")

	if err != nil {
		return err
	}

	return c.Redirect("/operator")
}
