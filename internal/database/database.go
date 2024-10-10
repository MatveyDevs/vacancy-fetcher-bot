package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"yandex-lms/internal/models"
)

type Database struct {
	DB *sql.DB
}

func Init() (*Database, error) {
	conn, err := sql.Open("sqlite3", "file:db.db")
	if err != nil {
		return nil, err
	}
	db := &Database{conn}
	return db, nil
}

func SaveVacancies(db *sql.DB, vacancies []models.Item) error {
	query := fmt.Sprintln("INSERT INTO vacancies (name, employer, area, salary, url, published_at, employment) VALUES (?, ?, ?, ?, ?, ?)")
	var wg sync.WaitGroup
	var errCh = make(chan error, len(vacancies))

	for _, value := range vacancies {
		wg.Add(1)
		go func(v models.Item) {
			defer wg.Done()
			_, err := db.Exec(query,
				v.Name,
				v.Employer.Name,
				v.Address.Raw,
				v.Salary.GetSalary(),
				v.URL,
				v.PublishedAt,
				v.Employment.Name,
			)

			if err != nil {
				errCh <- err
			}
		}(value)
	}
	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		} else {
			log.Printf("Error saving vacancy: %v", err)
		}
	}
	return finalErr
}
