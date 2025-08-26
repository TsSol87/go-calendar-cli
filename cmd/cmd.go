package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TsSol87/calendarApp/calendar"
	"github.com/TsSol87/calendarApp/events"
	"github.com/TsSol87/calendarApp/logger"
	"github.com/TsSol87/calendarApp/priority"
	"github.com/TsSol87/calendarApp/reminder"
	"github.com/TsSol87/calendarApp/storage"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"os"
	"strings"
)

type Cmd struct {
	calendar   *calendar.Calendar
	wg         sync.WaitGroup
	log        Log
	logStorage storage.Store
}

type LogEntry struct {
	Message   string
	Timestamp time.Time
}
type Log struct {
	entries []LogEntry
	mutex   sync.Mutex
}

func (l *Log) Print() {
	for _, e := range l.entries {
		fmt.Printf("CMD(Сообщение): %s\tCMD(Время): %s\n", e.Message, e.Timestamp.Format("2006-01-02T15:04:05"))
	}

}
func NewCmd(c *calendar.Calendar) *Cmd {
	logStorage := storage.NewJsonStorage("log_data.json")
	cmd := &Cmd{
		calendar:   c,
		log:        Log{entries: make([]LogEntry, 0), mutex: sync.Mutex{}},
		logStorage: logStorage,
	}
	cmd.loadLog()
	return cmd

}

func (c *Cmd) Save() error {
	data, err := json.Marshal(c.log.entries)
	if err != nil {

		return err
	}
	err = c.logStorage.Save(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cmd) loadLog() error {
	data, err := c.logStorage.Load()
	if err != nil {
		return err
	}
	c.log.mutex.Lock()
	defer c.log.mutex.Unlock()
	err = json.Unmarshal(data, &c.log.entries)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cmd) executor(input string) {
	parts, err := shlex.Split(input)
	if err != nil {
		fmt.Println("Ошибка разбора команды:", err)
		return
	}
	if len(parts) == 0 {
		return
	}
	logger.Info(input)
	c.log.mutex.Lock()
	defer c.log.mutex.Unlock()
	timestamp := time.Now()
	entry := LogEntry{input, timestamp}
	c.log.entries = append(c.log.entries, entry)
	c.Save()

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		if len(parts) < 4 {
			fmt.Println("Формат: add \"название события\" \"дата и время\" \"приоритет\"")
			return
		}

		title := parts[1]
		date := parts[2]
		priorityStr := (parts[3])

		e, err := c.calendar.AddEvent(title, date, priorityStr)
		if err != nil {
			logMessage := fmt.Sprintf("Error adding event (title: %s, date: %s, priority: %s): %v", title, date, priorityStr, err)
			logger.Error(logMessage)
			c.LogCapture(err)
			if errors.Is(err, events.ErrIsValidTitle) {
				fmt.Printf("Error: Invalid title '%s'. It must contain between 3 and 50 alphanumeric characters and spaces.\n", title)

			} else if errors.Is(err, events.ErrIsValidDate) {
				fmt.Printf("Error: Invalid date format. Please use the format: %s\n", events.DateFormat)

			} else if errors.Is(err, priority.ErrIsValidPriority) {
				fmt.Println("Error: Invalid priority. Please use 'high', 'medium', or 'low'.")

			} else {
				fmt.Printf("can't create event: %v\n", err)

			}
			return
		}

		fmt.Println("Событие", e.Title, "добавлено")

	case "remove":
		if len(parts) < 2 {
			fmt.Println("Формат: remove  \"название ID\"")
			return
		}
		id := parts[1]
		errDel := c.calendar.DeleteEvent(id)
		logMessage := fmt.Sprintf("Error delete event (id: %s): %v", id, errDel)
		logger.Error(logMessage)
		if errDel != nil {
			fmt.Println(errDel)
		} else {
			fmt.Printf("Событие c ключом '%s' удалено\n", id)
		}
	case "update":
		if len(parts) < 5 {
			fmt.Println("Формат: update \"название ID\" \"название события\" \"дата и время\" \"приоритет\"")
			return
		}
		id := parts[1]
		title := parts[2]
		date := parts[3]
		priorityStr := (parts[4])

		err := c.calendar.EditEvent(id, title, date, priorityStr)
		if err != nil {
			logMessage := fmt.Sprintf("Error update event (title: %s, date: %s, priority: %s): %v", title, date, priorityStr, err)
			logger.Error(logMessage)
			c.LogCapture(err)

			if errors.Is(err, events.ErrIsValidTitle) {
				fmt.Printf("Error: Invalid title '%s'. It must contain between 3 and 50 alphanumeric characters and spaces.\n", title)
			} else if errors.Is(err, events.ErrIsValidDate) {
				fmt.Printf("Error: Invalid date format. Please use the format: %s\n", events.DateFormat)
			} else if errors.Is(err, priority.ErrIsValidPriority) {
				fmt.Println("Error: Invalid priority. Please use 'high', 'medium', or 'low'.")
			} else {
				fmt.Printf("can't update event: %v\n", err)
			}
			return
		}

		fmt.Printf("Событие c ключом '%s' изменено\n", id)

	case "list":
		events := c.calendar.GetEvents()
		if len(events) == 0 {
			fmt.Println("Список событий пуст")
			return
		}
		for _, event := range events {
			event.Print()
		}
	case "reminder":
		if len(parts) < 4 {
			fmt.Println("Формат: reminder \"ID события\" \"сообщение\" \"дата и время\"")
			return
		}
		id := parts[1]
		message := parts[2]
		at := parts[3]
		//timer := parts[4]

		err := c.calendar.SetEventReminder(id, message, at)
		if err != nil {
			logMessage := fmt.Sprintf("Error adding reminder (id: %s, message: %s, at: %s): %v", id, message, at, err)
			logger.Error(logMessage)
			if errors.Is(err, reminder.ErrEmptyMessage) {
				fmt.Println("Can't set reminder with empty message")
			} else if errors.Is(err, events.ErrIsValidDate) {
				fmt.Printf("Error: Invalid date format. Please use the format: %s\n", events.DateFormat)
			} else {
				fmt.Println(err)
			}
			return
		}
		fmt.Printf("Напоминание для события c ключом '%s' добавлено\n", id)
	case "cancel-reminder":
		if len(parts) < 2 {
			fmt.Println("Формат: cancel-reminder \"ID события\"")
			return
		}
		id := parts[1]
		errCancelReminder := c.calendar.CancelEventReminder(id)
		if errCancelReminder != nil {
			fmt.Println(errCancelReminder)
		}
	case "history":

		c.log.Print()

	case "help":
		fmt.Println("Доступные команды:")
		fmt.Println("  Добавить событие:\t\tadd \"название события\" \"дата и время\" \"приоритет\"")
		fmt.Println("  Удалить событие:\t\tremove \"ID события\"")
		fmt.Println("  Обновить событие:\t\tupdate \"ID события\" \"название события\" \"дата и время\" \"приоритет\"")
		fmt.Println("  Показать список событий:\tlist")
		fmt.Println("  Установить напоминание:\treminder \"ID события\" \"сообщение\" \"дата и время\" \"таймер\"")
		fmt.Println("  Отменить напоминание:\t\tcancel-reminder \"ID события\"")
		fmt.Println("  Показать историю:\t\thistory")
		fmt.Println("  Выйти из программы:\t\texit")

	case "exit":
		logger.System("app is closed")
		c.calendar.Save()
		close(c.calendar.Notification)
		c.wg.Wait()
		os.Exit(0)

	default:
		fmt.Println("Неизвестная команда:")
		fmt.Println("Введите 'help' для списка команд")
	}
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	if strings.Contains(d.TextBeforeCursor(), " ") {
		return []prompt.Suggest{}
	}
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "update", Description: "Обновить событие"},
		{Text: "reminder", Description: "Добавить напоминание"},
		{Text: "cancel-reminder", Description: "Отменить напоминание"},
		{Text: "help", Description: "Показать справку"},
		{Text: "history", Description: "Показать историю"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) Run() {

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for msg := range c.calendar.Notification {
			fmt.Println(msg)

			timestamp := time.Now()
			entry := LogEntry{msg, timestamp}
			c.log.mutex.Lock()
			c.log.entries = append(c.log.entries, entry)
			c.log.mutex.Unlock()
			c.Save()
		}
		fmt.Println("Канал уведомлений закрыт")
	}()
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	p.Run()
}
func (c *Cmd) LogCapture(err error) {
	msg := err.Error()
	c.log.entries = append(c.log.entries, LogEntry{Message: msg, Timestamp: time.Now()})
	c.Save()
}
