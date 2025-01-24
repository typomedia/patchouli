package filter

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app/helper"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Operator(c *fiber.Ctx) error {
	id := c.Params("id")
	db := boltdb.New()

	err := db.SetBucket("machine")
	if err != nil {
		return err
	}

	machines, _ := db.GetAllByOperatorId(id, "machine")

	Machines := structs.Machines{}
	for _, v := range machines {
		machine := structs.Machine{}
		err = json.Unmarshal(v, &machine)
		if err != nil {
			return err
		}

		lastUpdate, _ := db.GetLastByName(machine.Id, "history")

		update := structs.Update{}
		err = json.Unmarshal(lastUpdate, &update)
		if err != nil {
			log.Error(err)
		}

		machine.Update = update

		Machines = append(Machines, machine)

	}

	// sort machines by oldest update first
	sort.Sort(structs.ByDate(Machines))

	err = db.SetBucket("config")
	if err != nil {
		return err
	}

	var config structs.Config
	db.Get("main", &config, "config")

	defer db.Close()

	interval := config.General.Interval
	for i := range Machines {
		currentDate := time.Now()

		if Machines[i].Update.Date == "" {
			Machines[i].Update.Date = "0000-00-00"
			Machines[i].Status = "danger"
			continue
		}

		Machines[i].Update.Date = helper.UnixToDateString(Machines[i].Update.Date)

		date, err := time.Parse("2006-01-02", Machines[i].Update.Date)
		if err != nil {
			log.Error(err)
		}

		Machines[i].Days = int(currentDate.Sub(date).Hours() / 24)
		if Machines[i].Days > interval {
			Machines[i].Status = "danger"
		} else if Machines[i].Days > interval/3 {
			Machines[i].Status = "warning"
		} else {
			Machines[i].Status = "success"
		}
		Machines[i].Days = interval - int(currentDate.Sub(date).Hours()/24)

	}
	return c.Render("app/views/dashboard/list", fiber.Map{
		"Machines": Machines,
	})
}
