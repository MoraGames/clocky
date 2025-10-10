package utils

import (
	"fmt"
	"strings"

	"github.com/go-co-op/gocron"
)

func StringifyJobs(jobs []*gocron.Job) string {
	var stringifiedJobs []string
	for _, job := range jobs {
		stringifiedJobs = append(stringifiedJobs, fmt.Sprintf("%s (%p)", job.GetName(), job))
	}
	return fmt.Sprintf("[%s]", strings.Join(stringifiedJobs, ", "))
}
