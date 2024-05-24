package services

import (
	"context"
	"github.com/L1LSunflower/lina/config"
	"log"
	"net/http"
	"time"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/repositories"
	"github.com/L1LSunflower/lina/internal/tools"

	"github.com/L1LSunflower/lina/pkg/db"
	"github.com/L1LSunflower/lina/pkg/telegram"
)

const (
	defaultTimeoutForJob      = 1 * time.Minute
	defaultTimeoutForExporter = 30 * time.Second
)

type Worker struct {
	dbConn     *db.Postgres
	queue      tools.Jobs
	itemsChan  chan *entities.Item
	cancelChan chan struct{}
	tgHandler  telegram.Handler
	usersRepo  repositories.Users
	itemsRepo  repositories.Items
	lastItem   *entities.Item
}

func NewWorker(cfg *config.Config, dbConn *db.Postgres, tHandler telegram.Handler, queue tools.Jobs) *Worker {
	cancelChan := make(chan struct{}, 1)
	return &Worker{
		dbConn:     dbConn,
		queue:      queue,
		itemsChan:  make(chan *entities.Item, cfg.Workers),
		cancelChan: cancelChan,
		tgHandler:  tHandler,
		usersRepo:  repositories.NewUsersRepository(),
		itemsRepo:  repositories.NewItemsRepository(),
		lastItem:   new(entities.Item),
	}
}

func (w *Worker) Start() {
	jobTicker := time.NewTicker(defaultTimeoutForJob)
	dataTicker := time.NewTicker(defaultTimeoutForExporter)
	for {
		select {
		case <-jobTicker.C:
			ctx := tools.CtxWithDepends(context.Background(), w.dbConn)
			w.queue.GoWithTimeoutWithContext(ctx, defaultTimeoutForJob, w.dataUploader)
		case <-dataTicker.C:
			ctx := tools.CtxWithDepends(context.Background(), w.dbConn)
			w.queue.GoWithContext(ctx, w.dataSender)
		case <-w.cancelChan:
			w.Stop()
			break
		}
	}
}

func (w *Worker) Stop() {
	close(w.cancelChan)
	close(w.itemsChan)
}

func (w *Worker) dataUploader(ctx context.Context) {
	log.Println("UPLOADING ITEM START")
	select {
	case item := <-w.itemsChan:
		users, err := w.usersRepo.GetAll(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		var resp *http.Response
		for _, user := range users {
			mg := &entities.MediaGroup{ChatId: user.Id}
			mg.SetImages(item.ImageLinks)
			resp, err = w.tgHandler.MediaGroup(mg) //nolint:bodyclose
			if err != nil {
				log.Println(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("failed to send photos with error, response: %v\n", resp)
			}
			resp, err = w.tgHandler.SendMessage(entities.NewMsg(item, user.Id)) //nolint:bodyclose
			if err != nil {
				log.Println(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("failed to send message with error, response: %v\n", resp)
			}
		}
		if err = w.itemsRepo.UpdateStatus(ctx, item.ID, entities.DoneStatus); err != nil {
			log.Println("ERROR: failed to update status with error: ", err)
		}
		log.Println("UPLOADING ITEM END")
	case <-ctx.Done():
		return
	}
}

func (w *Worker) dataSender(ctx context.Context) {
	log.Println("SENDING ITEM START")
	select {
	case <-ctx.Done():
		return
	default:
		items, err := w.itemsRepo.Items(ctx, w.lastItem.ID, entities.ReadyStatus, 2)
		if err != nil {
			log.Println("ERROR: failed to get items with error: ", err)
		}
		if len(items) > 0 {
			return
		}
		if _, isChannelOpen := <-w.cancelChan; !isChannelOpen {
			return
		}
		for i, item := range items {
			if i == len(items)-1 {
				w.lastItem = item
			}
			w.itemsChan <- item
		}
		log.Println("SENDING ITEM END")
	}
}
