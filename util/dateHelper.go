package util

import (
	"time"
	"strings"
	"strconv"
	"github.com/pkg/errors"
)

var dayOrdinals = map[string]string{ // map[ordinal]cardinal
	"1st": "1", "2nd": "2", "3rd": "3", "4th": "4", "5th": "5",
	"6th": "6", "7th": "7", "8th": "8", "9th": "9", "10th": "10",
	"11th": "11", "12th": "12", "13th": "13", "14th": "14", "15th": "15",
	"16th": "16", "17th": "17", "18th": "18", "19th": "19", "20th": "20",
	"21st": "21", "22nd": "22", "23rd": "23", "24th": "24", "25th": "25",
	"26th": "26", "27th": "27", "28th": "28", "29th": "29", "30th": "30",
	"31st": "31",
}

var monthStrings = map[string]int{
	"january":1, "february":2, "match":3, "april":4,
	"may":5, "june":6, "july":7, "august":8, "september":9, "octomber":10, "november":11, "december":12,
}

// ParseOrdinalDate parses a string time value using an ordinary package time layout.
// Before parsing, an ordinal day, [1st, 31st], is converted to a cardinal day, [1, 31].
// For example, "1st August 2017" is converted to "1 August 2017" before parsing, and
// "August 1st, 2017" is converted to "August 1, 2017" before parsing.
func ParseOrdinalDate(layout, value string) (time.Time, error) {
	value = strings.ToLower(strings.Replace(value, ",", "", -1))
	dateStringSplited := strings.Split(value, " ")

	day, err := strconv.ParseInt(dayOrdinals[dateStringSplited[0]], 10, 64)
	if err != nil {
		return time.Time{}, errors.New("Day out of range")
	}
	monthString := strings.ToLower(dateStringSplited[1])
	month, found := monthStrings[monthString]
	if !found {
		return time.Time{}, errors.New("Month out of range " + monthString)
	}
	year, err := strconv.ParseInt(dateStringSplited[2], 10, 64)

	//hadle and log err in format

	return time.Date(int(year), time.Month(month), int(day), 0,0,0,0, DefaultLocation), nil

	const ( // day number
		cardMinLen = len("1")
		cardMaxLen = len("31")
		ordSfxLen  = len("th")
		ordMinLen  = cardMinLen + ordSfxLen
	)

	for k := 0; k < len(value)-ordMinLen; {
		// i number start
		for ; k < len(value) && (value[k] > '9' || value[k] < '0'); k++ {
		}
		i := k
		// j cardinal end
		for ; k < len(value) && (value[k] <= '9' && value[k] >= '0'); k++ {
		}
		j := k
		if j-i > cardMaxLen || j-i < cardMinLen {
			continue
		}
		// k ordinal end
		// ASCII Latin (uppercase | 0x20) = lowercase
		for ; k < len(value) && (value[k]|0x20 >= 'a' && value[k]|0x20 <= 'z'); k++ {
		}
		if k-j != ordSfxLen {
			continue
		}

		// day ordinal to cardinal
		for ; i < j-1 && (value[i] == '0'); i++ {
		}
		o := strings.ToLower(value[i:k])
		c, ok := dayOrdinals[o]
		if ok {
			value = value[:i] + c + value[k:]
			break
		}
	}

	date, err := time.ParseInLocation(layout, value, DefaultLocation)
	return date, err
}

// Times without a timezone are Hong Kong times.
var DefaultLocation = func(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		loc = time.UTC
	}
	return loc
}(`Africa/Lagos`)
