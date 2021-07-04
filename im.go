package webmail

import "encoding/json"

type MessageId int

type ConversationId int

type ContactId string

type Status string

const (
	available Status = "available"
	offline   Status = "offline"
	dnd       Status = "dnd"
	away      Status = "away"
	invisible Status = "invisible"
)

type Presence struct {
	ContactId ContactId   `json:"contactId"`
	Status    Status      `json:"status"`
	Text      string      `json:"text"`
	Date      UtcDateTime `json:"date"`
}

type PresenceList []Presence

type ConversationEvent string

const (
	Created   ConversationEvent = "Created"   // when conversation id is created
	Updated   ConversationEvent = "Updated"   // when message flows into conversation
	Delivered ConversationEvent = "Delivered" // when messages within this conversation were delivered
	Read      ConversationEvent = "Read"      // when messages within this conversation has been read
)

type ContactIdList []ContactId

type Conversation struct {
	ConversationId      ConversationId    `json:"conversationId"`
	LastActivity        UtcDateTime       `json:"lastActivity"`        // used to find RECENT conversations
	SentLastDeliveredId MessageId         `json:"sentLastDeliveredId"` // all higher ids are undelivered (user-specific - messages from me)
	SentLastReadId      MessageId         `json:"sentLastReadId"`      // // all higher ids are undelivered (user-specific - messages from me)
	ReceivedLastReadId  MessageId         `json:"receivedLastReadId"`  // all higher ids are undelivered (user-specific - message to me)
	ReceivedUnreadCount int               `json:"receivedUnreadCount"` // user-specific - message to me
	Contacts            ContactIdList     `json:"contacts"`            // one for 1:1, more for groupchats
	Event               ConversationEvent `json:"event"`               // to notify other users/sessions about conversation activity
	Muted               bool              `json:"muted"`
}

type ConversationList []Conversation

type MessageEvent string

const (
	active   MessageEvent = "active"   // typing
	inactive MessageEvent = "inactive" // stopped typing
)

type Message struct {
	MessageId MessageId      `json:"messageId"`
	Text      string         `json:"text"`
	Event     MessageEvent   `json:"event"`
	To        ConversationId `json:"to"` // conversation ID
	From      ContactId      `json:"from"`
	Time      UtcDateTime    `json:"time"` // filled by server
}

type MessageList []Message

// imGetPresence - Retrieve presence for users.
//	contacts - (optional) list of all contacts
// Return
//	list - list of statuses of given contacts (or all non-offilne contacts if input array is empty)
func (c *ClientConnection) imGetPresence(contacts KIdList) (PresenceList, error) {
	params := struct {
		Contacts KIdList `json:"contacts"`
	}{contacts}
	data, err := c.CallRaw("im.getPresence", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List PresenceList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// imSubscribePresence - Note that callback will be executed periodically as changes from server arrive.
// Return
//	list - Presence status of all users. Offline users are always excluded. Missing means offline. (TODO)
func (c *ClientConnection) imSubscribePresence() (PresenceList, error) {
	data, err := c.CallRaw("im.subscribePresence", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List PresenceList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// imSetPresence - Update own presence. Server than resend such status to all interested clients (subscribed) as a Presence with current date.
//	status - new user status to be set (online, offline, ...)
//	text - status text
func (c *ClientConnection) imSetPresence(status Status, text string) error {
	params := struct {
		Status Status `json:"status"`
		Text   string `json:"text"`
	}{status, text}
	_, err := c.CallRaw("im.setPresence", params)
	return err
}

// imCreateConversation - Create new conversation (or return existing one)
//	contacts - required conversation, one for 1:1, more for groupchats
// Return
//	conversation - created conversation
func (c *ClientConnection) imCreateConversation(contacts ContactIdList) (*Conversation, error) {
	params := struct {
		Contacts ContactIdList `json:"contacts"`
	}{contacts}
	data, err := c.CallRaw("im.createConversation", params)
	if err != nil {
		return nil, err
	}
	conversation := struct {
		Result struct {
			Conversation Conversation `json:"conversation"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &conversation)
	return &conversation.Result.Conversation, err
}

// imSubscribeConversations - It provides list of all conversations in which the current user participates and also subscribes for further changes.
// Return
//	list - all conversations in which current user participes
func (c *ClientConnection) imSubscribeConversations() (ConversationList, error) {
	data, err := c.CallRaw("im.subscribeConversations", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List ConversationList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// imMuteConversation - Set conversation as (un)muted for the current user
//	conversationId - conversation to be set
func (c *ClientConnection) imMuteConversation(conversationId ConversationId, mute bool) error {
	params := struct {
		ConversationId ConversationId `json:"conversationId"`
		Mute           bool           `json:"mute"`
	}{conversationId, mute}
	_, err := c.CallRaw("im.muteConversation", params)
	return err
}

// imReadConversation - Set last read message in conversation for the current user
//	conversationId - conversation to be set
func (c *ClientConnection) imReadConversation(conversationId ConversationId, lastReadId MessageId) error {
	params := struct {
		ConversationId ConversationId `json:"conversationId"`
		LastReadId     MessageId      `json:"lastReadId"`
	}{conversationId, lastReadId}
	_, err := c.CallRaw("im.readConversation", params)
	return err
}

// imGetMessages - Returns ordered list of messages from single conversation. The list is always ordered by messageId and always returns older messages than currentMessageId parameter
//	conversationId - Identifier of a conversation.
//	currentMessageId - Messages older than currentMessageId are returned
//	count - Required messages count
// Return
//	list - Ordered list of messages for given conversation
func (c *ClientConnection) imGetMessages(conversationId ConversationId, currentMessageId MessageId, count int) (MessageList, error) {
	params := struct {
		ConversationId   ConversationId `json:"conversationId"`
		CurrentMessageId MessageId      `json:"currentMessageId"`
		Count            int            `json:"count"`
	}{conversationId, currentMessageId, count}
	data, err := c.CallRaw("im.getMessages", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List MessageList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// imSubscribeMessages - ordered by messageId. The continuous flow of messages includes all new, it is NOT filtered by query.
//	conversationId - Identifier of a conversation
//	currentMessageId - All newer (included currentMessageId) messages are returned + count of older messages
//	count - Required older messages count
func (c *ClientConnection) imSubscribeMessages(conversationId ConversationId, currentMessageId int, count int) (MessageList, error) {
	params := struct {
		ConversationId   ConversationId `json:"conversationId"`
		CurrentMessageId int            `json:"currentMessageId"`
		Count            int            `json:"count"`
	}{conversationId, currentMessageId, count}
	data, err := c.CallRaw("im.subscribeMessages", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List MessageList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// imUnsubscribeMessages - Stops listening on new messages within conversation.
//	conversationId - Identifier of conversation to unsubscribe.
func (c *ClientConnection) imUnsubscribeMessages(conversationId ConversationId) error {
	params := struct {
		ConversationId ConversationId `json:"conversationId"`
	}{conversationId}
	_, err := c.CallRaw("im.unsubscribeMessages", params)
	return err
}

// imSendMessage - It sends a message into a conversation. Field 'message.to' must be known before sending a message. Client either have it or must asks for it.
//	message - Message to be send. It already contains the destination (the 'to' field). It does not contain a 'messageId' as it is generated by server.
//	markAsRead - Message will be marked as read for to users.
// Return
//	messageId - Message id
//	time - Time of message
func (c *ClientConnection) imSendMessage(message Message, markAsRead bool) (*MessageId, *UtcDateTime, error) {
	params := struct {
		Message    Message `json:"message"`
		MarkAsRead bool    `json:"markAsRead"`
	}{message, markAsRead}
	data, err := c.CallRaw("im.sendMessage", params)
	if err != nil {
		return nil, nil, err
	}
	messageId := struct {
		Result struct {
			MessageId MessageId   `json:"messageId"`
			Time      UtcDateTime `json:"time"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &messageId)
	return &messageId.Result.MessageId, &messageId.Result.Time, err
}
