package main

import (
	"fmt"
	"log"
	"main/statuspage"
	"time"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

const (
	purple    = "99"
	gray      = "245"
	lightGray = "241"
)

var (
	// HeaderStyle is the lipgloss style used for the table headers.
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(purple)).
			Bold(true).
			Align(lipgloss.Left)
	// CellStyle is the base lipgloss style used for the table rows.
	CellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(14)
	// OddRowStyle is the lipgloss style used for odd-numbered table rows.
	OddRowStyle = CellStyle.Foreground(lipgloss.Color(gray))
	// EvenRowStyle is the lipgloss style used for even-numbered table rows.
	EvenRowStyle = CellStyle.Foreground(lipgloss.Color(lightGray))
	// BorderStyle is the lipgloss style used for the table border.
	BorderStyle = lipgloss.NewStyle().Foreground(lipgloss.White)
)

func askDate() (string, error) {
	validate := func(s string) error {
		_, err := time.Parse("2006-01-02", s)
		if err != nil {
			return fmt.Errorf("use YYYY-MM-DD")
		}
		return nil
	}
	var business_day string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("From date").
				Description("Please specify the business day to query").
				Placeholder("YYYY-MM-DD").
				Validate(validate).
				Value(&business_day),
		),
	).Run()

	return business_day, err
}

func printTable(table_data [][]string) {
	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			style = style.Align(lipgloss.Left)

			return style
		}).
		Headers(table_data[0]...).
		Rows(table_data[1:]...)

	lipgloss.Println(t)
}

func main() {

	business_day, err := askDate()
	if err != nil {
		log.Fatal(err)
	}

	tables, err := statuspage.GetData(business_day)
	if err != nil {
		log.Fatal(err)
	}

	lipgloss.Println("Flowbased Capacity Calculation")
	printTable(tables.FBCC)

	lipgloss.Println("Day-Ahead Market Coupling")
	printTable(tables.DAMC)

	lipgloss.Println("Intraday Market Coupling")
	printTable(tables.IDMC)

}
