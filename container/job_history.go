package container

import "time"

type JobHistory []Job

type Job struct {
	Company   string    `toml:"company"`
	Link      string    `toml:"link"`
	Title     string    `toml:"title"`
	StartDate time.Time `toml:"start"`
	EndDate   time.Time `toml:"end"`
}

func (jh JobHistory) Last() Job {
	var c Job
	for _, job := range jh {
		if job.StartDate.After(c.StartDate) {
			c = job
		}
	}

	return c
}
