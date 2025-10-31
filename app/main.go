package main

import (
	"log"
	"os"

	"grimoire/database"
	registertext "grimoire/handlers/registerText"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {
	viper.SetConfigName("botConfig.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	godotenv.Load()

	log.Println("Config loaded")
	log.Printf("Bot token: %s", os.Getenv("TELEGRAM_BOT_TOKEN")[:10]+"...")
	log.Printf("Public URL: %s", os.Getenv("TELEGRAM_BOT_PUBLIC_URL"))

	// Initialize database
	err = database.InitDb()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Println("Database initialized")

	pref := tele.Settings{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tele.Webhook{
			Listen: ":7760", // 7760
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: os.Getenv("TELEGRAM_BOT_PUBLIC_URL"),
			},
			MaxConnections: 100,
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}
	log.Println("Bot created successfully")

	log.Println("Bloodofspring telegram id:", viper.GetInt64("bloodofspring.telegram_id"))
	sourceGroupID := viper.GetInt64("bloodofspring.source_group_id")
	log.Println("Source group ID:", sourceGroupID)
	targetChannelID := viper.GetInt64("bloodofspring.target_channel_id")
	log.Println("Target channel ID:", targetChannelID)

	b.Use(middleware.Logger())

	// Обрабатываем сообщения из конкретной группы
	b.Use(tele.MiddlewareFunc(func(hf tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			// Проверяем, что сообщение пришло из нужной группы
			if c.Chat().ID != sourceGroupID {
				return nil
			}
			// Обрабатываем только групповые чаты или супергруппы
			if c.Chat().Type != tele.ChatGroup && c.Chat().Type != tele.ChatSuperGroup {
				return nil
			}
			return hf(c)
		}
	}))

	// Обрабатываем все типы сообщений (текст, медиа и т.д.)
	b.Handle(tele.OnText, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnPhoto, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnVideo, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnDocument, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnAudio, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnVoice, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnVideoNote, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnSticker, registertext.RegisterTextChain().Run)
	b.Handle(tele.OnAnimation, registertext.RegisterTextChain().Run)

	log.Println("Starting bot...")
	b.Start()
}
