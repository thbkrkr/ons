package dns

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ovh/go-ovh/ovh"
)

// OvhClient is a wrapper of an OVH API Client
type OvhClient struct {
	Client *ovh.Client
}

// Config for the OVH client
type Config struct {
	Zone              string `json:"zone"`
	Endpoint          string `json:"endpoint"`
	ApplicationKey    string `json:"ak"`
	ApplicationSecret string `json:"as"`
	ConsumerKey       string `json:"ck"`
}

type RecordState struct {
	Target    string `json:"target"`
	SubDomain string `json:"subDomain"`
	State     string `json:"state,omitempty"`
	ID        int64  `json:"id,omitempty"`
}

func (c *OvhClient) LoadState(stateFilename string) ([]*RecordState, error) {
	stateFile, err := os.Open(stateFilename)
	if err != nil {
		return nil, err
	}
	defer stateFile.Close()

	var state []*RecordState
	jsonParser := json.NewDecoder(stateFile)
	if err = jsonParser.Decode(&state); err != nil {
		return nil, err
	}

	return state, nil
}

func (c *OvhClient) Plan(recordStates []*RecordState, records []*Record) []*RecordState {
	for _, recordState := range recordStates {
		_, state := c.GetState(recordState, records)
		//recordState.ID = id
		recordState.State = state
		//fmt.Printf("%v\n", recordState)
	}
	return recordStates
}

func (c *OvhClient) Apply(zone string, recordStates []*RecordState, records []*Record) ([]*RecordState, error) {
	var modifications []*RecordState

	for _, recordState := range recordStates {
		_, state := c.GetState(recordState, records)
		if state == "to_add" {
			fmt.Printf("Add record: %s %s...\n", recordState.SubDomain, recordState.Target)
			_, err := c.AddRecord(zone, recordState.SubDomain, recordState.Target)
			if err != nil {
				return nil, err
			}
			recordState.State = "added"
			modifications = append(modifications, recordState)
		}
	}
	c.ApplyZoneModifications(zone)
	return modifications, nil
}

func (c *OvhClient) GetState(recordState *RecordState, records []*Record) (int64, string) {
	for _, record := range records {
		if recordState.SubDomain == record.SubDomain && recordState.Target == record.Target {
			return record.ID, "ok"
		}
	}
	return 0, "to_add"
}

func (c *OvhClient) ApplyZoneModifications(zone string) error {
	err := c.Client.Post(fmt.Sprintf("/domain/zone/%s/refresh", zone), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *OvhClient) ListRecords(zone string, fieldType string) ([]int64, error) {
	var records []int64

	err := c.Client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=%s", zone, fieldType), &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *OvhClient) ListFullARecords(zone string) ([]*Record, error) {
	records, err := c.ListRecords(zone, "A")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(records))

	var fullRecords []*Record
	for _, recordID := range records {
		go func(recordID int64) {
			defer wg.Done()
			record, err := c.GetRecord(zone, recordID)
			if err != nil {
				return
			}
			fullRecords = append(fullRecords, record)
		}(recordID)
	}
	wg.Wait()

	return fullRecords, nil
}

// Record represents a DNS record
type Record struct {
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
	Zone      string `json:"zone"`
	FieldType string `json:"fieldType"`
	ID        int64  `json:"id"`
	SubDomain string `json:"subDomain"`
}

func getRecordByID(id int64, records []*Record) (*Record, error) {
	for _, record := range records {
		if record.ID == id {
			return record, nil
		}
	}

	return nil, fmt.Errorf("No record found for id: %d", id)
}

func (c *OvhClient) GetRecord(zone string, id int64) (*Record, error) {
	var record = &Record{}

	err := c.Client.Get(fmt.Sprintf("/domain/zone/%s/record/%d", zone, id), record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (c *OvhClient) GetRecordIDBySubDomain(zone string, subDomain string) (*int64, error) {
	var records []int64

	err := c.Client.Get(fmt.Sprintf("/domain/zone/%s/record/?subDomain=%s", zone, subDomain), records)
	if err != nil {
		return nil, err
	}

	return &records[0], nil
}

// AddRecord represents the request to add a new record
type AddRecord struct {
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
}

func (c *OvhClient) AddRecord(zone string, subDomain string, target string) (*Record, error) {
	var record = &Record{}

	newRecord := &AddRecord{FieldType: "A", SubDomain: subDomain, Target: target}
	err := c.Client.Post(fmt.Sprintf("/domain/zone/%s/record", zone), newRecord, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (c *OvhClient) DeleteRecordByID(zone string, id int64) (bool, error) {
	var record = &Record{}

	err := c.Client.Delete(fmt.Sprintf("/domain/zone/%s/record/%d", zone, id), record)
	if err != nil {
		return false, err
	}

	return true, nil
}
