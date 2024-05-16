package services

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/repositories"
	"github.com/L1LSunflower/lina/pkg/telegram"
)

type Worker struct {
	workerPool *sync.WaitGroup
	itemsChan  chan *entities.Item
	cancelChan chan struct{}
	tgHandler  telegram.Handler
	usersRepo  repositories.Users
	itemsRepo  repositories.Items
}

func NewWorker() *Worker {
	return &Worker{
		workerPool: &sync.WaitGroup{},
		itemsChan:  nil,
		cancelChan: nil,
		tgHandler:  nil,
		usersRepo:  repositories.NewUsersRepository(),
		itemsRepo:  repositories.NewItemsRepository(),
	}
}

func (w *Worker) Start() {

}

func (w *Worker) Stop() {

}

func (w *Worker) dataUploader(ctx context.Context) {
	item := <-w.itemsChan
	users, err := w.usersRepo.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	var resp *http.Response
	for _, user := range users {
		mg := &entities.MediaGroup{ChatId: user.Id}
		mg.SetImages(item.ImageLinks)
		resp, err = w.tgHandler.MediaGroup(mg)
		if err != nil {
			log.Println(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to send photos with error, response: %v\n", resp)
		}
		resp, err = w.tgHandler.SendMessage(entities.NewMsg(item, user.Id))
		if err != nil {
			log.Println(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to send message with error, response: %v\n", resp)
		}
	}
}

func (w *Worker) dataSender() {

}
