package client

// DNSConfig represents a DNS zone records configuration
type DNSConfig struct {
	configPath string
	records    []Record
}

func loadConfig(configPath string) (*DNSConfig, error) {
	records, err := loadRecords(configPath)
	if err != nil {
		return nil, err
	}

	return &DNSConfig{
		configPath: configPath,
		records:    records,
	}, nil
}

func (c *DNSConfig) save() error {
	return saveRecords(c.configPath, c.records)
}
