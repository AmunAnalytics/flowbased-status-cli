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

type PrintMode int

const (
	BasePrintMode PrintMode = iota
	IvaPrintMode
)

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

func doApiCall(url string) (*http.Response, error) {
	Host := getHost()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", Host, url), nil)
	setUserAgent(req)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	if IsNewVersion(resp) {
		config.NewVersionAvailable = true
	}

	return resp, nil
}

func IsNewVersion(resp *http.Response) bool {
	return resp.Header.Get("x-cli-outdated") == "1"
}

func GetData(business_day string) (StatusTables, error) {
	Host := getHost()
	tables := StatusTables{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/table", Host, business_day), nil)
	if err != nil {
		return tables, err
	}
	setUserAgent(req)
	req.Header.Add("accept", "text/csv")
	resp, err := client.Do(req)
	if err != nil {
		return tables, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return tables, fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return tables, err
	}

	if IsNewVersion(resp) {
		config.NewVersionAvailable = true
	}

	for i, table := range strings.Split(string(bodyBytes), "\n\n") {
		r := csv.NewReader(strings.NewReader(table))
		t, errt := r.ReadAll()
		if err != nil {
			return tables, errt
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

	return tables, nil

}

func PrintTable(table_data [][]string, offsetColumns int, p PrintMode) {
	termWidth, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || termWidth == 0 {
		termWidth = 125
	}

	termWidth -= 5
	firstWidth := termWidth * 15 / 100
	columnWidth := (termWidth - firstWidth) / (len(table_data[0]) - 1 + offsetColumns)

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
				cell_data := table_data[row+1][col]
				if p == BasePrintMode {
					if len(table_data[row+1][col]) >= 3 {
						if strings.Contains(cell_data, "[✖]") {
							style = style.Foreground(lipgloss.Color(config.Red))
						} else if strings.Contains(cell_data, "[!]") {
							style = style.Foreground(lipgloss.Color(config.Yellow))
						} else if strings.Contains(cell_data, "[✔]") {
							style = style.Foreground(lipgloss.Color(config.Green))
						}
					}
				} else if p == IvaPrintMode {
					if cell_data[0] == 'F' {
						style = style.Foreground(lipgloss.Color(config.Red))
					} else if cell_data == "ok" {
						style = style.Foreground(lipgloss.Color(config.Green))
					} else {
						style = style.Foreground(lipgloss.Color(config.Yellow))
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
