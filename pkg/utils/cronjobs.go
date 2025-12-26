package utils

import (
	"fmt"
	"strings"

	"github.com/go-co-op/gocron/v2"
)

func GetJobByName(scheduler gocron.Scheduler, name string) gocron.Job {
	for _, job := range scheduler.Jobs() {
		if job.Name() == name {
			return job
		}
	}
	return nil
}

func StringifyJobs(jobs []gocron.Job) string {
	var stringifiedJobs []string
	for _, job := range jobs {
		stringifiedJobs = append(stringifiedJobs, fmt.Sprintf("%s (%v)", job.Name(), job.ID()))
	}
	return fmt.Sprintf("[%s]", strings.Join(stringifiedJobs, ", "))
}
