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
	"time"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
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

func printConfig(c config.Config) {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return config.HeaderStyle
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
		lipgloss.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(config.Red)).Render(">> New version of this cli is available!"))
	}

	lipgloss.Println(config.HeaderStyle.Render(">> Core Market Coupling Status\n" +
		fmt.Sprintf(">> For business day %s \n", business_day) +
		">> An Amun Analytics product\n"))

	lipgloss.Println(config.HeaderStyle.Render("> Flowbased Capacity Calculation"))
	statuspage.PrintTable(tables.FBCC)

	lipgloss.Println(config.HeaderStyle.Render("> Day-Ahead Market Coupling"))
	statuspage.PrintTable(tables.DAMC)

	lipgloss.Println(config.HeaderStyle.Render("> Intraday Market Coupling"))
	statuspage.PrintTable(tables.IDMC)
}

func main() {
	business_day := ""
	var err error
	short_table := false

	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			lipgloss.Println(config.HeaderStyle.Render(fmt.Sprintf("Version: \t\t %s", buildinfo.Version)))
			lipgloss.Println(config.HeaderStyle.Render(fmt.Sprintf("Git Commit: \t %s", buildinfo.GitCommit)))
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
		statuspage.PrintShortTable(business_day)
	}

}
