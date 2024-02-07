package csv

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper/csv"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Export(c *fiber.Ctx) error {
	db := boltdb.New()

	err := db.SetBucket("machine")
	if err != nil {
		return err
	}

	machines, _ := db.GetAll("machine")

	Machines := structs.Machines{}

	for _, v := range machines {
		machine := structs.Machine{}
		err = json.Unmarshal(v, &machine)
		if err != nil {
			return err
		}
		Machines = append(Machines, machine)

	}

	defer db.Close()

	csvString, err := csv.Marshal(Machines)
	if err != nil {
		log.Error(err)
	}

	csvBytes := []byte(csvString)

	c.Response().Header.Set("Content-Type", "text/csv")
	c.Response().Header.Set("Content-Disposition", "attachment; filename=machines.csv")
	c.Response().Header.Set("Content-Length", fmt.Sprintf("%d", len(csvBytes)))
	c.Response().SetBodyRaw(csvBytes)

	return nil
}
