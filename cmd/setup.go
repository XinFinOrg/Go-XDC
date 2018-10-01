// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"

	"bytes"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/spf13/viper"
)

// struct holds user provided input data
type Inputs struct {
	NetworkName    string `mapstructure:"network_name"`
	PublicIP       string `mapstructure:"public_ip"`
	DockerSubnetIP string `mapstructure:"docker_ip"`
	PortRange      int    `mapstructure:"port_range"`
	Nodes          int
}

// port selection struct
var PortSel struct {
	geth          []int
	rpc           []int
	raft          []int
	constellation []int
}

var GETH_PORT_OFFSET = 0
var RAFT_PORT_OFFSET = 100
var CONSTELLATION_PORT_OFFSET = 200
var RPC_PORT_OFFSET = 300

var image = "xinfinorg/quorum:v2.1.0"

var qd string
var separator string
var enode_id string
var enode_url string
var hostPath string
var containerPath string
var cmdString1 []string
var cmdString2 []string
var cmdString3 []string

//Test keystore, genesis file
var keystore_string = "{\"address\":\"0638e1574728b6d862dd5d3a3e0942c3be47d996\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"d8119d67cb134bc65c53506577cfd633bbbf5acca976cea12dd507de3eb7fd6f\",\"cipherparams\":{\"iv\":\"76e88f3f246d4bf9544448d1a27b06f4\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"6d05ade3ee96191ed73ea019f30c02cceb6fc0502c99f706b7b627158bfc2b0a\"},\"mac\":\"b39c2c56b35958c712225970b49238fb230d7981ef47d7c33c730c363b658d06\"},\"id\":\"00307b43-53a3-4e03-9d0c-4fcbb3da29df\",\"version\":3}"
var genesis_string = "{\"alloc\": {    \"0638e1574728b6d862dd5d3a3e0942c3be47d996\": {      \"balance\": \"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\"    }  },  \"coinbase\": \"0x0000000000000000000000000000000000000000\",  \"config\": {    \"homesteadBlock\": 0,    \"chainId\": 1,    \"eip155Block\": null,    \"eip158Block\": null,    \"isQuorum\": true  },  \"difficulty\": \"0x0\",  \"extraData\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",  \"gasLimit\": \"0x47b760\",  \"mixhash\": \"0x00000000000000000000000000000000000000647572616c65787365646c6578\",  \"nonce\": \"0x0\",  \"parentHash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",  \"timestamp\": \"0x00\"}"

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup XDC network",
	Long: `
Setup the XDC network`,
	Run: func(cmd *cobra.Command, args []string) {
		userInput := Inputs{}

		filename := "config.yml"

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			getUserInput(&userInput)
		} else {
			readConfigFile(&userInput)
		}

		setupNetwork(&userInput)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	/*setupCmd.Flags().StringVarP(&projectName,"projectname", "p", "XDC-NW", "Project name for your network")
	setupCmd.Flags().IntVarP(&nodes,"nodes",  "n", 2,"Enter number of inital static nodes")
	setupCmd.Flags().StringVarP(&PublicIP,"ipaddress",  "i", "","Enter Public IP of this machine")

	setupCmd.MarkFlagRequired("projectname")
	setupCmd.MarkFlagRequired("nodes")
	setupCmd.MarkFlagRequired("ipaddress")*/
}

func readConfigFile(s *Inputs) {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&s)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	fmt.Println("Network Name: ", s.NetworkName)
	fmt.Println("Public IP: ", s.PublicIP)
	fmt.Println("Docker Subnet: ", s.DockerSubnetIP)
	fmt.Println("Starting Port Range:", s.PortRange)
	fmt.Println("Number of Nodes: ", s.Nodes)
}

func getUserInput(s *Inputs) {
	c := color.New(color.FgWhite, color.BgHiBlue) // create a new color object

	fmt.Println()
	c.Println(" Define General Parameters ")
	fmt.Println()

	//define nodeCount prompt
	nodeCountPrompt := &survey.Input{Message: "Number Of Nodes To Create", Default: "4", Help: "No of Nodes to create"}

	//define network name prompt
	networkNamePrompt := &survey.Input{Message: "Network Name", Default: "XDC-Network"}

	survey.AskOne(nodeCountPrompt, &s.Nodes, nil)
	survey.AskOne(networkNamePrompt, &s.NetworkName, nil)

	fmt.Println()
	c.Println(" Define IP Address ")
	fmt.Println()

	//define PublicIP prompt
	publicIPPrompt := &survey.Input{Message: "Public IP Address", Default: "0.0.0.0"}

	//define docker subnet prompt
	dockerSubnetIPPrompt := &survey.Input{Message: "Docker Subnet IP", Default: "172.13.0.0/16"}

	survey.AskOne(publicIPPrompt, &s.PublicIP, nil)
	survey.AskOne(dockerSubnetIPPrompt, &s.DockerSubnetIP, nil)

	fmt.Println()
	c.Println(" Define Ports ")
	fmt.Println()

	//define port range prompt
	portPrompt := &survey.Input{Message: "Assign Ports from", Default: "20000"}

	survey.AskOne(portPrompt, &s.PortRange, nil)

	//fmt.Printf("%+v", s)
}

/* Get unused port when PortRange and portType is defined
host - host name (localhost)
portType - [GETH,RAFT,CONSTELLATION,RPC]
PortRange - scan free ports after this port
*/
func getUnusedPort(host string, portType int, portRange int) int {
	portt := 0

	for {
		switch portType { // check portType (rpc,geth,..)
		case 0: // GETH PORT
			portt = portRange + GETH_PORT_OFFSET
			GETH_PORT_OFFSET++
		case 1: // RAFT PORT
			portt = portRange + RAFT_PORT_OFFSET
			RAFT_PORT_OFFSET++
		case 2: // CONSTELLATION_PORT
			portt = portRange + CONSTELLATION_PORT_OFFSET
			CONSTELLATION_PORT_OFFSET++
		case 3: // RPC_PORT
			portt = portRange + RPC_PORT_OFFSET
			RPC_PORT_OFFSET++
		}

		_, err := exec.Command("nc", "-zv", host, strconv.Itoa(portt)).Output()
		if err == nil {
			//fmt.Println("Port unavailable... checking next", strconv.Itoa(portt), out, err)
			portt++
		} else {
			//fmt.Println("Port available ", strconv.Itoa(portt))
			if portType == 0 {
				PortSel.geth = append(PortSel.geth, portt)
			} else if portType == 1 {
				PortSel.raft = append(PortSel.raft, portt)
			} else if portType == 2 {
				PortSel.constellation = append(PortSel.constellation, portt)
			} else if portType == 3 {
				PortSel.rpc = append(PortSel.rpc, portt)
			}
			break
		}
	}
	return portt
}

func setupNetwork(s *Inputs) {
	scaffoldNodeDir(s)
	createTmFiles(s)
	createStartNodeScript(s)
	createDockerComposeFile(s)
	setupComplete()
}

func setupComplete() {
	c := color.New(color.FgWhite, color.BgHiGreen) // create a new color object

	fmt.Println()
	c.Println(" Setup complete ")
	fmt.Println()
}

func scaffoldNodeDir(s *Inputs) {

	//Get Working Directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n - Getting current working directory - " + dir)

	//Create static-nodes.json file
	staticnodesfilepath := filepath.Join(dir, "/static-nodes.json")
	staticnodesfile, err := os.Create(staticnodesfilepath)
	if err != nil {
		log.Fatal("Cannot create file static-nodes.json", err)
	}
	defer staticnodesfile.Close()

	fmt.Fprintf(staticnodesfile, "[\n")

	//Create dir structure/files & add enode urls to static-nodes.json [for each node]
	for i := 1; i <= s.Nodes; i++ {
		qd = "qdata_" + strconv.Itoa(i)

		fmt.Println(" - Creating dir structure for node " + strconv.Itoa(i) + " - " + qd)

		//create logs folder
		logpath := filepath.Join(dir, "/"+qd+"/logs")
		os.MkdirAll(logpath, os.ModePerm)

		//create keys folder
		keyspath := filepath.Join(dir, "/"+qd+"/keys")
		os.MkdirAll(keyspath, os.ModePerm)

		//create dd folder
		gethpath := filepath.Join(dir, "/"+qd+"/dd/geth")
		os.MkdirAll(gethpath, os.ModePerm)

		keystorepath := filepath.Join(dir, "/"+qd+"/dd/keystore")
		os.MkdirAll(keystorepath, os.ModePerm)

		//create password file
		passwordfilepath := filepath.Join(dir, "/"+qd+"/passwords.txt")
		os.Create(passwordfilepath)

		//create start-node.sh file - TO-DO Populate this file with custom parameters & then start geth/constellation
		startnodefilepath := filepath.Join(dir, "/"+qd+"/start-node.sh")
		os.Create(startnodefilepath)

		// Set permissions for start-node.sh
		err := os.Chmod(startnodefilepath, 0755)
		if err != nil {
			log.Println(err)
		}

		//Separator for static-nodes.json urls - no comma after last node address
		if i < s.Nodes {
			separator = ","
		} else {
			separator = ""
		}

		//Create Key File
		keyfilepath := filepath.Join(dir, "/"+qd+"/dd/keystore/key")
		keyfile, err := os.Create(keyfilepath)
		if err != nil {
			log.Fatal("Cannot create key file", err)
		}
		defer keyfile.Close()

		fmt.Fprintf(keyfile, keystore_string)

		//Create Genesis File
		genesisfilepath := filepath.Join(dir, "/"+qd+"/genesis.json")
		genesisfile, err := os.Create(genesisfilepath)
		if err != nil {
			log.Fatal("Cannot create genesis file", err)
		}
		defer genesisfile.Close()

		fmt.Fprintf(genesisfile, genesis_string)

		//Set docker host & container mount directories
		hostPath = filepath.Join(dir, "/"+qd)
		containerPath = "/qdata"

		//Generate nodekey
		cmdString1 = []string{"/usr/local/bin/bootnode", "-genkey", "/qdata/dd/nodekey"}
		runDockerContainer(hostPath, containerPath, cmdString1)

		//Return hash for generated key for use in static-nodes
		cmdString2 = []string{"/usr/local/bin/bootnode", "--nodekey", "/qdata/dd/nodekey", "-writeaddress"}
		enode_id = runDockerContainer(hostPath, containerPath, cmdString2)
		enode_id = strings.TrimRight(enode_id, "\r\n")
		fmt.Println(" - Generating nodeKey for node " + strconv.Itoa(i))
		//Construct enode url
		enode_url = "enode://" + enode_id + "@" + s.PublicIP + ":" + strconv.Itoa(getUnusedPort("localhost", 0, s.PortRange)) + "?discport=0&raftport=" + strconv.Itoa(getUnusedPort("localhost", 1, s.PortRange))
		fmt.Fprintf(staticnodesfile, strconv.Quote(enode_url)+separator+"\n")

		//assign/reserve rpc port
		getUnusedPort("localhost", 3, s.PortRange)

		// Generate Quorum-related keys (used by Constellation)
		// NOTE: using sh here as this command asks user input for the password
		// < /dev/null would set an empty password and > /dev/null  will set no output
		cmdString3 = []string{"sh", "-c", "/usr/local/bin/constellation-node --generatekeys=/qdata/keys/tm < /dev/null > /dev/null"}
		runDockerContainer(hostPath, containerPath, cmdString3)

		//To-Do
		// 1-Create Constellation config files & keys
		// 2-Update Ports in start-node.sh
		// 3-Create docker-compose.yaml
		// 4-Cleanup

	}

	fmt.Fprintf(staticnodesfile, "]")

	//copy static-nodes.json to each folder
	copyStaticNode(s)

	//YAML MockConfig file for docker-compose

	/*c, err1 := loadConfig("docker-compose.yml")
	if err1 != nil {
		panic(err1)
	}
	fmt.Printf("%+v\n", c)*/

}

//create docker-compose.yaml file
func createDockerComposeFile(s *Inputs) {
	fmt.Println(" - Generating docker-compose.yml file")
	err1 := saveConfig(createMockConfig(s), "docker-compose.yml")
	if err1 != nil {
		panic(err1)
	}
}

//copy static-nodes.json to each qdata folder
func copyStaticNode(s *Inputs) {
	for i := 1; i <= s.Nodes; i++ {
		_, err := exec.Command("cp", "static-nodes.json", "qdata_"+strconv.Itoa(i)+"/dd").Output()
		if err != nil {
			panic(err)
		}
	}

}

//create tm conf files [ TO-DO ]
func createTmFiles(s *Inputs) {

	// allocate constellation ports for all the Nodes
	for i := 1; i <= s.Nodes; i++ {
		getUnusedPort("localhost", 2, s.PortRange)
	}

	//Get Working Directory
	dir, _ := os.Getwd()

	for i := 1; i <= s.Nodes; i++ {
		fmt.Println(" - Generating tm.conf files for node " + strconv.Itoa(i))
		qd = "qdata_" + strconv.Itoa(i)
		//create tmconf
		tmfilepath := filepath.Join(dir, "/"+qd+"/tm.conf")
		tmfile, _ := os.Create(tmfilepath)
		defer tmfile.Close()
		fmt.Fprintln(tmfile, `url ="http://`+s.PublicIP+":"+strconv.Itoa(PortSel.constellation[i-1])+`"`)
		fmt.Fprintln(tmfile, `port =`+strconv.Itoa(PortSel.constellation[i-1]))
		fmt.Fprintln(tmfile, `socket = "/qdata/tm.ipc"`)
		fmt.Fprintln(tmfile, `othernodes = `+makeOtherNodesString(s.PublicIP))
		fmt.Fprintln(tmfile, `publickeys = ["/qdata/keys/tm.pub"]`)
		fmt.Fprintln(tmfile, `privatekeys = ["/qdata/keys/tm.key"]`)
		fmt.Fprintln(tmfile, `storage = "/qdata/constellation"`)
		fmt.Fprintln(tmfile, `verbosity = 3`)
	}
}

func makeOtherNodesString(publicIP string) string {
	var str string

	for _, elem := range PortSel.constellation {
		str = str + `"http://` + publicIP + `:` + strconv.Itoa(elem) + `",`
	}

	lastChar := str[len(str)-1:]

	if lastChar == `,` {
		str = strings.TrimRight(str, ",")
	}

	return `[` + str + `]`
}

func runDockerContainer(hostPath string, containerPath string, cmdString []string) string {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmdString,
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostPath,
				Target: containerPath,
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	//io.Copy(os.Stdout, out)
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	stdoutStr := buf.String()
	return stdoutStr
}

type Node struct {
	Image    string   `yaml:"image"`
	Restart  string   `yaml:"restart"`
	Volumes  []string `yaml:"volumes"`
	Networks []string `yaml:"networks"`
	Ports    []string `yaml:"ports"`
	User     string   `yaml:"user"`
}

type Network struct {
	Driver string `yaml:"driver"`
	IPAM   IPAM   `yaml:"ipam"`
}
type IPAM struct {
	Driver string   `yaml:"driver"`
	Config []Subnet `yaml:"config"`
}

type Subnet struct {
	Subnet string `yaml:"subnet"`
}

type Configuration struct {
	Version  string             `yaml:"version"`
	Services map[string]Node    `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
}

func createMockConfig(s *Inputs) Configuration {
	servtest1 := map[string]Node{}
	for i := 0; i < s.Nodes; i++ {
		rpcPortStr := strconv.Itoa(PortSel.rpc[i])
		gethPortStr := strconv.Itoa(PortSel.geth[i])
		constellationPortStr := strconv.Itoa(PortSel.constellation[i])
		raftPortStr := strconv.Itoa(PortSel.raft[i])
		servtest1["node_"+strconv.Itoa(i+1)] = Node{
			Image:    "xinfinorg/quorum:v2.1.0",
			Restart:  "always",
			Volumes:  []string{"./qdata_" + strconv.Itoa(i+1) + ":/qdata"},
			Networks: []string{"xdc_network"},
			Ports:    []string{gethPortStr + ":" + gethPortStr, rpcPortStr + ":" + rpcPortStr, constellationPortStr + ":" + constellationPortStr, raftPortStr + ":" + raftPortStr},
			User:     "0:0",
		}
	}
	return Configuration{
		Version:  "2",
		Services: servtest1,
		Networks: map[string]Network{
			"xdc_network": Network{
				Driver: "bridge",
				IPAM: IPAM{
					Driver: "default",
					Config: []Subnet{Subnet{Subnet: s.DockerSubnetIP}},
				},
			},
		},
	}
}

func createStartNodeScript(s *Inputs) {
	var script = `#!/bin/bash

	#
	# This is used at Container start up to run the constellation and geth Nodes
	#
	
	set -u
	set -e
	
	### Configuration Options
	TMCONF=/qdata/tm.conf
	
	GETH_ARGS="--datadir /qdata/dd --raft --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum,raft --nodiscover --unlock 0 --password /qdata/passwords.txt --rpcport 0000 --port 0000 --raftport 0000"
	
	if [ ! -d /qdata/dd/geth/chaindata ]; then
	  echo "[*] Mining Genesis block"
	  /usr/local/bin/geth --datadir /qdata/dd init /qdata/genesis.json
	fi
	
	echo "[*] Starting Constellation node"
	nohup /usr/local/bin/constellation-node $TMCONF 2>> /qdata/logs/constellation.log &
	
	sleep 2
	
	echo "[*] Starting node"
	PRIVATE_CONFIG=$TMCONF nohup /usr/local/bin/geth $GETH_ARGS 2>&1 >>/qdata/logs/geth.log | tee --append /qdata/logs/geth.log`

	dir, _ := os.Getwd()

	for i := 1; i <= s.Nodes; i++ {
		fmt.Println(" - Generating startup-node.sh file for node " + strconv.Itoa(i))
		tempScr := script
		tempScr = strings.Replace(tempScr, "--rpcport 0000", "--rpcport "+strconv.Itoa(PortSel.rpc[i-1]), 1)
		tempScr = strings.Replace(tempScr, "--port 0000", "--port "+strconv.Itoa(PortSel.geth[i-1]), 1)
		tempScr = strings.Replace(tempScr, "--raftport 0000", "--raftport "+strconv.Itoa(PortSel.raft[i-1]), 1)

		snfilepath := filepath.Join(dir, "/qdata_"+strconv.Itoa(i)+"/start-node.sh")
		snfile, _ := os.Create(snfilepath)
		defer snfile.Close()
		fmt.Fprint(snfile, tempScr)
	}
}

func saveConfig(c Configuration, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func loadConfig(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}
