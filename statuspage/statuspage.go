package statuspage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"main/buildinfo"
	"main/config"
	"net/http"
	"strings"
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

func GetDataShort(business_day string) (*StatusJson, error) {
	Host := getHost()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/json", Host, business_day), nil)
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

	status := &StatusJson{}
	json.NewDecoder(resp.Body).Decode(status)
	return status, nil
}
