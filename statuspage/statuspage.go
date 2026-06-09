package statuspage

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/buildinfo"
	"net/http"
	"os"
	"strings"
)

type StatusTables struct {
	FBCC [][]string
	DAMC [][]string
	IDMC [][]string
}

func GetData(business_day string) (StatusTables, bool, error) {
	var Host string
	if Debug := os.Getenv("AMUN_DEBUG"); Debug == "1" {
		Host = "http://localhost:5000"
	} else {
		Host = "https://status.coreflowbased.eu"
	}
	tables := StatusTables{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/table", Host, business_day), nil)
	if err != nil {
		return tables, false, err
	}
	req.Header.Add("accept", "text/csv")
	resp, err := client.Do(req)
	if err != nil {
		return tables, false, err
	}

	defer resp.Body.Close()

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
