package board

// messageType represents the type of message that can be sent to and from the client
type messageType string

const (
	messageTypeBoardUsers               messageType = "board.users"
	messageTypeBoardStatus              messageType = "board.status"
	messageTypeNotificationNotification messageType = "board.notification"
	messageTypeColumnNew                messageType = "column.new"
	messageTypeColumnUpdate             messageType = "column.update"
	messageTypeColumnDelete             messageType = "column.delete"
	messageTypeCardNew                  messageType = "card.new"
	messageTypeCardUpdate               messageType = "card.update"
	messageTypeCardDelete               messageType = "card.delete"
	messageTypeCardVote                 messageType = "card.vote"
	messageTypeTimerCmd                 messageType = "timer.cmd"
	messageTypeTimerState               messageType = "timer.state"
)

// message represents a message that can be sent to and from the client
type message struct {
	Type messageType `json:"type"`
	Data any         `json:"data"`

	// for inbound message, client is assigned when the message is received
	// for outbound message, client is used to exclude the sender from receiving the message
	client *Client
}
