package schedule

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

	log.Infof("Deployment indicators will be computed from following cron : %s", recurrenceCronString)
	log.Infof("Deployment indicators will computed next at %s", scheduler.Next(time.Now()))

	job.AddFunc(recurrenceCronString, func() {
		jobDeploy(scheduler)
	})
	job.Start()

	return
}
