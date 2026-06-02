package main

import (
	"fmt"
	"log"
	"main/statuspage"
	"os"
	"time"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var (
	CellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(14)
	RowStyle = CellStyle.Foreground(lipgloss.White).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true)
	BorderStyle = lipgloss.NewStyle().Foreground(lipgloss.White)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f59e0b")).
			Bold(true).
			Align(lipgloss.Left)
)

var (
	Green  = "#65a30d"
	Yellow = "#fbbf24"
	Red    = "#b91c1c"
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
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row == table.HeaderRow {
				style = HeaderStyle
			} else {
				style = RowStyle
			}

			if col == 0 {
				style = style.Bold(true).Foreground(lipgloss.Color("#FACD81"))
			}

			if row == len(table_data)-2 {
				style = style.BorderBottom(false)
			}

			style = style.Align(lipgloss.Left)

			return style
		}).
		Headers(table_data[0]...).
		Rows(table_data[1:]...)

	lipgloss.Println(t)
}

func main() {
	var business_day string
	var err error

	if len(os.Args) > 1 {
		business_day = os.Args[1]
	} else {
		business_day, err = askDate()
		if err != nil {
			log.Fatal(err)
		}
	}
	tables, err := statuspage.GetData(business_day)
	if err != nil {
		log.Fatal(err)
	}

	lipgloss.Println(HeaderStyle.Render("Core Market Coupling Status\n" +
		fmt.Sprintf("For business day %s \n", business_day) +
		"An Amun Analytics product\n"))

	lipgloss.Println(HeaderStyle.Render("> Flowbased Capacity Calculation"))
	printTable(tables.FBCC)

	lipgloss.Println(HeaderStyle.Render("> Day-Ahead Market Coupling"))
	printTable(tables.DAMC)

	lipgloss.Println(HeaderStyle.Render("> Intraday Market Coupling"))
	printTable(tables.IDMC)

}
