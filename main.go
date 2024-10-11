package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
	"yandex-lms/internal/bot"
	"yandex-lms/internal/database"
	"yandex-lms/internal/fetcher"
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

/*
{
    "items": [
        {
            "id": "104375930",
            "premium": false,
            "name": "Менеджер поддержки",
            "department": null,
            "has_test": false,
            "response_letter_required": false,
            "area": {
                "id": "72",
                "name": "Пермь",
                "url": "https://api.hh.ru/areas/72"
            },
            "salary": null,
            "type": {
                "id": "open",
                "name": "Открытая"
            },
            "address": null,
            "response_url": null,
            "sort_point_distance": null,
            "published_at": "2024-09-18T17:45:18+0300",
            "created_at": "2024-09-18T17:45:18+0300",
            "archived": false,
            "apply_alternate_url": "https://hh.ru/applicant/vacancy_response?vacancyId=104375930",
            "show_logo_in_search": null,
            "insider_interview": null,
            "url": "https://api.hh.ru/vacancies/104375930?host=hh.ru",
            "alternate_url": "https://hh.ru/vacancy/104375930",
            "relations": [],
            "employer": {
                "id": "87021",
                "name": "WILDBERRIES",
                "url": "https://api.hh.ru/employers/87021",
                "alternate_url": "https://hh.ru/employer/87021",
                "logo_urls": {
                    "original": "https://img.hhcdn.ru/employer-logo-original/1042808.jpg",
                    "240": "https://img.hhcdn.ru/employer-logo/5791989.jpeg",
                    "90": "https://img.hhcdn.ru/employer-logo/5791988.jpeg"
                },
                "vacancies_url": "https://api.hh.ru/vacancies?employer_id=87021",
                "accredited_it_employer": false,
                "trusted": true
            },
            "snippet": {
                "requirement": "Желание развиваться в сфере клиентской поддержки и предоставлять высокое качество услуг. Инициативность, ответственность и самостоятельность. Наличие стабильного интернета и компьютера...",
                "responsibility": "Внимательно изучать и отвечать на входящие сообщения от пользователей Wildberries в чате. Понимать проблемы и помогать пользователям решать их. "
            },
            "contacts": null,
            "schedule": {
                "id": "remote",
                "name": "Удаленная работа"
            },
            "working_days": [
                {
                    "id": "only_saturday_and_sunday",
                    "name": "По субботам и воскресеньям"
                }
            ],
            "working_time_intervals": [
                {
                    "id": "from_four_to_six_hours_in_a_day",
                    "name": "Можно сменами по 4-6 часов в день"
                }
            ],
            "working_time_modes": [
                {
                    "id": "start_after_sixteen",
                    "name": "С началом дня после 16:00"
                }
            ],
            "accept_temporary": true,
            "professional_roles": [
                {
                    "id": "121",
                    "name": "Специалист технической поддержки"
                }
            ],
            "accept_incomplete_resumes": true,
            "experience": {
                "id": "noExperience",
                "name": "Нет опыта"
            },
            "employment": {
                "id": "full",
                "name": "Полная занятость"
            },
            "adv_response_url": null,
            "is_adv_vacancy": false,
            "adv_context": null
        }
    ],
    "found": 17520,
    "pages": 2000,
    "page": 0,
    "per_page": 1,
    "clusters": null,
    "arguments": null,
    "fixes": null,
    "suggests": null,
    "alternate_url": "https://hh.ru/search/vacancy?area=72&enable_snippets=true&items_on_page=1"
}
*/

// https://api.hh.ru/vacancies?area=72&per_page=10&page=1
