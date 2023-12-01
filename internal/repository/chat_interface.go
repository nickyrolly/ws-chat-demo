package repository

import "time"

type ChatHistoryData struct {
	UserIDA      int
	UserIDB      int
	SenderUserID int
	Message      string
	ReplyTime    time.Time
}

type GroupChatHistoryData struct {
	GroupID      int
	SenderUserID int
	Message      string
	ReplyTime    time.Time
}
