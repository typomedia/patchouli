package system

import (
	"cmp"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"slices"
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
		machinesOfSystem, err := db.GetMachinesBySystem(system.Id)
		if err != nil {
			return err
		}
		system.MachineCount = len(machinesOfSystem)
		Systems = append(Systems, system)
	}
	slices.SortFunc(Systems, func(a, b structs.System) int {
		return cmp.Compare(b.MachineCount, a.MachineCount)
	})
	defer db.Close()

	return c.Render("app/views/system/list", fiber.Map{
		"Systems": Systems,
	})
}
