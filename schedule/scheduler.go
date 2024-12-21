package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/danielelegbe/discord-join-count/bot"
	"github.com/danielelegbe/discord-join-count/service"
	"github.com/danielelegbe/discord-join-count/sqlc"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	Scheduler gocron.Scheduler
	Ctx       context.Context
	Store     *sqlc.Queries
	Bot       *bot.Bot
	Service   *service.Service
}

func New(scheduler gocron.Scheduler, ctx context.Context, store *sqlc.Queries, bot *bot.Bot, service *service.Service) *Scheduler {
	return &Scheduler{
		Scheduler: scheduler,
		Ctx:       ctx,
		Store:     store,
		Bot:       bot,
		Service:   service,
	}
}

func (s *Scheduler) HandleJobs() error {
	interval := 30 * time.Second

	_, err := s.Scheduler.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(
			func() {

				service := service.New(s.Bot, s.Store, s.Ctx)

				userChannel, err := bot.CreateUserChannel(s.Bot.Discord)

				if err != nil {
					panic(err)
				}

				err = service.SendLeaderboardScores(userChannel.ID)

				if err != nil {
					panic(err)
				}
			},
		),
	)

	if err != nil {
		return err
	}

	return err
}

func (s *Scheduler) Start() {
	s.HandleJobs()

	s.Scheduler.Start()

	fmt.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	err := s.Scheduler.Shutdown()

	if err != nil {
		panic(err)
	}
}
