package machine

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"net/url"
)

func Save(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "new" {
		id = helper.GenerateId()
	}

	db := boltdb.New()
	err := db.SetBucket("machine")
	if err != nil {
		log.Error(err)
	}

	request := c.Request().PostArgs().String()

	var system structs.System
	err = db.Get(c.FormValue("system"), &system, "systems")
	if err != nil {
		log.Error(err)
	}

	var operator structs.Operator
	err = db.Get(c.FormValue("operator"), &operator, "operator")
	if err != nil {
		log.Error(err)
	}

	//machine := structs.Machine{}
	//err = c.BodyParser(&machine)
	//if err != nil {
	//	log.Error(err)
	//}

	machine := structs.Machine{}
	helper.DecodeQuery(request, &machine)

	values, _ := url.ParseQuery(request)
	if values.Get("internet_access") == "on" {
		machine.InternetAccess = true
	} else {
		machine.InternetAccess = false
	}

	machine.Id = id
	machine.System = system
	machine.Operator = operator

	db.Set(id, machine, "machine")

	defer db.Close()

	return c.Redirect(values.Get("_referer"))
}
