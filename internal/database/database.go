package database

import (
	"database/sql"
	"log"
	"vacancy-fetcher-bot/internal/models"
)
import _ "github.com/ncruces/go-sqlite3/driver"
import _ "github.com/ncruces/go-sqlite3/embed"

type Database struct {
	Conn *sql.DB
}

func Init() *Database {
	conn, err := sql.Open("sqlite3", "file:../db.db")
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to the database: %w", err)
		return nil
	}
	db := &Database{conn}
	return db
}

func (db *Database) SaveVacancies(vacancies []models.Item) error {
	var exists bool
	for _, v := range vacancies {
		err := db.Conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM vacancies WHERE url = ?)`, v.URL).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			_, err = db.Conn.Exec(`INSERT INTO vacancies (name, employer, area, salary, url, published_at, employment)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				v.Name,
				v.Employer.Name,
				v.Address.Raw,
				v.Salary.GetSalary(),
				v.URL,
				v.PublishedAt.Time,
				v.Employment.Name,
			)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Conn.Exec(`UPDATE vacancies
			SET name = ?, employer = ?, area = ?, salary = ?, published_at = ?, employment = ?
			WHERE url = ?`,
				v.Name,
				v.Employer.Name,
				v.Address.Raw,
				v.Salary.GetSalary(),
				v.PublishedAt.Time,
				v.Employment.Name,
				v.URL,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *Database) SaveVacancyPublication(vacancy *models.VacanciPublication) error {
	_, err := db.Conn.Exec(`INSERT INTO vacancy_publication 
    (vacancy_url, tg_bot_published_at) VALUES (?, ?)`, vacancy.URL, vacancy.TgBotPublishedAt)
	return err
}

func (db *Database) IsPublishedVacancy(vacancy models.Item) (bool, error) {
	var exists bool
	err := db.Conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM vacancy_publication WHERE vacancy_url = ?)`, vacancy.URL).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ClearVacancies(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM vacancies")
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM sqlite_sequence WHERE name='vacancies'")
	return err
}
