package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var dateFormatFull0 = regexp.MustCompile(`\d{1,2}\/\d{1,2}\/\d{4}$`)
var dateFormatFull1 = regexp.MustCompile(`\b(January|February|March|April|May|June|July|August|September|October|November|December) \d{1,2},\s?\d{4}`)
var dateFormatMonthDay0 = regexp.MustCompile(`\d{1,2}\/\d{1,2}$`)
var dateFormatMonthDay1 = regexp.MustCompile(`\b(January|February|March|April|May|June|July|August|September|October|November|December) \d{1,2}`)
var dateFormatYear = regexp.MustCompile(`\d{4}`)
var dateFormatMonthYear = regexp.MustCompile(`\b(January|February|March|April|May|June|July|August|September|October|November|December) \d{4}`)
var dateFormatToday = regexp.MustCompile(`today|Today`)
var dateFormatYesterday = regexp.MustCompile(`yesterday|Yesterday`)
var dateFormatThis = regexp.MustCompile(`this \b(week|month|year)`)
var dateFormatLast = regexp.MustCompile(`last \b(week|month|year)`)

func ProcessDate(dateCommand string) (from time.Time, to time.Time, err error) {
	var t time.Time
	if dateFormatFull0.MatchString(dateCommand) {
		// try both formats
		t, err = time.Parse("01/02/2006", dateCommand)
		if err != nil {
			t, err = time.Parse("1/2/2006", dateCommand)
			if err != nil {
				return
			}
		}
		return t, t.AddDate(0, 0, 1), nil

	} else if dateFormatFull1.MatchString(dateCommand) {
		if t, err = time.Parse("January 2, 2006", dateCommand); err != nil {
			return
		}
		return t, t.AddDate(0, 0, 1), nil
	} else if dateFormatMonthDay0.MatchString(dateCommand) {
		if t, err = time.Parse("1/2", dateCommand); err != nil {
			return
		}
		t = t.AddDate(time.Now().Year(), 0, 0)
		return t, t.AddDate(0, 0, 1), nil
	} else if dateFormatMonthDay1.MatchString(dateCommand) {
		if t, err = time.Parse("January 2", dateCommand); err != nil {
			return
		}
		t = t.AddDate(time.Now().Year(), 0, 0)
		return t, t.AddDate(0, 0, 1), nil

	} else if dateFormatYear.MatchString(dateCommand) {
		if t, err = time.Parse("2006", dateCommand); err != nil {
			return
		}
		return t, t.AddDate(1, 0, 0), nil
	} else if dateFormatMonthYear.MatchString(dateCommand) {
		if t, err = time.Parse("January 2006", dateCommand); err != nil {
			return
		}
		return t, t.AddDate(0, 1, 0), nil
	} else if dateFormatToday.MatchString(dateCommand) {

		return beginningOfDay(time.Now()), beginningOfDay(time.Now()).AddDate(0, 0, 1), nil

	} else if dateFormatYesterday.MatchString(dateCommand) {
		return beginningOfDay(time.Now()).AddDate(0, 0, -1), beginningOfDay(time.Now()), nil
	} else if dateFormatThis.MatchString(dateCommand) {
		thisCommandSplit := strings.Split(dateCommand, "this ")
		if len(thisCommandSplit) < 2 {
			err = errors.New("invalid command")
			return
		}
		switch strings.ToLower(thisCommandSplit[1]) {
		case "week":
			weekday := time.Now().Weekday()
			return beginningOfDay(time.Now()).AddDate(0, 0, int(-weekday+1)), beginningOfDay(time.Now()).AddDate(0, 0, int(-weekday+7)), nil
		case "month":
			return time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC), time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.UTC), nil
		case "year":
			return time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC), time.Time{}.AddDate(time.Now().Year()+1, 0, 0), nil
		}
	} else if dateFormatLast.MatchString(dateCommand) {
		agoCommandSplit := strings.Split(dateCommand, "last ")
		if len(agoCommandSplit) < 2 {
			err = errors.New("invalid command")
			return
		}
		switch strings.ToLower(agoCommandSplit[1]) {
		case "week":
			weekday := time.Now().Weekday()
			return beginningOfDay(time.Now()).AddDate(0, 0, int(-weekday+1-7)), beginningOfDay(time.Now()).AddDate(0, 0, int(-weekday+1)), nil
		case "month":
			return time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC), time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.UTC), nil
		case "year":
			return time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC), time.Date(time.Now().Year()+1, 0, 0, 0, 0, 0, 0, time.UTC), nil
		}

	}
	return time.Time{}, time.Time{}, errors.New("Command not supported")
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
