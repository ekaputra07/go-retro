package board

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTopics(t *testing.T) {
	id := uuid.New()
	assert.Equal(t, fmt.Sprintf("boards.%s.status", id), boardStatusTopic(id))
	assert.Equal(t, fmt.Sprintf("boards.%s.client-joined", id), clientJoinTopic(id))
	assert.Equal(t, fmt.Sprintf("boards.%s.client-leave", id), clientLeaveTopic(id))
	assert.Equal(t, fmt.Sprintf("boards.%s.msg-in", id), inboundMessageTopic(id))
	assert.Equal(t, fmt.Sprintf("boards.%s.msg-out", id), broadcastMessageTopic(id))
}
