package update

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Save(c *fiber.Ctx) error {
	id := c.Params("id")

	sendMail := false
	if id == "new" {
		id = helper.GenerateId()
		sendMail = true
	}

	db := boltdb.New()
	defer db.Close()
	db.SetBucket("history")

	var machine structs.Machine
	err := db.Get(c.FormValue("machine"), &machine, "machine")
	if err != nil {
		log.Error(err)
	}

	var operator structs.Operator
	err = db.Get(c.FormValue("operator"), &operator, "operator")
	if err != nil {
		log.Error(err)
	}

	request := c.Request().PostArgs().String()
	update := structs.Update{}
	helper.DecodeQuery(request, &update)
	update.Id = id
	// day.Date to unix timestamp
	update.Date = helper.DateToUnixString(update.Date)
	update.Operator = operator

	if update.Id == "" {
		update.Id = helper.GenerateId()
	}

	err = db.Set(update.Id, update, "history")
	if err != nil {
		return err
	}
	if sendMail {
		return c.Redirect("/machine/update/mail/send/" + update.Id)
	}
	return c.Redirect("/machine/update/list/" + machine.Id)
}
