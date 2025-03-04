package notifier

import (
	"github.com/XotoX1337/tinymail"
	"github.com/gofiber/fiber/v2/log"
	patchouli "github.com/typomedia/patchouli/app"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

const NOTIFY_DAYS = 7

type Notifier struct {
	Hostname  string
	Operators map[string]structs.Machines
}

func (n Notifier) Run() {
	db := boltdb.New()
	defer db.Close()

	config, err := db.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	machines, err := db.GetActiveMachines()
	if err != nil {
		log.Fatal(err)
	}
	n.Operators = make(map[string]structs.Machines)
	for _, machine := range machines {
		if machine.Days > NOTIFY_DAYS {
			continue
		}
		n.Operators[machine.Operator.Id] = append(n.Operators[machine.Operator.Id], machine)
	}

	opts := tinymail.MailerOpts{
		User:     config.Smtp.Username,
		Password: config.Smtp.Password,
		Host:     config.Smtp.Host,
		Port:     config.Smtp.Port,
	}

	mailer, err := tinymail.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	for operatorId, machines := range n.Operators {

		operator, err := db.GetOperatorById(operatorId)
		if err != nil {
			log.Fatal(err)
		}

		tplData := map[string]interface{}{
			"Machines": machines,
			"Operator": operator,
			"Host":     n.Hostname,
			"Version":  patchouli.GetApp().Version,
			"App":      patchouli.GetApp().Name,
		}

		msg, err := tinymail.FromTemplateString(tplData, patchouli.GetApp().NotifyTemplate)
		if err != nil {
			log.Error(err)
		}

		msg.SetFrom(config.Smtp.Sender)
		msg.SetTo(config.General.Email)
		msg.SetSubject("Patchmgmt: Maschinen benötigen Updates!")

		mailer.SetMessage(msg)

		err = mailer.SetMessage(msg).Send()
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}
