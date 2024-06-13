package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type WorkDay struct {
	StartDay time.Time
	EndDay   time.Time
	Breaks   []Break
	Comment  string
}

type Break struct {
	Start time.Time
	End   time.Time
}

type CSV [][]string

const (
	DateFormat     = "02.01.2006"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "02.01.2006 15:04:05"
	HoursInWorkDay = 8.0
	CSVFileName    = "worktime.csv"
)

var workDay WorkDay
var inBreak bool

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Available commands:")
	fmt.Println("  start day    - Start the work day")
	fmt.Println("  break start  - Start a break")
	fmt.Println("  break end    - End a break")
	fmt.Println("  end day      - End the work day")
	fmt.Println("  add free day - Add a free day")
	fmt.Println("  sort csv     - Sort the CSV file by date")
	fmt.Println("  take overtime- Take overtime")
	fmt.Println("  current overtime - Display current overtime")

	for {
		fmt.Print("Enter command: ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		switch cmd {
		case "start day":
			startDay()
		case "break start":
			startBreak()
		case "break end":
			endBreak()
		case "end day":
			endDay()
			return // End the program after ending the day
		case "add free day":
			addFreeDay()
			return // End the program after adding a free day
		case "sort csv":
			sortCSVByDate()
		case "take overtime":
			takeOvertime()
			return // End the program after taking overtime
		case "current overtime":
			currentOvertime()
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  start day    - Start the work day")
			fmt.Println("  break start  - Start a break")
			fmt.Println("  break end    - End a break")
			fmt.Println("  end day      - End the work day")
			fmt.Println("  add free day - Add a free day")
			fmt.Println("  sort csv     - Sort the CSV file by date")
			fmt.Println("  take overtime- Take overtime")
			fmt.Println("  current overtime - Display current overtime")
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func startDay() {
	workDay = WorkDay{StartDay: time.Now().Round(0)}
	inBreak = false
	fmt.Println("Work day started at", workDay.StartDay.Format(DateTimeFormat))
}

func startBreak() {
	if !inBreak {
		breakStart := time.Now().Round(0)
		workDay.Breaks = append(workDay.Breaks, Break{Start: breakStart})
		inBreak = true
		fmt.Println("Break started at", breakStart.Format(DateTimeFormat))
	} else {
		fmt.Println("You are already in a break.")
	}
}

func endBreak() {
	if inBreak {
		breakEnd := time.Now().Round(0)
		workDay.Breaks[len(workDay.Breaks)-1].End = breakEnd
		inBreak = false
		fmt.Println("Break ended at", breakEnd.Format(DateTimeFormat))
	} else {
		fmt.Println("You are not in a break.")
	}
}

func endDay() {
	workDay.EndDay = time.Now().Round(0)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a comment for the day (optional, max 500 characters): ")
	comment, _ := reader.ReadString('\n')
	comment = strings.TrimSpace(comment)
	if len(comment) > 500 {
		comment = comment[:500]
	}
	workDay.Comment = comment

	writeToCSV(workDay)
	fmt.Println("Work day ended at", workDay.EndDay.Format(DateTimeFormat))
}

func addFreeDay() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the date (DD.MM.YYYY) or date range (DD.MM.YYYY-DD.MM.YYYY) for the free day(s): ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = strings.TrimSpace(dateInput)

	var dates []string
	if strings.Contains(dateInput, "-") {
		dateRange := strings.Split(dateInput, "-")
		startDate, err1 := time.Parse(DateFormat, dateRange[0])
		endDate, err2 := time.Parse(DateFormat, dateRange[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid date range format.")
			return
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dates = append(dates, d.Format(DateFormat))
		}
	} else {
		date, err := time.Parse(DateFormat, dateInput)
		if err != nil {
			fmt.Println("Invalid date format.")
			return
		}
		dates = append(dates, date.Format(DateFormat))
	}

	fmt.Print("Enter a comment for the free day(s) (e.g., 'Vacation'): ")
	comment, _ := reader.ReadString('\n')
	comment = strings.TrimSpace(comment)
	if len(comment) > 500 {
		comment = comment[:500]
	}

	for _, date := range dates {
		workDay = WorkDay{
			StartDay: time.Now().Round(0),
			EndDay:   time.Now().Round(0),
			Comment:  comment,
		}
		writeFreeDayToCSV(date, workDay)
	}

	fmt.Println("Free day(s) added with comment:", comment)
}

func writeFreeDayToCSV(date string, day WorkDay) {
	file, err := os.OpenFile(CSVFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
		writer.Write([]string{"weekday", "date", "comment", "work hours", "net work hours", "start day", "end day", "break start", "break end"})
	}

	writer.Write([]string{
		time.Now().Weekday().String(),
		date,
		day.Comment,
		"0.00",
		"0.00",
		"",
		"",
		"",
		"",
	})
}

func takeOvertime() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the date (DD.MM.YYYY) for the overtime taken: ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = strings.TrimSpace(dateInput)

	_, err := time.Parse(DateFormat, dateInput)
	if err != nil {
		fmt.Println("Invalid date format.")
		return
	}

	fmt.Print("Enter the number of overtime hours taken (e.g., -8 for a full day): ")
	overtimeInput, _ := reader.ReadString('\n')
	overtimeInput = strings.TrimSpace(overtimeInput)
	overtimeHours, err := strconv.ParseFloat(overtimeInput, 64)
	if err != nil {
		fmt.Println("Invalid number format.")
		return
	}

	workDay = WorkDay{
		StartDay: time.Now().Round(0),
		EndDay:   time.Now().Round(0),
		Comment:  "Overtime taken",
	}
	writeOvertimeToCSV(dateInput, workDay, overtimeHours)

	fmt.Println("Overtime of", overtimeHours, "hours taken on", dateInput)
}

func writeOvertimeToCSV(date string, day WorkDay, overtimeHours float64) {
	file, err := os.OpenFile(CSVFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
		writer.Write([]string{"weekday", "date", "comment", "work hours", "net work hours", "start day", "end day", "break start", "break end"})
	}

	writer.Write([]string{
		time.Now().Weekday().String(),
		date,
		day.Comment,
		"0.00",
		formatFloat(overtimeHours),
		"",
		"",
		"",
		"",
	})
}

func currentOvertime() {
	file, err := os.Open(CSVFileName)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	if len(records) <= 1 {
		fmt.Println("No records to calculate overtime.")
		return
	}

	var totalOvertime float64
	for _, record := range records[1:] {
		if len(record) > 4 {
			netWorkHours, err := strconv.ParseFloat(strings.Replace(record[4], ",", ".", 1), 64)
			if err != nil {
				fmt.Println("Error parsing net work hours:", record[4])
				return
			}
			totalOvertime += netWorkHours
		}
	}

	fmt.Printf("Current total overtime: %.2f hours\n", totalOvertime)
}

func writeToCSV(day WorkDay) {
	file, err := os.OpenFile(CSVFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
		writer.Write([]string{"weekday", "date", "comment", "work hours", "net work hours", "start day", "end day", "break start", "break end"})
	}

	var breaks []string
	for _, b := range day.Breaks {
		breaks = append(breaks, b.Start.Format(DateTimeFormat), b.End.Format(DateTimeFormat))
	}

	workHours := workDay.EndDay.Sub(workDay.StartDay).Hours()
	netWorkHours := workHours
	for _, b := range day.Breaks {
		netWorkHours -= b.End.Sub(b.Start).Hours()
	}

	writer.Write(append([]string{
		day.StartDay.Weekday().String(),
		day.StartDay.Format(DateFormat),
		day.Comment,
		formatFloat(workHours),
		formatFloat(netWorkHours),
		day.StartDay.Format(DateTimeFormat),
		day.EndDay.Format(DateTimeFormat),
	}, breaks...))
}

func formatFloat(f float64) string {
	return strings.Replace(fmt.Sprintf("%.2f", f), ".", ",", 1)
}

func (data CSV) Less(i, j int) bool {
	dateColumnIndex := 1
	date1 := data[i][dateColumnIndex]
	date2 := data[j][dateColumnIndex]
	timeT1, _ := time.Parse(DateFormat, date1)
	timeT2, _ := time.Parse(DateFormat, date2)

	return timeT1.Before(timeT2)
}

func (data CSV) Len() int {
	return len(data)
}

func (data CSV) Swap(i, j int) {
	data[i], data[j] = data[j], data[i]
}

func sortCSVByDate() {
	file, err := os.Open(CSVFileName)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	if len(records) <= 1 {
		return // No need to sort if there are no records or only the header
	}

	header := records[0]
	data := records[1:]

	// Sort the data based on the second column (date)
	sort.Sort(CSV(data))

	// Create and open the output file
	outputFile, err := os.Create(CSVFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header and sorted data records
	writer.Write(header)
	writer.WriteAll(data)
}
