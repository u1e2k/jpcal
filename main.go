package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/y-yagi/goext/timeext"
)

var (
	version = "devel"
)

func showYearCalendar(specifyYear string, w io.Writer, center bool) {
	var calendar Calendar

	year, err := time.Parse("2006", specifyYear)
	if err != nil {
		fmt.Printf("Year parse error: %s\n", err)
		os.Exit(1)
	}

	for i := 1; i < 13; i++ {
		date := time.Date(year.Year(), time.Month(i), 1, 0, 0, 0, 0, time.Local)
		calendar.Generate(date)

		if i%3 == 0 {
			calendar.Show(w, center)
			calendar.Clear()
		}
	}
}

func showThreeMonthsCalendar(w io.Writer, center bool) {
	var calendar Calendar
	date := time.Now()

	calendar.Generate(timeext.BeginningOfMonth(date).AddDate(0, 0, -1))
	calendar.Generate(date)
	calendar.Generate(timeext.EndOfMonth(date).AddDate(0, 0, 1))
	calendar.Show(w, center)
}

func showOneMonthCalendar(specifyDate string, w io.Writer, center bool) {
	var calendar Calendar
	var err error

	date := time.Now()

	if len(specifyDate) > 0 {
		date, err = time.Parse("2006-01", specifyDate)
		if err != nil {
			fmt.Printf("Date parse error: %s\n", err)
			os.Exit(1)
		}
	}

	calendar.Generate(date)
	calendar.Show(w, center)
}

func showBeforeCalendar(number int, w io.Writer, center bool) {
	var calendar Calendar

	date := time.Now()
	count := number + 1
	date = date.AddDate(0, -count, 0)

	for i := 1; i <= count; i++ {
		date = timeext.EndOfMonth(date).AddDate(0, 0, 1)
		calendar.Generate(date)
		if i%3 == 0 {
			calendar.Show(w, center)
			calendar.Clear()
		}
	}

	if count%3 != 0 {
		calendar.Show(w, center)
	}
}

func showAfterCalendar(number int, w io.Writer, center bool) {
	var calendar Calendar
	date := time.Now()
	i := 1

	calendar.Generate(date)
	for ; i <= number; i++ {
		if i%3 == 0 {
			calendar.Show(w, center)
			calendar.Clear()
		}
		date = timeext.EndOfMonth(date).AddDate(0, 0, 1)
		calendar.Generate(date)
	}
	calendar.Show(w, center)
}

func run(args []string, out, err io.Writer) int {
	var showVersion bool
	var specifyDate string
	var showYear bool
	var three bool
	var before int
	var after int
	var centerOutput bool
	var test bool

	specifyYear := strconv.Itoa(time.Now().Year())

	flags := flag.NewFlagSet("jpcal", flag.ExitOnError)
	flags.SetOutput(err)
	flags.StringVar(&specifyDate, "d", "", "Use yyyy-mm as the date.")
	flags.BoolVar(&showYear, "y", false, "Use yyyy as the year.")
	flags.BoolVar(&three, "3", false, "Display the previous, current and next month surrounding today.")
	flags.BoolVar(&showVersion, "v", false, "show version")
	flags.IntVar(&before, "B", 0, "Display the number of months before the current month.")
	flags.IntVar(&after, "A", 0, "Display the number of months after the current month.")
	flags.BoolVar(&centerOutput, "c", false, "Center the calendar output.") 
	flags.BoolVar(&test, "t", false, "test mode.")
	flags.Parse(args[1:])

	if showVersion {
		fmt.Fprintln(out, "version:", version)
		return 0
	}

	if test {
		width := getTerminalWidth()
		fmt.Printf("現在のターミナル幅: %d 文字\n", width)
		return 0
	}

	if showYear {
		if len(args) > 2 {
			specifyYear = args[2]
		}

		showYearCalendar(specifyYear, out, centerOutput)
	} else if three {
		showThreeMonthsCalendar(out, centerOutput)
	} else if after > 0 {
		showAfterCalendar(after, out, centerOutput)
	} else if before > 0 {
		showBeforeCalendar(before, out, centerOutput)
	} else {
		showOneMonthCalendar(specifyDate, out, centerOutput)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
