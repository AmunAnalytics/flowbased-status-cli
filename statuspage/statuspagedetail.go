package statuspage

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type StatusIvaDetailJson map[string]map[string]TSOMetrics

type TSOMetrics struct {
	AvgIVAPct float64 `json:"avg_iva_pct"`
	AvgRAMPct float64 `json:"avg_ram_pct"`
	Fallback  bool    `json:"fallback"`
	NumCNECs  int     `json:"num_cnecs"`
}

var TSOs = []string{
	"50HERTZ", "AMPRION", "APG", "CEPS", "ELES", "ELIA", "HOPS",
	"MAVIR", "PSE", "RTE", "SEPS", "TENNETBV", "TENNETGMBH",
	"TRANSELECTRICA", "TRANSNETBW",
}

func GetDataIvaDetail(business_day string, timeframe string) (*StatusIvaDetailJson, error) {
	resp, err := doApiCall(fmt.Sprintf("%s/details/%s/iva/json", business_day, timeframe))
	if err != nil {
		return nil, err
	}
	status := &StatusIvaDetailJson{}
	json.NewDecoder(resp.Body).Decode(status)
	return status, nil
}

func GetTableIvaDetail(data *StatusIvaDetailJson) [][]string {
	timestamps := make([]string, 0, len(*data))
	for ts := range *data {
		timestamps = append(timestamps, ts)
	}
	sort.Strings(timestamps)

	table := make([][]string, 0, len(timestamps)+1)
	header := make([]string, 0, len(TSOs)+1)
	header = append(header, "MTU")
	header = append(header, TSOs...)
	table = append(table, header)

	for _, ts := range timestamps {
		row_data := (*data)[ts]
		row := make([]string, 0, len(TSOs)+1)
		row = append(row, ts)
		for _, tso := range TSOs {
			metrics := row_data[tso]
			if metrics.NumCNECs == 0 {
				row = append(row, "ok")
			} else if metrics.Fallback {
				row = append(row, fmt.Sprintf("F %d@%.2f%%", metrics.NumCNECs, metrics.AvgIVAPct))
			} else {
				row = append(row, fmt.Sprintf("%d@%.2f%%", metrics.NumCNECs, metrics.AvgIVAPct))
			}
		}
		table = append(table, row)
	}

	return table
}

func PrintTableIvaDetail(business_day string) {
	data, err := GetDataIvaDetail(business_day, "DACC")
	if err != nil {
		log.Fatal(err)
	}

	table := GetTableIvaDetail(data)

	PrintTable(table, 1, IvaPrintMode)
}
