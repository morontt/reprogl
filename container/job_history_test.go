package container

import (
	"strconv"
	"testing"
	"time"
)

func TestDurationStrings(t *testing.T) {
	cases := map[int]string{
		1:         "1 мес.",
		1*12 + 3:  "1 год 3 мес.",
		2 * 12:    "2 года",
		2*12 + 4:  "2 года 4 мес.",
		3*12 + 3:  "3 года 3 мес.",
		4*12 + 7:  "4 года 7 мес.",
		5*12 + 4:  "5 лет 4 мес.",
		6*12 + 2:  "6 лет 2 мес.",
		7*12 + 1:  "7 лет 1 мес.",
		8*12 + 9:  "8 лет 9 мес.",
		9*12 + 4:  "9 лет 4 мес.",
		10*12 + 5: "10 лет 5 мес.",
		14 * 12:   "14 лет",
		17*12 + 2: "17 лет 2 мес.",
		20*12 + 3: "20 лет 3 мес.",
		21*12 + 4: "21 год 4 мес.",
		22*12 + 9: "22 года 9 мес.",
		23*12 + 8: "23 года 8 мес.",
		24*12 + 7: "24 года 7 мес.",
		25*12 + 6: "25 лет 6 мес.",
		26*12 + 5: "26 лет 5 мес.",
		27*12 + 4: "27 лет 4 мес.",
		28*12 + 3: "28 лет 3 мес.",
		29*12 + 2: "29 лет 2 мес.",
		30*12 + 1: "30 лет 1 мес.",
	}

	for delta, expected := range cases {
		actual := durationString(delta)
		if actual != expected {
			t.Errorf("durationString error: got %s; want %s", actual, expected)
		}
	}
}

func TestDuration(t *testing.T) {
	cases := []struct {
		job  Job
		want string
	}{
		{
			job: Job{
				Company:   "Company_1",
				StartDate: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "1 мес.",
		},
		{
			job: Job{
				Company:   "Company_2",
				StartDate: time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2020, time.April, 2, 0, 0, 0, 0, time.UTC),
			},
			want: "1 год 1 мес.",
		},
		{
			job: Job{
				Company:   "Company_3",
				StartDate: time.Date(2023, time.June, 4, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, time.June, 18, 0, 0, 0, 0, time.UTC),
			},
			want: "0 мес.",
		},
		{
			job: Job{
				Company:   "Company_4",
				StartDate: time.Date(2023, time.June, 6, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, time.June, 21, 0, 0, 0, 0, time.UTC),
			},
			want: "1 мес.",
		},
		{
			job: Job{
				Company:   "Company_5",
				StartDate: time.Date(2021, time.May, 6, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2022, time.July, 16, 0, 0, 0, 0, time.UTC),
			},
			want: "1 год 2 мес.",
		},
	}

	for idx, item := range cases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			res := item.job.Duration()
			if res != item.want {
				t.Errorf("%s : got %s; want %s", item.job.Company, res, item.want)
			}
		})
	}
}
