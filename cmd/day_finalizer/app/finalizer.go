package app

import (
	"fmt"
	"log"
	"sync"
)

func (a *app) finalize() error {
	tasks, err := a.tasksClient.GetAllUpdatedTasksForToday()
	if err != nil {
		return fmt.Errorf("get tasks: %s", err)
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

	wg := &sync.WaitGroup{}

	for userID, balance := range tasksPerPopugID {
		wg.Add(1)
		go a.process(wg, userID, balance)
	}

	wg.Wait()

}

func (a *app) process(wg *sync.WaitGroup, userID string, balance int) {
	defer wg.Done()

	err := a.accountingClient.Checkout(userID, balance)
	if err != nil {
		log.Printf("checkout user's %s balance: %s", userID, err)
		return
	}

	email, err := a.usersClient.GetUserEmail(userID)
	if err != nil {
		log.Printf("get user's %s email: %s", userID, err)
		return
	}

	err = a.mailSender.Send(email, balance)
	if err != nil {
		log.Printf("sent email: %s", err)
	}
}
