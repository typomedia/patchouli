package operator

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	err := db.SetBucket("operator")
	if err != nil {
		return err
	}

	operators, _ := db.GetAll("operator")

	Operators := structs.Operators{}
	for _, v := range operators {
		operator := structs.Operator{}
		err = json.Unmarshal(v, &operator)
		if err != nil {
			return err
		}
		Operators = append(Operators, operator)

	}

	defer db.Close()

	return c.Render("app/views/operator/list", fiber.Map{
		"Operators": Operators,
	})
}
