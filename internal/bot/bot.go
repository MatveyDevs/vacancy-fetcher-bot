package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
	"vacancy-fetcher-bot/internal/database"
	"vacancy-fetcher-bot/internal/models"
)

type Bot struct {
	API       *tgbotapi.BotAPI
	ChannelID int64
	DB        *database.Database
}

func New(token string, channelID int64, db *database.Database) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get TG_BOT_TOKEN: %w", err)
		return nil
	}
	bot := &Bot{api, channelID, db}
	log.Printf("Authorized on account %s", bot.API.Self.UserName)
	return bot
}

func (b *Bot) PostVacanci(vacanci models.Item) (*models.VacanciPublication, error) {
	text := fmt.Sprintf(`
%s
От %d ₽ на руки

Обязанности:
%s

Требования:
%s

Мы предлагаем:
- Официальное трудоустройство.
- График работы: %s.
- Оплата: %d-%d %s.

Контакты:
%s
    `,
		vacanci.Name,
		vacanci.Salary.From,
		vacanci.Description.Responsibility,
		vacanci.Description.Requirement,
		vacanci.Schedule.Name,
		vacanci.Salary.From,
		vacanci.Salary.To,
		vacanci.Salary.Currency,
		vacanci.URL,
	)
	msg := tgbotapi.NewMessage(b.ChannelID, text)
	if _, err := b.API.Send(msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	publishedVacanci := models.VacanciPublication{
		vacanci.URL,
		time.Now(),
	}
	if err := b.DB.SaveVacancyPublication(&publishedVacanci); err != nil {
		return nil, fmt.Errorf("failed to save vacanci publication: %w", err)
	}

	return &publishedVacanci, nil
}
