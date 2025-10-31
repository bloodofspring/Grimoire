package registertext

import (
	"errors"
	"time"

	"grimoire/database"
	"grimoire/database/models"
	"grimoire/handlers"

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
	message := c.Message()
	if message == nil {
		return nil, errors.New("message is nil")
	}

	db := (*args)["db"].(*pg.DB)

	// Извлекаем текст сообщения или подпись к медиа
	text := message.Text
	if text == "" {
		text = message.Caption
	}

	// Сохраняем пользователя если его еще нет
	_, err := db.Model(&models.User{TgID: c.Sender().ID}).OnConflict("DO NOTHING").SelectOrInsert()
	if err != nil {
		return nil, errors.New("error updating user: " + err.Error())
	}

	// Сохраняем текст в БД только если он есть
	if text != "" {
		textDb := models.Text{
			Text:   text,
			UserID: c.Sender().ID,
		}
		_, err = db.Model(&textDb).Insert()
		if err != nil {
			return nil, errors.New("error inserting text: " + err.Error())
		}
	}

	// Сохраняем сообщение в args для пересылки (даже если текста нет)
	(*args)["message"] = message
	if text != "" {
		(*args)["text"] = text
	}

	return args, nil
}

func forwardToChannel(c tele.Context, args *handlers.Arg) (*handlers.Arg, error) {
	messageRaw, ok := (*args)["message"]
	if !ok {
		return nil, errors.New("message not found in args")
	}
	message, ok := messageRaw.(*tele.Message)
	if !ok || message == nil {
		return nil, errors.New("message is nil or invalid type in args")
	}

	targetChannelID := viper.GetInt64("bloodofspring.target_channel_id")
	if targetChannelID == 0 {
		return nil, errors.New("target_channel_id is not set in config")
	}

	// Пересылаем сообщение в канал
	if message.OriginalSender == nil {
		_, err := c.Bot().Copy(
			&tele.Chat{ID: targetChannelID},
			message,
		)
		if err != nil {
			return nil, errors.New("error forwarding message: " + err.Error())
		}
	} else {
		_, err := c.Bot().Forward(
			&tele.Chat{ID: targetChannelID},
			message,
		)
		if err != nil {
			return nil, errors.New("error forwarding message: " + err.Error())
		}
	}

	return args, nil
}

func RegisterTextChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(10*time.Second, connectToDatabase, registerText, forwardToChannel)
}
