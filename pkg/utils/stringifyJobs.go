package utils

import (
	"fmt"
	"strings"

	"github.com/go-co-op/gocron/v2"
)

func StringifyJobs(jobs []gocron.Job) string {
	var stringifiedJobs []string
	for _, job := range jobs {
		stringifiedJobs = append(stringifiedJobs, fmt.Sprintf("%s (%v)", job.Name(), job.ID()))
	}
	return fmt.Sprintf("[%s]", strings.Join(stringifiedJobs, ", "))
}
