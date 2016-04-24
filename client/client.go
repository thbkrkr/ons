package client

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ovh/go-ovh/ovh"
)

var (
	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	cyan    = color.New(color.FgCyan).PrintfFunc()
)

// OnsClient is a wrapper of an OVH API Client, a state and a config
type OnsClient struct {
	client     *ovh.Client
	config     *DNSConfig
	configPath string
	state      *DNSState
	statePath  string
}

// NewOnsClient creates a new ONS client
func NewOnsClient(statePath string, configPath string,
	endpoint string, ak string, as string, ck string) (*OnsClient, error) {

	ovhClient, err := ovh.NewClient(endpoint, ak, as, ck)
	if err != nil {
		return nil, err
	}

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	state, err := loadState(statePath)
	if err != nil {
		return nil, err
	}

	return &OnsClient{
		client:     ovhClient,
		configPath: configPath,
		config:     config,
		statePath:  statePath,
		state:      state,
	}, nil
}

// --

// Ls lists all records from a DNS zone by marking configured record with a star
func (c *OnsClient) Ls(zone string) (Records, error) {
	records, err := c.ListARecords(zone)
	if err != nil {
		return nil, err
	}

	for i, r := range records {
		if r.ExistsInBySubDomainAndTarget(c.config.records) {
			r.Managed = "*"
			records[i] = r
		}
	}

	return records, nil
}

// Add adds a new record in the config and plans the DNS config
func (c *OnsClient) Add(zone string, subDomain string, target string) error {
	record := Record{Zone: zone, SubDomain: subDomain, Target: target}

	if record.ExistsInBySubDomainAndTarget(c.config.records) {
		return fmt.Errorf("Record `%s.%s %s` already added", record.SubDomain, record.Zone, record.Target)
		//return nil
	}

	c.config.records = append(c.config.records, record)

	err := c.config.save()
	if err != nil {
		return err
	}

	return nil
}

// Rm removes records from the config given a sub domain and plans the DNS config
func (c *OnsClient) Rm(zone string, subDomain string) error {
	record := Record{Zone: zone, SubDomain: subDomain}

	if !record.ExistsInBySubDomain(c.state.records) && !record.ExistsInBySubDomain(c.config.records) {
		return fmt.Errorf("Record `%s.%s` not managed.", record.SubDomain, record.Zone)
	}

	for i := len(c.config.records) - 1; i >= 0; i-- {
		r := c.config.records[i]
		if r.Zone == zone && r.SubDomain == subDomain {
			c.config.records = append(c.config.records[:i],
				c.config.records[i+1:]...)
		}
	}

	err := c.config.save()
	if err != nil {
		return err
	}

	return nil
}

// Plan shows the DNS zone modifications to apply
func (c *OnsClient) Plan(zone string) ([]Record, []Record, error) {
	var toAdd []Record
	var toRm []Record

	dns, err := c.ListARecords(zone)
	if err != nil {
		return nil, nil, err
	}

	touchState := false

	// Plan to addd record if it exists in the config
	for _, r := range c.config.records {

		// and not in the dns zone
		if !r.ExistsInBySubDomainAndTarget(dns) {

			toAdd = append(toAdd, r)

		} else
		// and in the dns zone but not in the state
		if !r.ExistsInBySubDomainAndTarget(c.state.records) &&
			r.ExistsInBySubDomainAndTarget(dns) {

			record := r.GetBySubDomainAndTarget(dns)
			c.state.records = append(c.state.records, *record)
			touchState = true
		}
	}

	// Remove record from the state
	for i, r := range c.state.records {

		// if not in the DNS zone
		record := r.GetBySubDomainAndTarget(dns)
		if record == nil {
			c.state.records = append(c.state.records[:i], c.state.records[i+1:]...)
			touchState = true

		} else
		// Refresh record if ID is absent
		if r.ID == 0 {
			c.state.records[i] = *record
			touchState = true
		}

		// Plan to remove record if it exists in the state
		// and not in the config but in the dns zone
		if !r.ExistsInBySubDomainAndTarget(c.config.records) &&
			r.ExistsInBySubDomainAndTarget(dns) {

			toRm = append(toRm, r)
		}
	}

	if touchState {
		err = c.state.save()
		if err != nil {
			return nil, nil, err
		}
	}

	return toAdd, toRm, nil
}

var (
	printAdd = color.New(color.Bold, color.FgGreen).PrintfFunc()
	printRm  = color.New(color.Bold, color.FgRed).PrintfFunc()
)

// Apply applies the zone DNS configuration on the DNS zone
func (c *OnsClient) Apply(zone string) (int, int, error) {
	added := 0
	removed := 0

	toAdd, toRm, err := c.Plan(zone)
	if err != nil {
		return 0, 0, err
	}

	for _, r := range toAdd {
		newRecord, err := c.AddRecord(zone, r.SubDomain, r.Target)
		if err != nil {
			return 0, 0, err
		}

		c.state.records = append(c.state.records, *newRecord)

		printAdd("%-16s %s.%s  added\n", r.Target, r.SubDomain, zone)
		added++
	}

	for _, r := range toRm {
		_, err := c.DeleteRecordByID(zone, r.ID)
		if err != nil {
			return 0, 0, err
		}

		for i, sr := range c.state.records {
			if sr.Zone == r.Zone && sr.SubDomain == r.SubDomain {
				c.state.records = append(c.state.records[:i], c.state.records[i+1:]...)
			}
		}

		printAdd("%-16s %s.%s  removed\n", r.Target, r.SubDomain, zone)
		removed++
	}

	if (len(toAdd) + len(toRm)) == 0 {
		// No modification
		return 0, 0, nil
	}

	err = c.RefreshZone(zone)
	if err != nil {
		return 0, 0, err
	}

	err = c.state.save()
	if err != nil {

		return 0, 0, err
	}

	return added, removed, nil
}
