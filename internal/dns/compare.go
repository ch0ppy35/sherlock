package dns

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// CompareRecords compares expected and actual DNS records, printing the results in a formatted table and returning an error if mismatches are found
func CompareRecords(expected []string, actual []string) error {
	expectedMap := make(map[string]struct{}, len(expected))
	for _, val := range expected {
		expectedMap[val] = struct{}{}
	}

	actualMap := make(map[string]struct{}, len(actual))
	for _, val := range actual {
		actualMap[val] = struct{}{}
	}

	matchedRecords := []string{}
	unexpectedRecords := []string{}
	missingRecords := []string{}

	for _, val := range actual {
		if _, found := expectedMap[val]; !found {
			unexpectedRecords = append(unexpectedRecords, val)
		} else {
			matchedRecords = append(matchedRecords, val)
		}
	}

	for _, val := range expected {
		if _, found := actualMap[val]; !found {
			missingRecords = append(missingRecords, val)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	if len(matchedRecords) > 0 {
		t.AppendHeader(table.Row{"Type", "Record"})
		t.AppendRows([]table.Row{
			{"Matched", ""},
		})
		for _, record := range matchedRecords {
			t.AppendRows([]table.Row{
				{"", record},
			})
		}
	} else {
		t.AppendHeader(table.Row{"Type", "Record"})
		t.AppendRows([]table.Row{
			{"Matched", "None Found"},
		})
	}

	if len(unexpectedRecords) > 0 {
		t.AppendRows([]table.Row{
			{"Unexpected", ""},
		})
		for _, record := range unexpectedRecords {
			t.AppendRows([]table.Row{
				{"", record},
			})
		}
	} else {
		t.AppendRows([]table.Row{
			{"Unexpected", "None"},
		})
	}

	if len(missingRecords) > 0 {
		t.AppendRows([]table.Row{
			{"Missing", ""},
		})
		for _, record := range missingRecords {
			t.AppendRows([]table.Row{
				{"", record},
			})
		}
	} else {
		t.AppendRows([]table.Row{
			{"Missing", "None"},
		})
	}
	t.Render()

	if len(unexpectedRecords) > 0 || len(missingRecords) > 0 {
		return fmt.Errorf("mismatched records found")
	}
	return nil
}
