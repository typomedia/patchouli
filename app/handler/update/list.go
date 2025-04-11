package update

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	// Todo: Remove mandatory bucket setting
	db.SetBucket("history")

	var machine structs.Machine
	err := db.Get(id, &machine, "machine")
	if err != nil {
		log.Error(err)
	}

	days, _ := db.GetAllUpdatesByMachineId(id)

	Updates := structs.Updates{}

	for _, v := range days {
		update := structs.Update{}
		err = json.Unmarshal(v, &update)
		if err != nil {
			return err
		}

		update.Date = helper.UnixToDateString(update.Date)
		Updates = append(Updates, update)

	}

	defer db.Close()

	return c.Render("app/views/update/list", fiber.Map{
		"Updates": Updates,
		"Machine": machine,
	})
}
