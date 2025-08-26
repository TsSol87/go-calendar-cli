package events

import (
	"errors"
	"fmt"
	"github.com/TsSol87/calendarApp/priority"
	"github.com/TsSol87/calendarApp/reminder"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var ErrIsValidTitle = errors.New("Title does not match the required pattern")
var ErrIsValidDate = errors.New("Invalid date format")

const TimeZone = "Asia/Irkutsk"
const DateFormat = "2006-01-02 15:04"

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority priority.Priority  `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

func getNextID() string {
	return uuid.New().String()
}

func IsValidTitle(title string) error {
	pattern := "^[a-zA-Z0-9а-яА-Я ]{3,50}$"
	matched, err := regexp.MatchString(pattern, title)
	if err != nil {
		return fmt.Errorf("regexp error: %w", err)
	}
	if !matched {
		return ErrIsValidTitle
	}
	return nil
}

func TimeParse(dataStr string) (time.Time, error) {
	location, err := time.LoadLocation(TimeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load time zone '%s': %w", TimeZone, err)
	}

	at, err := time.ParseInLocation(DateFormat, dataStr, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date string '%s' with format '%s'", dataStr, DateFormat)
	}

	return at, nil
}

func NewEvent(title string, dateStr string, priorityStr string) (*Event, error) {

	err := IsValidTitle(title)
	if err != nil {
		return nil, fmt.Errorf("can't create event: %w", err)
	}

	t, errTimeParse := TimeParse(dateStr)
	if errTimeParse != nil {
		return nil, fmt.Errorf("can't create date: %w", ErrIsValidDate)
	}

	p := priority.Priority(priorityStr)
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: p,
		Reminder: nil,
	}, nil

}

func (e *Event) Update(title string, dateStr string, priorityStr string) error {
	err := IsValidTitle(title)
	if err != nil {
		return fmt.Errorf("can't create event: %w", err)
	}

	time, errTimeParse := TimeParse(dateStr)
	if errTimeParse != nil {
		return fmt.Errorf("can't create date: %w", errTimeParse)
	}

	p := priority.Priority(priorityStr)
	if err := p.Validate(); err != nil {
		return err
	}

	e.Title = title
	e.StartAt = time
	e.Priority = p
	return nil
}

func (e Event) Print() {
	fmt.Printf("ID: %s  Событие: %s  Дата: %s  Приоритет: %s (Напоминание: %s)\n", e.ID, e.Title, e.StartAt.Format("2006-01-02T15:04:05"), e.Priority, e.Reminder)
}

func (e *Event) AddReminder(message string, at time.Time, notify func(msg string)) error {
	var err error
	e.Reminder, err = reminder.NewReminder(message, at, notify)
	if err != nil {
		return err
	}
	e.Reminder.Start()
	return nil
}

func (e *Event) RemoveReminder() {
	if e.Reminder != nil {
		e.Reminder.Stop()
		e.Reminder = nil
		fmt.Println("Напоминание удалено")
		return
	}
	fmt.Println("Напоминание не найдено")

}
