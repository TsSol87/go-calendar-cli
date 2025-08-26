package reminder

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Reminder struct {
	Message string      `json:"message"`
	At      time.Time   `json:"at"`
	Sent    bool        `json:"sent"`
	Timer   *time.Timer `json:"-"`
	notify  func(msg string)
}

var ErrEmptyMessage = errors.New("message is empty")

func (r *Reminder) String() string {
	if r == nil {
		return "не установлено"
	}
	status := "не отправлено"
	if r.Sent {
		status = "отправлено"
	}
	return fmt.Sprintf("\"%s\", Время: %s, Статус: %s",
		r.Message,
		r.At.Format("2006-01-02 15:04:05"),
		status,
	)
}

func NewReminder(message string, at time.Time, notify func(msg string)) (*Reminder, error) {
	if len(strings.TrimSpace(message)) == 0 {
		return nil, fmt.Errorf("can't create reminder: %w", ErrEmptyMessage)
	}
	return &Reminder{
		Message: message,
		At:      at,
		Sent:    false,
		notify:  notify,
	}, nil

}

func (r *Reminder) Start() {

	duration := time.Until(r.At)
	fmt.Println("Напоминание сработает через:", duration)
	if duration > 0 {
		r.Timer = time.AfterFunc(duration, func() {
			r.Send()
		})
	} else {
		r.Send()
	}
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	r.notify(r.Message)

	r.Sent = true
	fmt.Println("Проверка статуса", r.Sent)
}

func (r *Reminder) Stop() {
	if r.Timer != nil {
		r.Timer.Stop()
		fmt.Println("Таймер остановлен")
	} else {
		fmt.Println("Таймер уже был nil")
	}
}
