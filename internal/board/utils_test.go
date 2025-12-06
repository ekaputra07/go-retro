package board

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTopics(t *testing.T) {
	id := uuid.New()
	assert.Equal(t, fmt.Sprintf("boards.%s.timer.cmd", id), timerCmdTopic(id))
	assert.Equal(t, fmt.Sprintf("boards.%s.msg.out", id), broadcastMessageTopic(id))
}
