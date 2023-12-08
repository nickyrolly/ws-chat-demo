package postgre

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
)

var (

	// SELECT QUERY
	QuerySelectChatHistory      string = "select sender_user_id, message, reply_time from chat_history where user_id_a = $1 and user_id_b = $2 order by reply_time desc"
	QuerySelectGroupChatHistory string = "select sender_user_id, message, reply_time from group_chat_history where group_id = $1 order by reply_time desc"

	// INSERT QUERY
	QueryInsertChatHistory      string = "insert into chat_history (user_id_a, user_id_b, sender_user_id, message, reply_time) values ($1, $2, $3, $4, $5)"
	QueryInsertGroupChatHistory string = "insert into group_chat_history (group_id, sender_user_id, message, reply_time) values ($1, $2, $3, $4)"
)

func InsertChatHistory(ctx context.Context, params repository.ChatHistoryData) error {
	_, err := DBChat.ExecContext(ctx, QueryInsertChatHistory, params.UserIDA, params.UserIDB, params.SenderUserID, params.Message, params.ReplyTime)
	if err != nil {
		log.Printf("Error insert chat history: %s\n", err.Error())
		return err
	}

	return nil
}

func InsertGroupChatHistory(ctx context.Context, params repository.GroupChatHistoryData) error {
	_, err := DBChat.ExecContext(ctx, QueryInsertGroupChatHistory, params.GroupID, params.SenderUserID, params.Message, params.ReplyTime)
	if err != nil {
		log.Printf("Error insert chat history: %s\n", err.Error())
		return err
	}

	return nil
}
