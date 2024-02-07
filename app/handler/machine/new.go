package machine

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func New(c *fiber.Ctx) error {
	db := boltdb.New()
	db.SetBucket("machine")

	defer db.Close()

	return c.Render("app/views/machine/new", fiber.Map{})
}
