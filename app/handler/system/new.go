package system

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
)

func New(c *fiber.Ctx) error {
	db := boltdb.New()
	db.SetBucket("systems")

	defer db.Close()

	return c.Render("app/views/system/new", fiber.Map{})
}
