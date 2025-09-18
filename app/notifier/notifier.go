package notifier

import (
	"github.com/XotoX1337/tinymail"
	"github.com/gofiber/fiber/v2/log"
	"github.com/typomedia/patchouli/app"

	"github.com/typomedia/patchouli/app/encryption"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

const NOTIFY_DAYS = 7

type Notifier struct {
	Operators map[string]structs.Machines
}

func (n Notifier) Run() {
	db := boltdb.New()
	defer db.Close()

	config := app.GetApp().Config

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

	smtpPasswd, err := encryption.DecryptString(config.Smtp.Password)
	if err != nil {
		return
	}
	opts := tinymail.MailerOpts{
		User:     config.Smtp.Username,
		Password: smtpPasswd,
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
			"Hostname": config.General.Hostname,
			"Version":  app.GetApp().Version,
			"App":      app.GetApp().Name,
		}

		msg, err := tinymail.FromTemplateString(tplData, app.GetApp().NotifyTemplate)
		if err != nil {
			log.Error(err)
		}

		msg.SetFrom(config.Smtp.Sender)
		if operator.Email == "" {
			msg.SetTo(config.General.Email)
		} else {
			msg.SetTo(operator.Email)
		}

		msg.SetSubject("Patchmgmt: Maschinen benötigen Updates!")

		mailer.SetMessage(msg)

		err = mailer.SetMessage(msg).Send()
		if err != nil {
			log.Fatal(err)
		}
	}
}
