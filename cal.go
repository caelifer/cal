package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	WeekDayNames = []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	DaysOfMonth  = []string{
		" 1", " 2", " 3", " 4", " 5", " 6", " 7", " 8", " 9",
		"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
		"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
		"30", "31",
	}
)

func main() {
	flag.Parse()

	switch flag.NArg() {
	case 1:
		year, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatalf("bad year: %v", err)
		}
		fmt.Printf("%v", printYear(year))

	case 2:
		month, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatalf("bad month: %v", err)
		}

		year, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Fatalf("bad year: %v", err)
		}

		if month < 1 || month > 12 {
			log.Fatalf("bad month: %v", month)
		}

		fmt.Printf("%v", printMonth(year, time.Month(month)))

	default:
		t0 := time.Now()
		fmt.Printf("%v", printMonth(t0.Year(), t0.Month()))
	}
}

func printMonth(year int, month time.Month) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", Calendar(year, time.Month(month)))
	return buf.String()
}

func printYear(year int) string {
	var buf bytes.Buffer

	ss := [][]*M{
		{Calendar(year, time.January), Calendar(year, time.February), Calendar(year, time.March)},
		{Calendar(year, time.April), Calendar(year, time.May), Calendar(year, time.June)},
		{Calendar(year, time.July), Calendar(year, time.August), Calendar(year, time.September)},
		{Calendar(year, time.October), Calendar(year, time.November), Calendar(year, time.December)},
	}

	// Year heading
	fmt.Fprintln(&buf, mid(flag.Arg(0), 64))

	// N-month sections
	for i, s := range ss {
		if i > 0 {
			fmt.Fprintln(&buf)
		}
		printSection(&buf, s...)
	}
	return buf.String()
}

func printSection(buf *bytes.Buffer, ms ...*M) {
	mc := len(ms)

	mh := make([]string, mc)
	wh := make([]string, mc)
	dr := make([][][]string, mc)

	for i, m := range ms {
		mh[i] = mid(m.Month(), 20)
		wh[i] = strings.Join(WeekDayNames, " ")
		dr[i] = m.Grid()
	}

	// 3-month name headings
	fmt.Fprintln(buf, strings.Join(mh, "  "))
	// 3-month weekdays headings
	fmt.Fprintln(buf, strings.Join(wh, "  "))

	// Print the rest of the lines
	for i := 0; i < 6; i++ {
		row := make([]string, mc)
		for j, c := range dr {
			// subrow for each month
			row[j] = strings.Join(c[i], " ")
		}
		// print formatted row
		fmt.Fprintln(buf, strings.Join(row, "  "))
	}
}

func Calendar(year int, month time.Month) *M {
	return &M{
		year:  year,
		month: month,
	}
}

type M struct {
	year  int
	month time.Month
}

func (m *M) String() string {
	var buf bytes.Buffer

	// Write title (Month YEAR)
	buf.WriteString(mid(fmt.Sprintf("%s %d", m.month, m.year), 20))
	buf.WriteString("\n")

	// Write header
	buf.WriteString(strings.Join(WeekDayNames, " "))
	buf.WriteString("\n")

	for _, v := range m.Grid() {
		buf.WriteString(strings.Join(v, " "))
		buf.WriteString("\n")
	}

	return buf.String()
}

func (m M) Month() string {
	return m.month.String()
}

// Two-dimentional grid of strings representing month's calendar
func (m *M) Grid() [][]string {
	// Clone days limited by number of days in this month.
	days := clone(DaysOfMonth[:daysIn(m.year, m.month)])

	// Highlight current day if present
	t0 := time.Now()
	if t0.Year() == m.year && t0.Month() == m.month {
		// Day is specified as a zero-based offset
		highlight(days, t0.Day()-1)
	}

	// Calculate first day of a given month
	d1 := time.Date(m.year, m.month, 1, 0, 0, 0, 0, time.Local)

	// Create Sun-first padded callendar representation
	cal := append(repeat("  ", int(d1.Weekday())), days...)

	// Right pad if required
	rpad := 42 - len(cal)
	cal = append(cal, repeat("  ", rpad)...)

	// Construct a 6x7 grid of days
	grid := make([][]string, 6)
	for i := 0; i < 6; i++ {
		grid[i] = cal[i*7 : (i+1)*7]
	}

	return grid
}

func daysIn(year int, month time.Month) int {
	// Calculates last day in a month by exploiting zero-day calculation returning last day of a previous month.
	//   Taken from: https://brandur.org/fragments/go-days-in-month
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

// mid returns a new string with centered `s` based on `width` by left padding it with blank characters.
func mid(s string, width int) string {
	// Split str for convenience accessing individual rune
	runes := []rune(s)

	// Calculate distance
	ln := len(runes)

	// If msg don't fit, return truncated msg; special case - return original msg if fit exactly.
	dis := width - ln
	if dis <= 0 {
		return string(runes[:width])
	}

	// Calculate left offset using integer division (left padding)
	off := dis / 2

	// Prepare new blank filled string
	res := bytes.Runes(bytes.Repeat([]byte{' '}, width))

	// Splice str's runes starting at the calculated offset
	for i, r := range runes {
		res[i+off] = r
	}

	return string(res)
}

func clone[T any](args []T) []T {
	return append(make([]T, 0, len(args)), args...)
}

func repeat(s string, n int) []string {
	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, s)
	}
	return res
}

func highlight(dates []string, idx int) {
	// Select highlight color scheme
	hl := color.New(color.Bold, color.FgBlack, color.BgHiWhite).SprintFunc()
	// Highlight specified element
	dates[idx] = hl(dates[idx])
}

// vim: :ts=4:sw=4:noexpandtab:ai
