package client

import "os"

// DNSState represents a DNS zone configuration
type DNSState struct {
	statePath string
	records   []Record
}

func loadState(statePath string) (*DNSState, error) {

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		records := []Record{}
		saveRecords(statePath, records)
	}

	records, err := loadRecords(statePath)
	if err != nil {
		return nil, err
	}

	return &DNSState{
		statePath: statePath,
		records:   records,
	}, nil
}

func (s *DNSState) save() error {
	return saveRecords(s.statePath, s.records)
}
