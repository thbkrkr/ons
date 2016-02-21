package dns

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/thbkrkr/go-ovh-dns/dns"
	"github.com/yadutaf/go-ovh"
)

var (
	ConfigFilename = flag.String("f", "/config.json", "OVH API config file")
	Command        = flag.String("c", "plan", "Command: plan, apply, show, delete")
	Zone           = flag.String("z", "", "DNS zone to manage")
	ID             = flag.Int64("id", 0, "DNS record id (to delete)")
)

func loadConfig() (c *dns.Config, err error) {
	configFile, err := os.Open(*ConfigFilename)
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
		exit(fmt.Sprintf("Cannot parse config file `%s`:\n %v\n", *ConfigFilename, err))
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

	if *Zone == "" {
		exit(fmt.Sprint("-z is required: see `dns help`\n"))
	}

	if *Command == "show" {
		records, err := ovhClient.ListFullARecords(*Zone)
		b, err := json.Marshal(records)
		if err != nil {
			exit(fmt.Sprintf("Error on show: %v\n", err))
		}
		fmt.Printf("%v", string(b))
		return
	}

	stateFile := "dns.state"
	state, err := ovhClient.LoadState(stateFile)
	if err != nil {
		exit(fmt.Sprintf("Cannot load the state file: %s\n %v\n", stateFile, err))
	}

	if *Command == "plan" {
		records, err := ovhClient.ListFullARecords(*Zone)
		newState := ovhClient.Plan(state, records)

		b, err := json.Marshal(newState)
		if err != nil {
			exit(fmt.Sprintf("Error: %v\n", err))
		}
		fmt.Printf("%v", string(b))
		return
	}

	if *Command == "apply" {
		records, err := ovhClient.ListFullARecords(*Zone)
		modifications, err := ovhClient.Apply(*Zone, state, records)
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

		_, err := ovhClient.DeleteRecordByID(*Zone, *ID)
		if err != nil {
			exit(fmt.Sprintf("Error on delete: %v\n", err))
		}
		fmt.Printf("    ok: record %d deleted", *ID)
		return
	}

	exit(fmt.Sprintf("Invalid Command: %s\n", *Command))

	/*	record, err := ovhClient.GetRecord("blurb.space", 1323781110)
		if err != nil {
			exit(fmt.Sprintf("Error: %v\n", err))
		}
		b, err := json.Marshal(record)
		fmt.Printf("%v", string(b))
	*/

	//subDomain := "ovh"

	//_, err = ovhClient.addRecord(zone, "pof", "149.202.169.125")

	//records, err := ovhClient.listFullARecords(zone)
	//recordID, err := ovhClient.getRecordIDBySubDomain(zone, subDomain)
	//record, err := ovhClient.getRecord(zone, *recordID)
	//_, err = ovhClient.deleteRecordByID(*Zone, 1354103752)

	/*if err != nil {
		exit(fmt.Sprintf("Error: %v\n", err))
	}
	*/
}
