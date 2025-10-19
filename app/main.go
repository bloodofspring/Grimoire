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
			Listen: ":3333",
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: os.Getenv("TELEGRAM_BOT_PUBLIC_URL"),
			},
			MaxConnections: 100,
		},
		Verbose: true,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}
	log.Println("Bot created successfully")

	log.Println("Bloodofspring telegram id:", viper.GetInt64("bloodofspring.telegram_id"))
	b.Use(middleware.Logger())
	bloodofspringOnly := b.Group()
	bloodofspringOnly.Use(middleware.Whitelist(viper.GetInt64("bloodofspring.telegram_id")))
	bloodofspringOnly.Handle(tele.OnText, registertext.RegisterTextChain().Run)

	log.Println("Starting bot...")
	b.Start()
}
