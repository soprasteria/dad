package jobs

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

// RunBackgroundJobs schedules background tasks as cron jobs
func RunBackgroundJobs() {

	recurrenceCronString := viper.GetString("tasks.recurrence")

	job := cron.New()

	scheduler, err := cron.Parse(recurrenceCronString)
	if err != nil {
		log.WithError(err).WithField("Deployment indicators recurrence", recurrenceCronString).Error("Unable to parse indicators recurrence")
		return
	}

	log.WithFields(log.Fields{
		"cron":          recurrenceCronString,
		"nextExecution": scheduler.Next(time.Now()),
	}).Info("Cron configuration")

	job.AddFunc(recurrenceCronString, func() {
		jobDeploy(scheduler)
	})
	job.Start()

	return
}
