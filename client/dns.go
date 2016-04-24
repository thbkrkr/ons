package client

import (
	"fmt"
	"sort"
	"sync"
)

// ListARecords lists all A DNS zone records
func (c *OnsClient) ListARecords(zone string) (Records, error) {
	records, err := c.ListRecordsByType(zone, "A")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(records))

	var fullRecords Records
	for _, recordID := range records {
		go func(recordID int64) {
			defer wg.Done()
			record, err := c.GetRecordByID(zone, recordID)
			if err != nil {
				return
			}
			fullRecords = append(fullRecords, *record)
		}(recordID)
	}
	wg.Wait()

	sort.Sort(fullRecords)

	return fullRecords, nil
}

// ListRecordsByType lists all DNS zone records given a type (A, MX, SRV, NS, ...)
func (c *OnsClient) ListRecordsByType(zone string, fieldType string) ([]int64, error) {
	var records []int64

	err := c.client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=%s", zone, fieldType), &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordByID gets a DNS zone record by its ID
func (c *OnsClient) GetRecordByID(zone string, id int64) (*Record, error) {
	var record = &Record{}

	err := c.client.Get(fmt.Sprintf("/domain/zone/%s/record/%d", zone, id), record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetRecordIDBySubDomain gets a DNS zone record ID given a record subdomain
func (c *OnsClient) GetRecordIDBySubDomain(zone string, subDomain string) (*int64, error) {
	var records []int64

	err := c.client.Get(fmt.Sprintf("/domain/zone/%s/record?subDomain=%s", zone, subDomain), records)
	if err != nil {
		return nil, err
	}

	if len(records) != 1 {
		return nil, nil
	}

	return &records[0], nil
}

// addRecord represents the request to add a new DNS zone record
type addRecord struct {
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
}

// AddRecord create a new DNS zone record
func (c *OnsClient) AddRecord(zone string, subDomain string, target string) (*Record, error) {
	var record = &Record{}

	newRecord := &addRecord{FieldType: "A", SubDomain: subDomain, Target: target}
	err := c.client.Post(fmt.Sprintf("/domain/zone/%s/record", zone), newRecord, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// DeleteRecordBySubDomain deletes a DNS zone record given a record subdomain
/*func (c *OnsClient) DeleteRecordBySubDomain(zone string, subDomain string) (bool, error) {
	id, err := c.GetRecordIDBySubDomain(zone, subDomain)
	if err != nil {
		return false, err
	}

	return c.DeleteRecordByID(zone, *id)
}*/

// DeleteRecordByID deletes a DNS zone record given a record ID
func (c *OnsClient) DeleteRecordByID(zone string, id int64) (bool, error) {
	var record = &Record{}

	err := c.client.Delete(fmt.Sprintf("/domain/zone/%s/record/%d", zone, id), record)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RefreshZone applies the DNS zone configuration to DNS servers
func (c *OnsClient) RefreshZone(zone string) error {
	err := c.client.Post(fmt.Sprintf("/domain/zone/%s/refresh", zone), nil, nil)
	if err != nil {
		return err
	}

	return nil
}
