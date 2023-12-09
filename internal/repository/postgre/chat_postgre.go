package postgre

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
	"github.com/tokopedia/sqlt"
)

var (
	DBChat   *sqlt.DB
	DBMaster string = "user=secure_random_username password=secure_random_password dbname=chat host=postgres port=5432 sslmode=disable"

	// SELECT QUERY
	QuerySelectChatHistory      string = "select sender_user_id, message, reply_time from chat_history where user_id_a = $1 and user_id_b = $2 order by reply_time desc"
	QuerySelectGroupChatHistory string = "select sender_user_id, message, reply_time from group_chat_history where group_id = $1 order by reply_time desc"

	// INSERT QUERY
	QueryInsertChatHistory      string = "insert into chat_history (user_id_a, user_id_b, sender_user_id, message, reply_time) values ($1, $2, $3, $4, $5)"
	QueryInsertGroupChatHistory string = "insert into group_chat_history (group_id, sender_user_id, message, reply_time) values ($1, $2, $3, $4)"
)

func InitPostgreSqltDB(connMaster string, connSlave string) error {
	var (
		dbURL string
		err   error
	)

	if connSlave == "" {
		dbURL = connMaster
	} else {
		dbURL = fmt.Sprintf("%s;%s", connMaster, connSlave)
	}

	DBChat, err = sqlt.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	if err := DBChat.Ping(); err != nil {
		return err
	}
	DBChat.SetConnMaxLifetime(9 * time.Second)

	return nil
}

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
