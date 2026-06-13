package main

import (
	"fmt"
	"log"
	"main/buildinfo"
	"main/config"
	"main/statuspage"
	"main/telemetry"
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/charmbracelet/x/term"
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

func validateDate(s string) error {
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("use YYYY-MM-DD")
	}
	return nil
}

func askDate() (string, error) {
	var business_day string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("From date").
				Description("Please specify the business day to query").
				Placeholder("YYYY-MM-DD").
				Validate(validateDate).
				Value(&business_day),
		),
	).Run()

	return business_day, err
}

func printTable(table_data [][]string) {
	termWidth, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || termWidth == 0 {
		termWidth = 125
	}

	termWidth -= 5
	firstWidth := termWidth * 15 / 100
	columnWidth := (termWidth - firstWidth) / (len(table_data[0]) - 1)

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
				style = style.Bold(true).Foreground(lipgloss.Color("#FACD81")).
					Width(firstWidth)
			} else {
				style = style.Width(columnWidth)
				if len(table_data[row+1][col]) >= 3 {
					cell_data := table_data[row+1][col]
					if strings.Contains(cell_data, "[✖]") {
						style = style.Foreground(lipgloss.Color(Red))
					} else if strings.Contains(cell_data, "[!]") {
						style = style.Foreground(lipgloss.Color(Yellow))
					} else if strings.Contains(cell_data, "[✔]") {
						style = style.Foreground(lipgloss.Color(Green))
					}
				}
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

func printConfig(c config.Config) {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderStyle
			} else {
				return lipgloss.NewStyle().
					Padding(0, 1).
					Width(25).Foreground(lipgloss.White)
			}
		}).
		Headers("Setting", "Value").
		Row("disable_telemetry", strconv.FormatBool(c.General.DisableTelemetry)).
		Row("debug_backend", strconv.FormatBool(c.FBStatus.Debug)).
		Row("surpress_version_check", strconv.FormatBool(c.FBStatus.SuppressVersionCheck))

	lipgloss.Println(t)
}

func DefaultTable(business_day string) {
	tables, new_version_available, err := statuspage.GetData(business_day)
	if err != nil {
		log.Fatal(err)
	}

	if new_version_available && !config.GetConfig().FBStatus.SuppressVersionCheck {
		lipgloss.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(Red)).Render(">> New version of this cli is available!"))
	}

	lipgloss.Println(HeaderStyle.Render(">> Core Market Coupling Status\n" +
		fmt.Sprintf(">> For business day %s \n", business_day) +
		">> An Amun Analytics product\n"))

	lipgloss.Println(HeaderStyle.Render("> Flowbased Capacity Calculation"))
	printTable(tables.FBCC)

	lipgloss.Println(HeaderStyle.Render("> Day-Ahead Market Coupling"))
	printTable(tables.DAMC)

	lipgloss.Println(HeaderStyle.Render("> Intraday Market Coupling"))
	printTable(tables.IDMC)
}

func ShortTable(business_day string) {
	data, err := statuspage.GetDataShort(business_day)
	if err != nil {
		log.Fatal(err)
	}

	table_data := statuspage.GetTableDataShort(data)

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle()

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
	business_day := ""
	var err error
	short_table := false

	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			lipgloss.Println(HeaderStyle.Render(fmt.Sprintf("Version: \t\t %s", buildinfo.Version)))
			lipgloss.Println(HeaderStyle.Render(fmt.Sprintf("Git Commit: \t %s", buildinfo.GitCommit)))
			return
		} else if os.Args[1] == "config" {
			c := config.GetConfig()
			printConfig(c)
			return
		} else if os.Args[1] == "short" {
			short_table = true
			if len(os.Args) > 2 {
				business_day = os.Args[2]
			}
		} else {
			business_day = os.Args[1]
		}
	}

	if business_day == "" {
		business_day, err = askDate()
		if err != nil {
			log.Fatal(err)
		}
	} else if business_day == "today" || business_day == "D" {
		business_day = time.Now().Format("2006-01-02")
	} else if business_day == "tomorrow" || business_day == "D+1" {
		business_day = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	} else {
		if validateDate(business_day) != nil {
			log.Fatal("Please use format YYYY-MM-DD")
		}
	}

	telemetry.Register("fbstatuscli", buildinfo.Version)

	if !short_table {
		DefaultTable(business_day)
	} else {
		ShortTable(business_day)
	}

}
