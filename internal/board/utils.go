package board

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func timerCmdTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.timer.cmd", boardID)
}

func broadcastMessageTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.msg.out", boardID)
}

func queryTimerStatus(conn *nats.Conn, boardID uuid.UUID) (*nats.Msg, error) {
	cmdMsg := message{
		Type: messageTypeTimerCmd,
		Data: timerCmd{Cmd: "status"},
	}
	cmd, err := cmdMsg.encode()
	if err != nil {
		return nil, err
	}
	return conn.Request(timerCmdTopic(boardID), cmd, 1*time.Second)
}
