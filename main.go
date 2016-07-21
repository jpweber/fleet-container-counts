/*
* @Author: Jim Weber
* @Date:   2016-05-18 22:07:31
* @Last Modified by:   Jim Weber
* @Last Modified time: 2016-07-20 22:33:59
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// FleetStates struct to hold all the data for a machine state
type FleetStates struct {
	States []struct {
		SystemdActiveState string `json:"systemdActiveState"`
		MachineID          string `json:"machineID"`
		Hash               string `json:"hash"`
		SystemdSubState    string `json:"systemdSubState"`
		Name               string `json:"name"`
		SystemdLoadState   string `json:"systemdLoadState"`
	}
}

func getInstanceStates(deployment string, params map[string]string) FleetStates {
	url := "http://coreos." + deployment + ".crosschx.com:49153/fleet/v1/state"
	// loop through params to append to the url if they exist
	if len(params) > 0 {
		url = url + "?"
		for key, value := range params {
			// as of now we are only ever expecting a single k,v pair
			// for parameters
			url = url + key + "=" + value + ".service"
		}
	}

	response, err := http.Get(url)
	fleetStates := FleetStates{}

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(contents, &fleetStates); err != nil {
			panic(err)
		}

	}

	return fleetStates
}

func filterInstances(fleetStates FleetStates, appName string) []string {

	var instances []string
	for _, state := range fleetStates.States {
		if strings.Contains(state.Name, appName) {
			// exclude any presence or discovery units
			if strings.Contains(state.Name, "presence") || strings.Contains(state.Name, "discovery") {
				continue
			}
			instances = append(instances, state.Name)
		}
	}

	return instances
}

func getContainerCount(fleetUnits FleetStates) map[string]int {
	containerCount := make(map[string]int)
	for _, fleetUnit := range fleetUnits.States {
		shortNameParts := strings.Split(fleetUnit.Name, "@")
		shortName := shortNameParts[0]
		containerCount[shortName] += 1
	}

	return containerCount
}

func main() {
	fleetStates := getInstanceStates("dev", nil)
	containerCounts := getContainerCount(fleetStates)
	jsonCounts, err := json.Marshal(containerCounts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonCounts))
	// for k, v := range containerCounts {
	// 	fmt.Println(k, ":", v)
	// }
}
