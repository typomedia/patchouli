package dashboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()
	defer db.Close()

	machines, err := db.GetActiveMachines()
	if err != nil {
		return err
	}

	return c.Render("app/views/dashboard/list", fiber.Map{
		"Machines": machines,
	})
}
