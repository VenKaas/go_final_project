package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/VenKaas/go_final_project/dformat"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	dateInTimeFormat, err := time.Parse(dformat.DFormat, date)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга формата заданной даты: %w", err)
	}
	if repeat == "" {
		return "", fmt.Errorf("ошибка формата repeat, пустая строка")
	}

	switch repeat[0] {
	case 'd':
		return addDays(now, dateInTimeFormat, repeat)
	case 'y':
		return addYear(now, dateInTimeFormat)
	case 'w':
		return addWeek(now, dateInTimeFormat, repeat)
	case 'm':
		return addMonth(now, dateInTimeFormat, repeat)
	default:
		return "", fmt.Errorf("ошибка формата repeat, первый символ не допустим")
	}

}

// функция поиска числа для переноса задачи, если надо переносить на заданные дни
func addDays(now time.Time, dateInTimeFormat time.Time, repeat string) (string, error) {
	daySlice := strings.Split(repeat, " ")
	if len(daySlice) != 2 {
		return "", fmt.Errorf("ошибка формата repeat, количество дней некорректно задано")
	}
	//переводим в число параметр, на сколько дней надо перенести задачу
	dayCount, err := strconv.Atoi(daySlice[1])
	if err != nil {
		return "", fmt.Errorf("ошибка формата repeat, количество дней не число: %w", err)
	}
	if dayCount > 400 {
		return "", fmt.Errorf("ошибка формата repeat, количество дней больше 400")
	}

	//устанавливаем новую дату для задачи
	dateInTimeFormat = dateInTimeFormat.AddDate(0, 0, dayCount)
	for dateInTimeFormat.Format(dformat.DFormat) <= now.Format(dformat.DFormat) {
		dateInTimeFormat = dateInTimeFormat.AddDate(0, 0, dayCount)
	}

	return dateInTimeFormat.Format(dformat.DFormat), nil
}

// функция поиска числа для переноса задачи, если надо переносить на недели
func addWeek(now time.Time, dateInTimeFormat time.Time, repeat string) (string, error) {
	daySlice := strings.Split(repeat, " ")
	if len(daySlice) != 2 {
		return "", fmt.Errorf("ошибка формата repeat, дни недели некорректно заданы")
	}
	weekDaySlice := strings.Split(daySlice[1], ",")
	//делаем мапу доспустимых значений дней, на которые можно перенести задачу
	validDays := make(map[int]string)
	for _, day := range weekDaySlice {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return "", fmt.Errorf("ошибка формата repeat, день недели не число")
		}
		if dayInt > 7 || dayInt < 1 {
			return "", fmt.Errorf("ошибка формата repeat, день недели не от 1 до 7")
		}
		validDays[dayInt] = day
	}

	//если дата меньше текущей, то устанавливаем текущую
	if dateInTimeFormat.Format(dformat.DFormat) < now.Format(dformat.DFormat) {
		dateInTimeFormat = now
	}

	//если в мапе есть допустимый ключ для переноса дня заполняем данные для этого ключа правильной строкой с датой
	for i := 1; i <= 7; i++ {
		dateInTimeFormat = dateInTimeFormat.AddDate(0, 0, 1)
		weekDay := weekDayNumber(dateInTimeFormat)
		_, ok := validDays[weekDay]
		if ok {
			validDays[weekDay] = dateInTimeFormat.Format(dformat.DFormat)
		}
	}
	//ищем ближайшую дату
	targetDay := dateInTimeFormat.Format(dformat.DFormat)
	for _, validDay := range validDays {
		if validDay < targetDay {
			targetDay = validDay
		}
	}
	return targetDay, nil
}

// поиск числа дня недели (пн - 1, вс - 7, и т.д.)
func weekDayNumber(day time.Time) int {
	weekDay := day.Weekday()
	if int(weekDay) == 0 {
		return 7
	}
	return int(weekDay)
}

// функция поиска числа для переноса задачи, когда надо переносить на год
func addYear(now time.Time, dateInTimeFormat time.Time) (string, error) {
	//устанавливаем новую дату для задачи
	dateInTimeFormat = dateInTimeFormat.AddDate(1, 0, 0)
	for dateInTimeFormat.Format(dformat.DFormat) <= now.Format(dformat.DFormat) {
		dateInTimeFormat = dateInTimeFormat.AddDate(1, 0, 0)
	}
	return dateInTimeFormat.Format(dformat.DFormat), nil
}

// функция поиска числа для переноса задачи, когда надо переносить на числа месяца
func addMonth(now time.Time, dateInTimeFormat time.Time, repeat string) (string, error) {
	monthSlice := strings.Split(repeat, " ")
	var monthNumDaySlice []string
	var monthNumSlice []string
	var validDays map[int]string
	var validMonths map[int]string
	var err error
	// проверяем формат месяцев
	if len(monthSlice) == 2 {
		monthNumDaySlice = strings.Split(monthSlice[1], ",")
		validDays, err = mapValidDays(monthNumDaySlice)
		if err != nil {
			return "", fmt.Errorf("в параметрах несть невалидные номера дней %v", err)
		}
	} else if len(monthSlice) == 3 {
		monthNumDaySlice = strings.Split(monthSlice[1], ",")
		validDays, err = mapValidDays(monthNumDaySlice)
		if err != nil {
			return "", fmt.Errorf("в параметрах несть невалидные номера дней %v", err)
		}
		monthNumSlice = strings.Split(monthSlice[2], ",")
		validMonths, err = mapValidMonth(monthNumSlice)
		if err != nil {
			return "", fmt.Errorf("в параметрах несть невалидные номера месяцев %v", err)
		}
	} else {
		return "", fmt.Errorf("неизвестная ошибка формата repeat для месяцев")
	}

	//проверяем, что считаем от даты которая дальше чем сегодня, иначе считаем с сегодня
	if dateInTimeFormat.Format(dformat.DFormat) < now.Format(dformat.DFormat) {
		dateInTimeFormat = now
	}

	if len(validMonths) == 0 {
		return setDateFromCurrentMonth(validDays, dateInTimeFormat)
	} else {
		return setDateForSpecificMonths(validDays, validMonths, dateInTimeFormat)
	}
}

// парсинг параметра с номерами месяцев в мапу
func mapValidMonth(monthNumSlice []string) (map[int]string, error) {
	//делаем мапу доспустимых значений месяцев, на которые можно перенести задачу
	validMonth := make(map[int]string)
	for _, month := range monthNumSlice {
		monthInt, err := strconv.Atoi(month)
		if err != nil {
			return nil, fmt.Errorf("ошибка формата repeat, номер месяца не число")
		}
		if monthInt > 12 || monthInt < 1 {
			return nil, fmt.Errorf("ошибка формата repeat, месяц задан некорректным числом")
		}
		validMonth[monthInt] = month
	}
	return validMonth, nil
}

// парсинг параметра с номерами дней в мапу
func mapValidDays(monthNumDaySlice []string) (map[int]string, error) {
	//делаем мапу доспустимых значений дней, на которые можно перенести задачу
	validDays := make(map[int]string)
	for _, day := range monthNumDaySlice {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return nil, fmt.Errorf("ошибка формата repeat, день недели не число")
		}
		if dayInt > 31 || dayInt < -2 || dayInt == 0 {
			return nil, fmt.Errorf("ошибка формата repeat, день задан некорректным числом")
		}
		validDays[dayInt] = day
	}
	return validDays, nil
}

// поиск даты если не заданы месяца в параметрах
func setDateFromCurrentMonth(validDays map[int]string, dateInTimeFormat time.Time) (string, error) {
	amountOfDaysInMonth := daysInMonth(dateInTimeFormat)
	//если в мапе есть допустимый ключ для переноса дня, заполняем данные для этого ключа правильной строкой с датой
	for i := 1; i <= 31; i++ {
		_, ok := validDays[i]
		if ok {
			if (i <= int(dateInTimeFormat.Day())) || (i > amountOfDaysInMonth) {
				varDate := dateInTimeFormat.AddDate(0, 1, 0)
				validDays[i] = varDate.Format("200601") + fmt.Sprintf("%02d", i)
			} else {
				validDays[i] = dateInTimeFormat.Format("200601") + fmt.Sprintf("%02d", i)
			}
		}
	}

	if _, ok := validDays[-1]; ok {
		validDays[-1] = dateInTimeFormat.AddDate(0, 1, -dateInTimeFormat.Day()).Format(dformat.DFormat)
	}

	if _, ok := validDays[-2]; ok {
		validDays[-2] = dateInTimeFormat.AddDate(0, 1, -dateInTimeFormat.Day()-1).Format(dformat.DFormat)
	}

	//возвращаем ближайшую дату
	return minDate(validDays), nil
}

// поиск даты если заданы месяца в параметрах
func setDateForSpecificMonths(validDays map[int]string, validMonths map[int]string, dateInTimeFormat time.Time) (string, error) {
	resultDate := make(map[int]string)
	var addYear bool
	for j, validMonth := range validMonths {

		if fmt.Sprintf("%02s", validMonth) < dateInTimeFormat.Format("01") {
			for i, validDay := range validDays {
				addYear = true
				resultDate[i+j*100] = dateFormat(dateInTimeFormat, validMonth, validDay, addYear)
			}
		} else if fmt.Sprintf("%02s", validMonth) == dateInTimeFormat.Format("01") {
			for i, validDay := range validDays {
				if fmt.Sprintf("%02s", validDay) > dateInTimeFormat.Format("02") {
					addYear = false
					resultDate[i+j*100] = dateFormat(dateInTimeFormat, validMonth, validDay, addYear)
				} else {
					addYear = true
					resultDate[i+j*100] = dateFormat(dateInTimeFormat, validMonth, validDay, addYear)
				}
			}
		} else {
			for i, validDay := range validDays {
				addYear = false
				resultDate[i+j*100] = dateFormat(dateInTimeFormat, validMonth, validDay, addYear)
			}
		}

	}

	return minDate(resultDate), nil
}

// поиск минимально значения в мапе с подходящими датами
func minDate(validDays map[int]string) string {
	targetDay := time.Now().AddDate(2, 0, 0).Format(dformat.DFormat)
	for _, validDay := range validDays {
		if (validDay < targetDay) && (validDay != "") {
			targetDay = validDay
		}
	}
	return targetDay
}

// фукнция вычисления номера последнего дня в месяце
func daysInMonth(date time.Time) int {
	return int(date.AddDate(0, 1, -date.Day()).Day())
}

// форматирование даты
func dateFormat(date time.Time, month string, day string, addYear bool) string {
	if addYear {
		return date.AddDate(1, 0, 0).Format("2006") + fmt.Sprintf("%02s", month) + fmt.Sprintf("%02s", day)
	}
	return date.Format("2006") + fmt.Sprintf("%02s", month) + fmt.Sprintf("%02s", day)
}
