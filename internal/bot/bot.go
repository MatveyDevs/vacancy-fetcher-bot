package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
	"yandex-lms/internal/database"
	"yandex-lms/internal/models"
)

type Bot struct {
	API       *tgbotapi.BotAPI
	ChannelID int64
	DB        *database.Database
}

func New(token string, channelID int64, db *database.Database) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{api, channelID, db}, nil
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
