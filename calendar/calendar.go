package calendar

import (
	"encoding/json"
	"github.com/TsSol87/calendarApp/logger"
	"github.com/TsSol87/calendarApp/storage"
	"time"

	//"github.com/TsSol87/calendarApp/storage"

	//"errors"
	"fmt"
	"github.com/TsSol87/calendarApp/events"
	//"time"
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c.calendarEvents)
	if err != nil {

		return err
	}
	err = c.storage.Save(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.calendarEvents)
	if err != nil {
		return err
	}
	return nil
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{calendarEvents: make(map[string]*events.Event), storage: s, Notification: make(chan string)}
}

func (c *Calendar) AddEvent(title string, dateStr string, priorityStr string) (*events.Event, error) {
	e, err := events.NewEvent(title, dateStr, priorityStr)
	if err != nil {
		return nil, err
	}

	c.calendarEvents[e.ID] = e
	errSave := c.Save()
	if errSave != nil {
		return nil, errSave
	}
	return e, nil
}
func (c *Calendar) GetEvents() map[string]*events.Event {
	eventsCopy := make(map[string]*events.Event)
	for key, value := range c.calendarEvents {
		eventsCopy[key] = value
	}
	return eventsCopy
}

func (c *Calendar) DeleteEvent(id string) error {

	_, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}

	delete(c.calendarEvents, id)

	errSave := c.Save()
	if errSave != nil {
		return fmt.Errorf("error saving after deletion: %w", errSave)
	}
	return nil

}

func (c *Calendar) EditEvent(id, title string, date string, priorityStr string) error {

	e, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}

	err := e.Update(title, date, priorityStr)
	if err != nil {
		return err
	}
	errSave := c.Save()
	if errSave != nil {
		return fmt.Errorf("error saving after event change: %w", errSave)
	}
	return nil
}

func (c *Calendar) SetEventReminder(id, message string, dateStr string) error {

	e, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}

	at, errDateStr := events.TimeParse(dateStr)
	if errDateStr != nil {
		return fmt.Errorf("can't create date: %w", events.ErrIsValidDate)
	}

	t := at.UTC()
	now := time.Now().UTC()
	if t.Before(now) {
		return fmt.Errorf("no reminder has been added: time %q has already passed", t)
	}

	err := e.AddReminder(message, t, c.Notify)
	if err != nil {
		return err
	}
	errSave := c.Save()
	if errSave != nil {
		return fmt.Errorf("error saving the calendar: %w", errSave)
	}

	return nil
}

func (c *Calendar) CancelEventReminder(id string) error {

	e, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}
	e.RemoveReminder()
	errSave := c.Save()
	if errSave != nil {
		logMessage := fmt.Sprintf("error saving the calendar: (id: %s): %v", e.ID, errSave)
		logger.Error(logMessage)
		return fmt.Errorf("error saving the calendar: %w", errSave)
	}
	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg

}
