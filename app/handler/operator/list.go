package operator

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	operators, err := db.GetAllOperators()
	if err != nil {
		return err
	}

	defer db.Close()

	return c.Render("app/views/operator/list", fiber.Map{
		"Operators": operators,
	})
}
