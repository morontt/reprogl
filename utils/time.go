package utils

import "time"

func RuMonthName(m time.Month, genitive bool) (res string) {
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

	if genitive {
		months[2] = "Марта"
		months[4] = "Мая"
		months[5] = "Июня"
		months[6] = "Июля"
	}

	return months[int(m)-1]
}
