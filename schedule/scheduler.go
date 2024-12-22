package schedule

import (
	"context"
	"log/slog"

	"github.com/danielelegbe/discord-join-count/bot"
	"github.com/danielelegbe/discord-join-count/config"
	"github.com/danielelegbe/discord-join-count/sqlc"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	Scheduler gocron.Scheduler
	Ctx       context.Context
	Store     *sqlc.Queries
	Bot       *bot.Bot
}

func New(scheduler gocron.Scheduler, ctx context.Context, store *sqlc.Queries, bot *bot.Bot) *Scheduler {
	return &Scheduler{
		Scheduler: scheduler,
		Ctx:       ctx,
		Store:     store,
		Bot:       bot,
	}
}

func (s *Scheduler) HandleJobs() error {
	everySundayCron := "0 21 * * 0"

	_, err := s.Scheduler.NewJob(
		gocron.CronJob(everySundayCron, false),
		gocron.NewTask(
			func() {
				err := s.Bot.SendWeeklyLeaderboardScores(config.ConfigInstance.SPOST_CHANNEL_ID)

				if err != nil {
					slog.Error(err.Error())
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

	slog.Info("Scheduler started")
}

func (s *Scheduler) Stop() {
	err := s.Scheduler.Shutdown()

	if err != nil {
		slog.Error(err.Error())
	}
}
