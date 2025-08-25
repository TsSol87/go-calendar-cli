package main

import (
	"fmt"
	"github.com/TsSol87/calendarApp/calendar"
	"github.com/TsSol87/calendarApp/cmd"
	"github.com/TsSol87/calendarApp/logger"
	"github.com/TsSol87/calendarApp/storage"
	//"github.com/TsSol87/calendarApp/events"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	defer logger.Close()
	logger.System("app is started")
	fmt.Println("Введите команду... или введите help для справки")
	s := storage.NewJsonStorage("calendar_data.json")

	//zs := storage.NewZipStorage("calendar_data.zip")
	c := calendar.NewCalendar(s)
	err := c.Load()
	if err != nil {
		logMessage := fmt.Sprintf("Data upload error: (file: %s): %v", s.GetFilename(), err)
		logger.Error(logMessage)
		fmt.Println("Data upload error:", err)
		return
	}

	cli := cmd.NewCmd(c)
	cli.Run()
	defer func() {
		err := c.Save()
		if err != nil {
			logMessage := fmt.Sprintf("Data saving error: (file: %s): %v", s.GetFilename(), err)
			logger.Error(logMessage)
			fmt.Println("Error:", err)
		}
	}()

}
