package statuspage

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/buildinfo"
	"main/config"
	"net/http"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/charmbracelet/x/term"
)

type StatusTables struct {
	FBCC [][]string
	DAMC [][]string
	IDMC [][]string
}

func getHost() string {
	if config.GetConfig().FBStatus.Debug {
		return "http://localhost:5000"
	} else {
		return "https://status.coreflowbased.eu"
	}
}

func setUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", fmt.Sprintf("fbstatus cli %s (%s)", buildinfo.Version, buildinfo.GitCommit))
}

func GetData(business_day string) (StatusTables, bool, error) {
	Host := getHost()
	tables := StatusTables{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/table", Host, business_day), nil)
	if err != nil {
		return tables, false, err
	}
	setUserAgent(req)
	req.Header.Add("accept", "text/csv")
	resp, err := client.Do(req)
	if err != nil {
		return tables, false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return tables, false, fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return tables, false, err
	}

	new_version_available := resp.Header.Get("x-cli-version") > buildinfo.Version

	for i, table := range strings.Split(string(bodyBytes), "\n\n") {
		r := csv.NewReader(strings.NewReader(table))
		t, errt := r.ReadAll()
		if err != nil {
			return tables, new_version_available, errt
		}
		switch i {
		case 0:
			tables.FBCC = t
		case 1:
			tables.DAMC = t
		case 2:
			tables.IDMC = t

		}
	}

	return tables, new_version_available, nil

}

func PrintTable(table_data [][]string) {
	termWidth, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || termWidth == 0 {
		termWidth = 125
	}

	termWidth -= 5
	firstWidth := termWidth * 15 / 100
	columnWidth := (termWidth - firstWidth) / (len(table_data[0]) - 1)

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(config.BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row == table.HeaderRow {
				style = config.HeaderStyle
			} else {
				style = config.RowStyle
			}

			if col == 0 {
				style = style.Bold(true).Foreground(lipgloss.Color("#FACD81")).
					Width(firstWidth)
			} else {
				style = style.Width(columnWidth)
				if len(table_data[row+1][col]) >= 3 {
					cell_data := table_data[row+1][col]
					if strings.Contains(cell_data, "[✖]") {
						style = style.Foreground(lipgloss.Color(config.Red))
					} else if strings.Contains(cell_data, "[!]") {
						style = style.Foreground(lipgloss.Color(config.Yellow))
					} else if strings.Contains(cell_data, "[✔]") {
						style = style.Foreground(lipgloss.Color(config.Green))
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
