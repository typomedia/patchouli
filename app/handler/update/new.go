package update

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"time"
)

func New(c *fiber.Ctx) error {
	machine := c.Params("machine")
	db := boltdb.New()
	db.SetBucket("history")

	var Machine structs.Machine
	err := db.Get(machine, &Machine, "machine")
	if err != nil {
		log.Error(err)
	}

	defer db.Close()

	today := time.Now().Format("2006-01-02")

	return c.Render("app/views/update/new", fiber.Map{
		"Machine": Machine,
		"Date":    today,
	})
}
