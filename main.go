package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
	"vacancy-fetcher-bot/internal/bot"
	"vacancy-fetcher-bot/internal/database"
	"vacancy-fetcher-bot/internal/fetcher"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get .env: %v", err)
	}
	db, err := database.Init()
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to the database: %w", err)
	}
	defer db.Conn.Close()

	TGBotToken := os.Getenv("TG_BOT_TOKEN")
	TGChannelID, _ := strconv.ParseInt(os.Getenv("TG_CHANNEL_ID"), 10, 64)

	bot, err := bot.New(TGBotToken, TGChannelID, db)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get TG_BOT_TOKEN: %w", err)
	}
	log.Printf("Authorized on account %s", bot.API.Self.UserName)

	f := fetcher.New(os.Getenv("BASE_URL"))
	t, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	timeout := time.Duration(t) * time.Second
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		vacs, err := f.Fetch(ctx, "1")
		if err != nil {
			log.Printf("[ERROR] Failed to fetch vacancies: %w", err)
			time.Sleep(time.Minute)
			continue
		}
		err = db.SaveVacancies(vacs)
		if err != nil {
			log.Println(err)
		}

		if err != nil {
			log.Println(err)
		}
		for _, vac := range vacs {
			isPublished, err := bot.DB.IsPublishedVacancy(vac)
			if err != nil {
				log.Println(err)
				continue
			}
			if !isPublished {
				postedVacanci, err := bot.PostVacanci(vac)
				if err != nil {
					log.Println(err)
				}
				bot.DB.SaveVacancyPublication(postedVacanci)
				break
			}
		}

		time.Sleep(time.Hour * 2)
	}
}
