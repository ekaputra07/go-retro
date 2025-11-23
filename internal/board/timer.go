package board

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ekaputra07/go-retro/internal/models"
)

type timerStatus string

const (
	timerStatusRunning timerStatus = "running"
	timerStatusPaused  timerStatus = "paused"
	timerStatusStopped timerStatus = "stopped"
	timerStatusDone    timerStatus = "done"
)

type timerCmd struct {
	cmd   string
	value string
	user  models.User
}

func (tc timerCmd) is(cmd string) bool {
	return tc.cmd == cmd
}

type timer struct {
	Status        timerStatus `json:"status"`
	Display       string      `json:"display"`
	duration      time.Duration
	elapsed       time.Duration
	statusMessage string
	lastCmdUser   models.User
	logger        *slog.Logger
	cmd           chan timerCmd
	state         chan *timer
}

func (t *timer) updateDisplay() {
	rem := t.duration - t.elapsed
	m := int(rem.Minutes()) % 60
	s := int(rem.Seconds()) % 60
	t.Display = fmt.Sprintf("%02d:%02d", m, s)
}

func (t *timer) run(ctx context.Context) {
	t.logger.Info("timer started")
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			t.logger.Info("timer destroyed")
			return
		case <-tick:
			if t.Status == timerStatusRunning {
				t.elapsed += 1 * time.Second

				// done
				if (t.duration - t.elapsed) == 0 {
					t.Status = timerStatusDone
					t.logger.Info("timer done")
				}
				t.updateDisplay()
				t.state <- t
			}
			// reset status message after tick
			t.lastCmdUser = models.User{}
			t.statusMessage = ""

		case cmd := <-t.cmd:
			t.lastCmdUser = cmd.user

			if cmd.is("start") && (t.Status == timerStatusStopped || t.Status == timerStatusDone) {
				d, err := time.ParseDuration(cmd.value)
				if err != nil {
					t.logger.Error("unable to parse timer duration", "err", err)
					continue
				}
				t.duration = d
				t.elapsed = 0
				t.Status = timerStatusRunning
				t.updateDisplay()
				t.statusMessage = fmt.Sprintf("%s started the timer", cmd.user.Name)
				t.state <- t
				t.logger.Info("timer running")

			} else if cmd.is("start") && t.Status == timerStatusPaused {
				t.Status = timerStatusRunning
				t.statusMessage = fmt.Sprintf("%s resumed the timer", cmd.user.Name)
				t.state <- t
				t.logger.Info("timer resumed")

			} else if cmd.is("stop") {
				t.Status = timerStatusStopped
				t.duration = 0
				t.elapsed = 0
				t.statusMessage = fmt.Sprintf("%s stopped the timer", cmd.user.Name)
				t.state <- t
				t.logger.Info("timer stopped")

			} else if cmd.is("pause") {
				t.statusMessage = fmt.Sprintf("%s paused the timer", cmd.user.Name)
				t.Status = timerStatusPaused
				t.state <- t
				t.logger.Info("timer paused")
			}
		}
	}
}

func newTimer(logger *slog.Logger) *timer {
	cmd := make(chan timerCmd)
	state := make(chan *timer)

	return &timer{
		Status:  timerStatusStopped,
		Display: "00:00",
		cmd:     cmd,
		state:   state,
		logger:  logger,
	}
}
