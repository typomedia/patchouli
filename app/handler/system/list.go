package system

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	err := db.SetBucket("systems")
	if err != nil {
		return err
	}

	systems, _ := db.GetAll("systems")

	Systems := structs.Systems{}
	for _, v := range systems {
		system := structs.System{}
		err = json.Unmarshal(v, &system)
		if err != nil {
			return err
		}
		Systems = append(Systems, system)
	}

	defer db.Close()

	return c.Render("app/views/system/list", fiber.Map{
		"Systems": Systems,
	})
}
