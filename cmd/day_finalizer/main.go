package main

import (
	"log"

	"uber-popug/cmd/day_finalizer/accounting_client"
	"uber-popug/cmd/day_finalizer/app"
	"uber-popug/cmd/day_finalizer/mail_sender"
	"uber-popug/cmd/day_finalizer/tasks_client"
	"uber-popug/cmd/day_finalizer/users_client"
)

func main() {
	accountingClient := accounting_client.New()
	usersClient := users_client.New()
	tasksClient := tasks_client.New()

	mailSender := mail_sender.NewSender()

	app := app.NewApp(accountingClient, tasksClient, usersClient, mailSender)

	if err := app.FinalizeDay(); err != nil {
		log.Fatal(err)
	}
}
