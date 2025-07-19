package bluetasks

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/bushubdegefu/m-playground/configs"
	"github.com/bushubdegefu/m-playground/database"
	"github.com/bushubdegefu/m-playground/logs"

	"github.com/madflojo/tasks"
)

func ScheduledTasks() *tasks.Scheduler {

	//  initalizing scheduler for regullarly running tasks
	scheduler := tasks.New()

	// // Add a task to move to Logs Directory Every Interval, Interval to Be Provided From Configuration File
	databaseLoggerFile, _ := database.LoggerFile("django_auth")
	//  App should not start
	log_file, _ := logs.Logfile("django_auth")
	// Getting clear log interval from env
	clearIntervalLogs, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("CLEAR_LOGS_INTERVAL", "1440"))
	if _, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(clearIntervalLogs) * time.Minute,
		TaskFunc: func() error {
			currentTime := time.Now()
			FileName := fmt.Sprintf("%v-%v-%v-%v-%v", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute())
			//  make sure to replace the names of log files correctly here
			Command := fmt.Sprintf("cp django_auth_blue.log logs/django_auth_blue-%v.log", FileName)
			//Command2 := fmt.Sprintf("cp django_auth_gorm.log logs/django_auth_gorm-%v.log", FileName)
			if _, err := exec.Command("bash", "-c", Command).Output(); err != nil {
				fmt.Printf("error: %v\n", err)
			}

			//if _, err := exec.Command("bash", "-c", Command2).Output(); err != nil {
			//	fmt.Printf("error: %v\n", err)
			//}

			err := databaseLoggerFile.Truncate(0)
			if err != nil {
				fmt.Println("Error truncating gorm logger file:", err)
			}
			lerr := log_file.Truncate(0)
			if lerr != nil {
				fmt.Println("Error truncating log file:", err)
			}
			return nil
		},
	}); err != nil {
		fmt.Println(err)
	}

	return scheduler
}
