package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
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

const (
	DateFormat       = "02.01.2006"
	TimeFormat       = "15:04:05"
	DateTimeFormat   = "02.01.2006 15:04:05"
	HoursInWorkDay   = 8.0
	CSVFileName      = "worktime.csv"
	OvertimeFileName = "overtime.txt"
)

var workDay WorkDay
var inBreak bool

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Available commands:")
	fmt.Println("  start day  - Start the work day")
	fmt.Println("  break start - Start a break")
	fmt.Println("  break end  - End a break")
	fmt.Println("  end day    - End the work day")

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
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  start day  - Start the work day")
			fmt.Println("  break start - Start a break")
			fmt.Println("  break end  - End a break")
			fmt.Println("  end day    - End the work day")
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
	updateOvertime(workDay)
	fmt.Println("Work day ended at", workDay.EndDay.Format(DateTimeFormat))
}

func writeToCSV(day WorkDay) {
	file, err := os.OpenFile(CSVFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Set semicolon as the delimiter
	defer writer.Flush()

	if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
		writer.Write([]string{"date", "comment", "work hours", "net work hours", "start day", "end day", "break start", "break end"})
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

func updateOvertime(day WorkDay) {
	netWorkHours := day.EndDay.Sub(day.StartDay).Hours()
	for _, b := range day.Breaks {
		netWorkHours -= b.End.Sub(b.Start).Hours()
	}

	overtime := netWorkHours - HoursInWorkDay
	currentOvertime := readCurrentOvertime()
	newOvertime := currentOvertime + overtime

	file, err := os.OpenFile(OvertimeFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening overtime file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Current Overtime: %.2f hours\n", newOvertime))
	if err != nil {
		fmt.Println("Error writing to overtime file:", err)
	}
}

func readCurrentOvertime() float64 {
	file, err := os.OpenFile(OvertimeFileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening overtime file:", err)
		return 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentOvertime float64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Current Overtime:") {
			fmt.Sscanf(line, "Current Overtime: %f hours", &currentOvertime)
		}
	}
	return currentOvertime
}
