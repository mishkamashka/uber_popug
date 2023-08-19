package app

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
)

func (a *app) finalize() error {
	tasks, err := a.tasksClient.GetAllUpdatedTasksForToday()
	if err != nil {
		return fmt.Errorf("get tasks: %s", err)
	}

	if len(tasks) == 0 {
		return nil
	}

	tasksPerPopugID := map[string]int{}

	for _, task := range tasks {
		var value int
		switch task.Status {
		case "closed":
			value = int(task.PriceForClosing)
		case "open":
			value = -int(task.PriceForClosing)
		default:
			log.Printf("unknown task status: %s", task.Status)
			continue
		}

		tasksPerPopugID[task.AssigneeId] += value
	}

	wg := errgroup.Group{}

	for userID, balance := range tasksPerPopugID {
		if balance > 0 {
			wg.Go(func() error {
				return a.process(userID, balance)
			})
		}
	}

	if err := wg.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *app) process(userID string, balance int) error {
	err := a.accountingClient.Checkout(userID, balance)
	if err != nil {
		msg := fmt.Sprintf("checkout user's %s balance: %s", userID, err)
		log.Println(msg)

		return errors.New(msg)
	}

	email, err := a.usersClient.GetUserEmail(userID)
	if err != nil {
		msg := fmt.Sprintf("get user's %s email: %s", userID, err)
		log.Println(msg)

		return errors.New(msg)
	}

	err = a.mailSender.Send(email, balance)
	if err != nil {
		msg := fmt.Sprintf("sent email: %s", err)
		log.Println(msg)

		return errors.New(msg)
	}

	return nil
}
