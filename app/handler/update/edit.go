package update

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Edit(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	// set bucket
	err := db.SetBucket("history")
	if err != nil {
		return err
	}

	var update structs.Update
	err = db.Get(id, &update, "history")
	if err != nil {
		log.Error(err)
	}
	update.Date = helper.UnixToDateString(update.Date)

	defer db.Close()

	return c.Render("app/views/update/edit", fiber.Map{
		"Update": update,
	})
}
