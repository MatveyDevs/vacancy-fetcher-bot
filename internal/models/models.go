package models

import (
	"fmt"
	"time"
)

type HHResponse struct {
	HHItems []Item `json:"items"`
}

type Item struct {
	Name        string     `json:"name" db:"name"`
	Salary      Salary     `json:"salary" db:"salary"`
	Address     Address    `json:"address"`
	PublishedAt CustomTime `json:"published_at" db:"published_at"`
	CreatedAt   CustomTime `json:"created_at"`
	URL         string     `json:"alternate_url" db:"url"`
	Employer    Employer   `json:"employer"`
	Description Snippet    `json:"snippet"`
	Schedule    Schedule   `json:"schedule"`
	Experience  Experience `json:"experience"`
	Employment  Employment `json:"employment" `
}

type CustomTime struct {
	time.Time
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
	Raw      string `json:"raw" db:"area"`
}
type Employer struct {
	Name string `json:"name" db:"name"`
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
	Name string `json:"name" db:"name"`
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

func (s *Salary) GetSalary() string {
	return fmt.Sprintf("%d - %d%s", s.From, s.To, s.Currency)
}
