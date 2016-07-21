/*
* @Author: Jim Weber
* @Date:   2016-05-18 22:07:31
* @Last Modified by:   Jim Weber
* @Last Modified time: 2016-07-20 23:09:57
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
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

func getContainerCount(fleetUnits FleetStates) map[string]int {
	containerCount := make(map[string]int)
	for _, fleetUnit := range fleetUnits.States {
		shortNameParts := strings.Split(fleetUnit.Name, "@")
		shortName := shortNameParts[0]
		containerCount[shortName] += 1
	}

	return containerCount
}

func getDeployEnv(fleetHost string) string {
	// fleetURL := "http://172.17.0.1"
	fleetURL := "http://" + fleetHost
	cfg := client.Config{
		Endpoints: []string{fleetURL + ":4001"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// get "/foo" key's value
	kapi := client.NewKeysAPI(etcdClient)
	resp, err := kapi.Get(context.Background(), "/deployment", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print value
		log.Println("deployment envrionment is", resp.Node.Value)
	}

	return resp.Node.Value
}

func main() {
	deploymentPtr := flag.String("e", "172.17.0.1", "ETCD Host Address")
	prettyPrintPtr := flag.Bool("p", false, "Human readble pretty print rather than json output for application")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	environ := getDeployEnv(*deploymentPtr)
	fleetStates := getInstanceStates(environ, nil)
	containerCounts := getContainerCount(fleetStates)
	jsonCounts, err := json.Marshal(containerCounts)
	if err != nil {
		fmt.Println(err)
	}
	if *prettyPrintPtr == true {
		for k, v := range containerCounts {
			fmt.Println(k, ":", v)
		}
	} else {
		fmt.Println(string(jsonCounts))
	}
}
