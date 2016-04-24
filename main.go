package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/thbkrkr/go-ovh-dns/dns"
)

var (
	ConfigFile = flag.String("f", "/ops/dns/config.json", "OVH API config file")
	StateFile  = flag.String("s", "/ops/dns/dns.state", "OVH API config file")
	Command    = flag.String("c", "plan", "Command: plan, apply, show, delete")

	ID = flag.Int64("id", 0, "DNS record id (to delete)")

	Zone = flag.String("z", "", "DNS zone to manage")
)

func loadConfig() (c *dns.Config, err error) {
	configFile, err := os.Open(*ConfigFile)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	var config *dns.Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func exit(msg string) {
	fmt.Print(msg)
	os.Exit(1)
}

func config() *dns.Config {
	config, err := loadConfig()
	if err != nil {
		exit(fmt.Sprintf("Cannot parse config file `%s`:\n %v\n", *ConfigFile, err))
	}

	return config
}

func createClient(config *dns.Config) (*dns.OvhClient, error) {
	ovhClient, err := ovh.NewClient(
		config.Endpoint,
		config.ApplicationKey,
		config.ApplicationSecret,
		config.ConsumerKey,
	)
	if err != nil {
		exit(fmt.Sprintf("Bad configuration:\n %v\n", err))
	}

	return &dns.OvhClient{Client: ovhClient}, nil
}

func main() {
	flag.Parse()
	config := config()

	if *Command == "help" {
		flag.PrintDefaults()
		exit("")
	}

	ovhClient, err := createClient(config)
	if err != nil {
		exit(fmt.Sprintf("Cannot initialize the client API:\n %v\n", err))
	}

	zone := *Zone
	if zone == "" {
		zone = config.Zone
	}

	if *Command == "show" {
		records, err := ovhClient.ListFullARecords(zone)
		b, err := json.Marshal(records)
		if err != nil {
			exit(fmt.Sprintf("Error on show: %v\n", err))
		}
		fmt.Printf("%v", string(b))
		return
	}

	stateFile, err := ovhClient.LoadState(*StateFile)
	if err != nil {
		exit(fmt.Sprintf("Cannot load the state file: %s\n %v\n", stateFile, err))
	}

	if *Command == "plan" {
		records, err := ovhClient.ListFullARecords(zone)
		newState := ovhClient.Plan(stateFile, records)

		b, err := json.Marshal(newState)
		if err != nil {
			exit(fmt.Sprintf("Error: %v\n", err))
		}
		fmt.Printf("%v", string(b))
		return
	}

	if *Command == "apply" {
		records, err := ovhClient.ListFullARecords(zone)
		modifications, err := ovhClient.Apply(zone, stateFile, records)
		if err != nil {
			exit(fmt.Sprintf("Error on apply: %s\n %v\n", stateFile, err))
		}

		if len(modifications) == 0 {
			exit(fmt.Sprint("    ok: 0 modification"))
		} else {
			totalModifs := len(modifications)
			if totalModifs == 1 {
				exit(fmt.Sprint("    ok: 1 modification"))
			} else {
				exit(fmt.Sprintf("    ok: %d modifications", totalModifs))

			}
		}
		return
	}

	if *Command == "delete" {
		if *ID == 0 {
			exit(fmt.Sprint("-id is required for command `delete`: see `dns help`\n"))
		}

		_, err := ovhClient.DeleteRecordByID(zone, *ID)
		if err != nil {
			exit(fmt.Sprintf("Error on delete: %v\n", err))
		}
		fmt.Printf("    ok: record %d deleted", *ID)
		return
	}

	exit(fmt.Sprintf("Invalid Command: %s\n", *Command))
}
