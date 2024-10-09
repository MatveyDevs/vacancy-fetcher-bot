package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get .env: %v", err)
	}

	bot, err := tg.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("[ERROR] Failed to get TG_BOT_TOKEN: %v", err)
	}
	_ = bot

	db, err := sql.Open("sqlite3", "file:db.db")
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to the database: %v", err)
	}
	defer db.Close()

	v, err := HHRequest()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(v.HHItems)-1; i++ {
		fmt.Printf("%v\n", v.HHItems[i])
	}
}

type HeadHunterResp struct {
	HHItems []Item `json:"items"`
}

type Item struct {
	Name           string     `json:"name"`
	SalaryInfo     Salary     `json:"salary"`
	AddressInfo    Address    `json:"address"`
	PublishedAt    CustomTime `json:"published_at"`
	CreatedAt      CustomTime `json:"created_at"`
	URL            string     `json:"alternate_url"`
	EmployerInfo   Employer   `json:"employer"`
	Description    Snippet    `json:"snippet"`
	ScheduleInfo   Schedule   `json:"schedule"`
	ExperienceInfo Experience `json:"experience"`
	EmploymentInfo Employment `json:"employment"`
}

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	str := string(b[1 : len(b)-1])
	t, err := time.Parse("2006-01-02T15:04:05+0300", str)
	if err != nil {
		return err
	}
	c.Time = t
	return nil
}

type Salary struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Currency string `json:"currency"`
}
type Address struct {
	City     string `json:"city"`
	Street   string `json:"street"`
	Building string `json:"building"`
	Raw      string `json:"raw"`
}
type Employer struct {
	Name string `json:"name"`
	//? trusted ?
}
type Snippet struct {
	Requirement    string `json:"requirement"`
	Responsibility string `json:"responsibility"`
}
type Schedule struct {
	Name string `json:"name"`
}
type Experience struct {
	Name string `json:"name"`
}
type Employment struct {
	Name string `json:"name"`
}

func HHRequest() (*HeadHunterResp, error) {
	url := "https://api.hh.ru/vacancies?area=72&page=1"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	cl := &http.Client{}
	resp, err := cl.Do(r)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	b := &HeadHunterResp{}
	w, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	for i := 0; i < len(b.HHItems)-1; i++ {

	}

	err = json.Unmarshal(w, &b)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return b, nil
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
