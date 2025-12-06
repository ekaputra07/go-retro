package board

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type timerStatus string

const (
	timerStatusRunning timerStatus = "running"
	timerStatusPaused  timerStatus = "paused"
	timerStatusStopped timerStatus = "stopped"
	timerStatusDone    timerStatus = "done"
)

type timerCmd struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}

func (tc *timerCmd) is(cmd string) bool {
	return tc.Cmd == cmd
}

type timer struct {
	BoardID uuid.UUID   `json:"board_id"`
	Status  timerStatus `json:"status"`
	Display string      `json:"display"`

	nats     *natsutil.NATS
	duration time.Duration
	elapsed  time.Duration
	logger   *slog.Logger
	cmdChan  chan *nats.Msg
	stopChan chan bool
}

func (t *timer) getStateMessage() message {
	return message{
		BoardID: t.BoardID,
		Type:    messageTypeTimerState,
		Data:    t,
	}
}

func (t *timer) getNotificationMessage(msg string, user models.User) message {
	return message{
		BoardID: t.BoardID,
		Type:    messageTypeBoardNotification,
		Data:    msg,
		User:    user,
	}
}

// parseCmd parses timer command from message instance
func (t *timer) parseCmd(msg message) (*timerCmd, error) {
	var command string
	var value string
	err := msg.stringVar(&command, "cmd")
	if err != nil {
		return nil, err
	}
	cmd := timerCmd{Cmd: command}

	err = msg.stringVar(&value, "value")
	if err == nil {
		cmd.Value = value
	}
	return &cmd, nil
}

// broadcast publish timer related message to board's broadcast channel/topic
func (t *timer) broadcast(topic string, msgs ...message) error {
	if len(msgs) == 0 {
		return fmt.Errorf("no message to broadcast")
	}

	var msgJson []byte
	if len(msgs) > 1 {
		m := newMessageList(t.BoardID, msgs...)
		msg, err := m.encode()
		if err != nil {
			return fmt.Errorf("failed to encode timer messageList: %s", err.Error())
		}
		msgJson = msg
	} else {
		msg, err := msgs[0].encode()
		if err != nil {
			return fmt.Errorf("failed to encode timer message: %s", err.Error())
		}
		msgJson = msg
	}
	return t.nats.Conn.Publish(topic, msgJson)
}

func (t *timer) updateDisplay() {
	rem := t.duration - t.elapsed
	m := int(rem.Minutes()) % 60
	s := int(rem.Seconds()) % 60
	t.Display = fmt.Sprintf("%02d:%02d", m, s)
}

// run starts the timer process.
// When started, it subscribe to timer command topic and react when new command received.
func (t *timer) run() {
	cmdSub, err := t.nats.Conn.ChanSubscribe(timerCmdTopic(t.BoardID), t.cmdChan)
	if err != nil {
		t.logger.Error("timer failed to subscribe cmd topic", "id", t.BoardID, "err", err.Error())
		return
	}

	tick := time.NewTicker(1 * time.Second)

	defer func() {
		cmdSub.Unsubscribe()
		close(t.cmdChan)
		tick.Stop()
		t.logger.Info("timer stopped")
	}()

	t.logger.Info("timer started")
	for {
		select {
		case <-t.stopChan:
			t.logger.Info("timer stop signal received")
			return

		case <-tick.C:
			if t.Status == timerStatusRunning {
				t.elapsed += 1 * time.Second

				// done
				if (t.duration - t.elapsed) == 0 {
					t.Status = timerStatusDone
					t.logger.Info("timer done")
				}
				t.updateDisplay()
				t.broadcast(broadcastMessageTopic(t.BoardID), t.getStateMessage())
			}

		case msg := <-t.cmdChan:
			var m message
			err := json.Unmarshal(msg.Data, &m)
			if err != nil {
				t.logger.Error("failed decode timer cmd message", "err", err.Error())
				continue
			}
			cmd, err := t.parseCmd(m)
			if err != nil {
				t.logger.Error("failed parsing timer cmd", "err", err.Error())
				continue
			}
			if err = t.handleCommand(msg, cmd, m.User); err != nil {
				t.logger.Error("failed handling timer cmd", "err", err.Error())
				continue
			}
		}
	}
}

// handleCommand handles each command
func (t *timer) handleCommand(nmsg *nats.Msg, cmd *timerCmd, user models.User) error {
	switch {
	case cmd.is("status"):
		stateMsg := t.getStateMessage()
		b, err := stateMsg.encode()
		if err != nil {
			return fmt.Errorf("failed to encode timer state: %s", err.Error())
		}
		if err = nmsg.Respond(b); err != nil {
			return fmt.Errorf("failed responding to timer status: %s", err.Error())
		}

	case cmd.is("start") && (t.Status == timerStatusStopped || t.Status == timerStatusDone):
		d, err := time.ParseDuration(cmd.Value)
		if err != nil {
			return fmt.Errorf("unable to parse timer duration: %s", err.Error())
		}
		t.duration = d
		t.elapsed = 0
		t.Status = timerStatusRunning
		t.updateDisplay()

		statusMessage := fmt.Sprintf("%s started the timer", user.Name)
		t.logger.Info("timer running")
		return t.broadcast(
			broadcastMessageTopic(t.BoardID),
			t.getStateMessage(),
			t.getNotificationMessage(statusMessage, user),
		)

	case cmd.is("start") && t.Status == timerStatusPaused:
		t.Status = timerStatusRunning

		statusMessage := fmt.Sprintf("%s resumed the timer", user.Name)
		t.logger.Info("timer resumed")
		return t.broadcast(
			broadcastMessageTopic(t.BoardID),
			t.getStateMessage(),
			t.getNotificationMessage(statusMessage, user),
		)

	case cmd.is("stop"):
		t.Status = timerStatusStopped
		t.duration = 0
		t.elapsed = 0

		statusMessage := fmt.Sprintf("%s stopped the timer", user.Name)
		t.logger.Info("timer stopped")
		return t.broadcast(
			broadcastMessageTopic(t.BoardID),
			t.getStateMessage(),
			t.getNotificationMessage(statusMessage, user),
		)

	case cmd.is("pause"):
		t.Status = timerStatusPaused

		statusMessage := fmt.Sprintf("%s paused the timer", user.Name)
		t.logger.Info("timer paused")
		return t.broadcast(
			broadcastMessageTopic(t.BoardID),
			t.getStateMessage(),
			t.getNotificationMessage(statusMessage, user),
		)
	}
	return nil
}

func newTimer(boardID uuid.UUID, nats_ *natsutil.NATS, logger *slog.Logger) *timer {
	return &timer{
		BoardID:  boardID,
		Status:   timerStatusStopped,
		Display:  "00:00",
		nats:     nats_,
		logger:   logger,
		cmdChan:  make(chan *nats.Msg, 256),
		stopChan: make(chan bool),
	}
}
