package board

import (
	"fmt"
	"log/slog"
	"time"
)

type timerStatus string

const (
	timerStatusRunning timerStatus = "running"
	timerStatusPaused  timerStatus = "paused"
	timerStatusStopped timerStatus = "stopped"
	timerStatusDone    timerStatus = "done"
)

type timerCmd struct {
	cmd    string
	value  string
	client *Client
}

func (tc timerCmd) is(cmd string) bool {
	return tc.cmd == cmd
}

type timer struct {
	Status            timerStatus `json:"status"`
	Display           string      `json:"display"`
	duration          time.Duration
	elapsed           time.Duration
	cmd               chan timerCmd
	state             chan *timer
	statusMessage     string
	lastCommandClient *Client
	logger            *slog.Logger
}

func (t *timer) updateDisplay() {
	rem := t.duration - t.elapsed
	m := int(rem.Minutes()) % 60
	s := int(rem.Seconds()) % 60
	t.Display = fmt.Sprintf("%02d:%02d", m, s)
}

func (t *timer) run() {
	t.logger.Info("timer started")
	tick := time.Tick(1 * time.Second)
	for {
		select {
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
			t.lastCommandClient = nil
			t.statusMessage = ""

		case cmd := <-t.cmd:
			t.lastCommandClient = cmd.client

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
				t.statusMessage = fmt.Sprintf("%s started the timer", cmd.client.User.Name)
				t.state <- t
				t.logger.Info("timer running")

			} else if cmd.is("start") && t.Status == timerStatusPaused {
				t.Status = timerStatusRunning
				t.statusMessage = fmt.Sprintf("%s resumed the timer", cmd.client.User.Name)
				t.state <- t
				t.logger.Info("timer resumed")

			} else if cmd.is("stop") {
				t.Status = timerStatusStopped
				t.duration = 0
				t.elapsed = 0
				t.statusMessage = fmt.Sprintf("%s stoped the timer", cmd.client.User.Name)
				t.state <- t
				t.logger.Info("timer stopped")

			} else if cmd.is("pause") {
				t.statusMessage = fmt.Sprintf("%s paused the timer", cmd.client.User.Name)
				t.Status = timerStatusPaused
				t.state <- t
				t.logger.Info("timer paused")

			} else if cmd.is("destroy") {
				t.logger.Info("timer destroyed")
				return
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
