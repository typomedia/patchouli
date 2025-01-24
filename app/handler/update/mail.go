package update

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"

	"github.com/XotoX1337/tinymail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	patchouli "github.com/typomedia/patchouli/app"
	"github.com/typomedia/patchouli/app/store/boltdb"
	"github.com/typomedia/patchouli/app/structs"
)

func Mail(c *fiber.Ctx) error {
	updateId := c.Params("id")

	db := boltdb.New()
	defer db.Close()
	db.SetBucket("history")

	var update structs.Update
	err := db.Get(updateId, &update, "history")
	if err != nil {
		log.Error(err)
	}

	var machine structs.Machine
	err = db.Get(update.Machine, &machine, "machine")
	if err != nil {
		log.Error(err)
	}

	var operator structs.Operator
	err = db.Get(update.Operator.Id, &operator, "operator")
	if err != nil {
		log.Error(err)
	}

	var config structs.Config
	err = db.Get("main", &config, "config")
	if err != nil {
		log.Error(err)
	}

	opts := tinymail.MailerOpts{
		User:     config.Smtp.Username,
		Password: config.Smtp.Password,
		Host:     config.Smtp.Host,
		Port:     config.Smtp.Port,
	}

	mailer, err := tinymail.New(opts)
	if err != nil {
		return err
	}

	resp, err := http.Get(c.BaseURL() + "/html/mail/update.html")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	tplData := map[string]interface{}{
		"Update":   update,
		"Machine":  machine,
		"Operator": operator,
		"Version":  patchouli.GetApp().Version,
		"App":      patchouli.GetApp().Name,
	}

	msg, err := tinymail.FromTemplateString(tplData, string(body))
	if err != nil {
		log.Error(err)
	}

	msg.SetFrom(config.Smtp.Sender)
	msg.SetTo(config.General.Email)
	msg.SetSubject(fmt.Sprintf("Patchmgmt: %s %s - %s", machine.System.Name, machine.Name, machine.Fqdn))

	err = mailer.SetMessage(msg).Send()
	if err != nil {
		log.Error(err)
		update.Mail = false

	}
	update.Mail = true
	err = db.Set(update.Id, update, "history")
	if err != nil {
		log.Error(err)
	}
	return c.Redirect("/machine/update/list/" + machine.Id)
}
