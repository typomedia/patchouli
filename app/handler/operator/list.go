package operator

import (
	"cmp"
	"github.com/gofiber/fiber/v2"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
	"slices"
	"strings"
)

func List(c *fiber.Ctx) error {
	db := boltdb.New()

	operators, err := db.GetAllOperators()
	var Operators structs.Operators
	if err != nil {
		return err
	}
	for _, operator := range operators {
		operator.MachineCount = 0
		machinesOfOperator, err := db.GetMachinesByOperator(operator.Id)
		if err != nil {
			return err
		}
		operator.MachineCount = len(machinesOfOperator)
		Operators = append(Operators, operator)
	}

	slices.SortFunc(Operators, func(a, b structs.Operator) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	defer db.Close()

	return c.Render("app/views/operator/list", fiber.Map{
		"Operators": Operators,
	})
}
