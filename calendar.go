package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	holiday "github.com/holiday-jp/holiday_jp-go"
	"github.com/y-yagi/goext/timeext"
)

// Calendar struct
type Calendar struct {
	DateHeader string
	WeekHeader string
	Body       [6]string
}

var (
	red          = color.New(color.FgRed, color.Bold).SprintFunc()
	blue         = color.New(color.FgBlue, color.Bold).SprintFunc()
	reverseVideo = color.New(color.ReverseVideo).SprintFunc()
	daySpace     = "   "
)

func (calendar *Calendar) setHeader(date time.Time) {
	calendar.DateHeader += fmt.Sprintf("     %d年 %02d月       ", date.Year(), date.Month())
	calendar.WeekHeader += fmt.Sprintf("%s %s %s %s %s %s %s   ", red("日"), "月", "火", "水", "木", "金", blue("土"))
}

func (calendar *Calendar) isNeedNewLine(date time.Time) bool {
	return date.Weekday().String() == "Saturday"
}

func (calendar *Calendar) decoratedDate(date time.Time) string {
	var decoratedDate string
	space := ""
	if date.Day() < 10 {
		space = " "
	}

	if holiday.IsHoliday(date) {
		decoratedDate = red(date.Day())
	} else if date.Weekday().String() == "Sunday" {
		decoratedDate = red(date.Day())
	} else if date.Weekday().String() == "Saturday" {
		decoratedDate = blue(date.Day())
	} else {
		decoratedDate = fmt.Sprintf("%d", date.Day())
	}

	today := time.Now()
	if (date.Year() == today.Year()) && (date.Month() == today.Month()) && date.Day() == today.Day() {
		decoratedDate = reverseVideo(decoratedDate)
	}

	return space + decoratedDate
}

func (calendar *Calendar) adjustSpace(line int) {
	calendar.Body[line] += "  "
}

// Generate calendar
func (calendar *Calendar) Generate(date time.Time) {
	var calDate time.Time

	calendar.setHeader(date)
	firstDate := timeext.BeginningOfMonth(date)
	lastDate := timeext.EndOfMonth(date)

	wday := int(firstDate.Weekday())
	calendar.Body[0] += strings.Repeat(daySpace, wday)

	line := 0

	for i := 1; i < lastDate.Day()+1; i++ {
		calDate = time.Date(date.Year(), date.Month(), i, 0, 0, 0, 0, time.Local)
		calendar.Body[line] += fmt.Sprintf("%2s ", calendar.decoratedDate(calDate))

		if calendar.isNeedNewLine(calDate) {
			calendar.adjustSpace(line)
			line++
		}
	}

	wday = int(lastDate.Weekday())
	if wday == 6 && line == 4 {
		calendar.Body[line] += strings.Repeat(daySpace, 7)
	} else {
		calendar.Body[line] += strings.Repeat(daySpace, 6-wday)
	}
	calendar.adjustSpace(line)

	for line++; line < len(calendar.Body); line++ {
		calendar.Body[line] += strings.Repeat(daySpace, 7)
	}
	calendar.adjustSpace(line - 1)
}

// Show calendar
func (calendar *Calendar) Show(w io.Writer, center bool) {
	termWidth := getTerminalWidth() // ターミナル幅をここで取得

    if center {
        // 各行を中央寄せして出力
        fmt.Fprint(w, centerLine(calendar.DateHeader, termWidth)+"\n")
        fmt.Fprint(w, centerLine(calendar.WeekHeader, termWidth)+"\n")
        for _, element := range calendar.Body {
            fmt.Fprint(w, centerLine(element, termWidth)+"\n")
        }
    } else {
        // center が false の場合は元の表示ロジック
        fmt.Fprint(w, calendar.DateHeader+"\n")
        fmt.Fprint(w, calendar.WeekHeader+"\n")
        for _, element := range calendar.Body {
            fmt.Fprint(w, element+"\n")
        }
    }
}

// Clear calendar
func (calendar *Calendar) Clear() {
	calendar.DateHeader = ""
	calendar.WeekHeader = ""
	for i := range calendar.Body {
		calendar.Body[i] = ""
	}
}
