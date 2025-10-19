package registertext

import (
	"errors"
	"time"

	"grimoire/database"
	"grimoire/database/models"
	"grimoire/handlers"
	"grimoire/util"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v4"
)

func connectToDatabase(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	db := database.GetDB()

	newArgs := make(handlers.Arg)
	newArgs["db"] = db
	
	return &newArgs, nil
}

func registerText(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	text := c.Message().Text
	if text == "" {
		return nil, errors.New("text is empty")
	}

	db := (*args)["db"].(*pg.DB)
	
	textDb := models.Text{
		Text: text,
		UserID: c.Sender().ID,
	}
	_, err := db.Model(&textDb).Insert()
	if err != nil {
		return nil, errors.New("error inserting text: " + err.Error())
	}

	_, err = db.Model(&models.User{TgID: c.Sender().ID}).OnConflict("DO NOTHING").SelectOrInsert()
	if err != nil {
		return nil, errors.New("error updating user: " + err.Error())
	}

	args = util.UpdateArgs(args, "text", text)

	return args, nil
}

func resendText(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	_, err := c.Bot().Send(
		&tele.Chat{ID: viper.GetInt64("bloodofspring.archive_group_id")},
		(*args)["text"],
		&tele.SendOptions{ThreadID: viper.GetInt("bloodofspring.thread_id")},
	)
	if err != nil {
		return nil, errors.New("error sending text: " + err.Error())
	}
	return args, nil
}

func sendResultMessage(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	err := c.Send("Сохранено!")
	if err != nil {
		return nil, errors.New("error sending result message: " + err.Error())
	}

	return args, nil
}


func RegisterTextChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(10 * time.Second, connectToDatabase, registerText, resendText, sendResultMessage)
}
