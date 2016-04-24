package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Record represents a DNS zone record
type Record struct {
	Zone      string `json:"zone"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`

	ID        int64  `json:"id,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	FieldType string `json:"fieldType,omitempty"`

	Managed string `json:"-"`
}

// GetBySubDomainAndTarget gets a record from a list of records  by comparing records
// zone, sub domain and target
func (r Record) GetBySubDomainAndTarget(records Records) *Record {
	for _, re := range records {
		if re.Zone == r.Zone && re.SubDomain == r.SubDomain && re.Target == r.Target {
			return &re
		}
	}
	return nil
}

// ExistsInBySubDomain returns true if a record exists in a list of records
// by comparing records zone and sub domain
func (r Record) ExistsInBySubDomain(records Records) bool {
	for _, re := range records {
		if re.Zone == r.Zone && re.SubDomain == r.SubDomain {
			return true
		}
	}
	return false
}

// ExistsInBySubDomainAndTarget returns true if a record exists in a list of records
// by comparing records zone, sub domain and target
func (r Record) ExistsInBySubDomainAndTarget(records Records) bool {
	record := r.GetBySubDomainAndTarget(records)
	return record != nil
}

// ExistsInByID returns true if a record exists in a list of records
// by comparing records ids
func (r Record) ExistsInByID(records Records) bool {
	for _, re := range records {
		if r.ID == re.ID {
			return true
		}
	}
	return false
}

// Print prints a record with fixed indentation and colors
func (r Record) Print() {
	fmt.Printf("%-30s %-1s %s\n", magenta(r.Target), r.Managed, green(r.SubDomain+"."+r.Zone))
}

// Records represents a list of DNS zone record
type Records []Record

// loadRecords loads records from a file in JSON format
func loadRecords(filepath string) ([]Record, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// saveRecords loads records in a file in JSON format
func saveRecords(filepath string, records []Record) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Records is sortable

func (r Records) Len() int {
	return len(r)
}

func (r Records) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Records) Less(i, j int) bool {
	return r[i].SubDomain < r[j].SubDomain
}
