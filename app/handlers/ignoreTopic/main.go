package ignoretopic

import (
	"errors"
	"grimoire/database"
	"grimoire/database/models"
	"grimoire/handlers"
	"time"

	"github.com/go-pg/pg/v10"
	tele "gopkg.in/telebot.v4"
)

func connectToDatabase(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	db := database.GetDB()
	newArgs := make(handlers.Arg)
	newArgs["db"] = db
	return &newArgs, nil
}

func markTopicAsIgnored(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	db := (*args)["db"].(*pg.DB)
	message := c.Message()
	if message == nil {
		return nil, errors.New("message is nil")
	}
	threadID := message.ThreadID
	if threadID == 0 {
		return nil, errors.New("no thread ID found in message")
	}

	_, err := db.Model(&models.IgnoredTopic{ChatID: message.Chat.ID, ThreadID: threadID}).OnConflict("DO NOTHING").SelectOrInsert()
	if err != nil {
		return nil, errors.New("error marking topic as ignored: " + err.Error())
	}

	return args, nil
}

func IgnoreTopicChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(10*time.Second, connectToDatabase, markTopicAsIgnored)
}