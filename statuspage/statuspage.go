package statuspage

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StatusTables struct {
	FBCC [][]string
	DAMC [][]string
	IDMC [][]string
}

func GetData(business_day string) (StatusTables, error) {
	tables := StatusTables{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://status.coreflowbased.eu/%s/table", business_day), nil)
	if err != nil {
		return tables, err
	}
	req.Header.Add("accept", "text/csv")
	resp, err := client.Do(req)
	if err != nil {
		return tables, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return tables, err
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
