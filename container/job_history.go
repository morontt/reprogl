package container

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"time"
)

type JobHistory []Job

type Job struct {
	Company   string    `toml:"company"`
	Link      string    `toml:"link,omitempty"`
	Title     string    `toml:"title"`
	StartDate time.Time `toml:"start"`
	EndDate   time.Time `toml:"end"`

	CustomEmoji string `toml:"emoji,omitempty"`
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

func (jh JobHistory) Sort() JobHistory {
	jhCopy := slices.Clone(jh)
	slices.SortFunc(
		jhCopy,
		func(a, b Job) int {
			return cmp.Compare(b.StartDate.Unix(), a.StartDate.Unix())
		},
	)

	return jhCopy
}

func (j *Job) LinkShort() (result string) {
	matches := regexp.MustCompile(`^https?:\/\/([^\/]+)`).FindStringSubmatch(j.Link)
	if matches != nil && matches[1] != "" {
		result = matches[1]
	}

	return
}

func (j *Job) Start() string {
	return ruMonthName(j.StartDate.Month()) + " " + j.StartDate.Format("2006")
}

func (j *Job) End() string {
	now := time.Now()
	if j.EndDate.After(now) {
		return "текущее"
	}

	return ruMonthName(j.EndDate.Month()) + " " + j.EndDate.Format("2006")
}

func (j *Job) Duration() string {
	var end time.Time
	now := time.Now()
	if j.EndDate.After(now) {
		end = now
	} else {
		end = j.EndDate
	}

	y0 := j.StartDate.Year()
	y1 := end.Year()

	m0 := int(j.StartDate.Month())
	m1 := int(end.Month())

	return durationString(12*(y1-y0) + (m1 - m0) + 1)
}

func (j *Job) Emoji() string {
	var symbols []rune
	if j.CustomEmoji == "electrical" {
		// https://emojipedia.org/man-mechanic
		symbols = []rune{
			rune(0x1F468),
			rune(0x200D),
			rune(0x1F527),
		}
	} else {
		// https://emojipedia.org/man-technologist
		symbols = []rune{
			rune(0x1F468),
			rune(0x200D),
			rune(0x1F4BB),
		}
	}

	return string(symbols)
}

func durationString(delta int) string {
	m := delta % 12
	y := delta / 12

	if y == 0 {
		return fmt.Sprintf("%d мес.", m)
	}

	var monthPart string
	if m > 0 {
		monthPart = fmt.Sprintf(" %d мес.", m)
	}

	if y == 1 {
		return "1 год" + monthPart
	}

	if y == 2 || y == 3 || y == 4 {
		return fmt.Sprintf("%d года", y) + monthPart
	}

	if y >= 5 && y <= 20 {
		return fmt.Sprintf("%d лет", y) + monthPart
	}

	if y%10 == 1 {
		return fmt.Sprintf("%d год", y) + monthPart
	}

	if y%10 == 0 || y%10 >= 5 {
		return fmt.Sprintf("%d лет", y) + monthPart
	}

	return fmt.Sprintf("%d года", y) + monthPart
}

func ruMonthName(m time.Month) (res string) {
	months := []string{
		"Янв",
		"Февр",
		"Март",
		"Апр",
		"Май",
		"Июнь",
		"Июль",
		"Авг",
		"Сент",
		"Окт",
		"Нояб",
		"Дек",
	}

	return months[int(m)-1]
}
