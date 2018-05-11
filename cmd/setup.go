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
	"os/exec"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/fatih/color"
)


// struct holds user provided input data
type Inputs struct {
	nodeCount string
	networkName string
	publicIP string
	dockerSubnetIP string
	raftPort string
	rpcPort string
	constellationPort string
	gethPort string
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup XDC network",
	Long: `
Setup the XDC network`,
	Run: func(cmd *cobra.Command, args []string) {
		userInput:= Inputs{}
		portSelection(userInput)
		checkPortInUse("3000","localhost")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func portSelection(s Inputs) {
	c:= color.New(color.FgWhite,color.BgBlue) // create a new color object 
	
	fmt.Println()
	c.Println(" Define General Parameters ")
	fmt.Println()

	//define nodeCount prompt
	nodeCountPrompt := &survey.Input{Message: "Number Of Nodes To Create",Default:"4",Help:"No of nodes to create"}
	
	//define network name prompt
	networkNamePrompt := &survey.Input{Message: "Network Name",Default:"XDC-Network"}
	
	survey.AskOne(nodeCountPrompt, &s.nodeCount, nil)
	survey.AskOne(networkNamePrompt, &s.networkName, nil)

	fmt.Println()
	c.Println(" Define IP Address ")
	fmt.Println()

	//define publicIP prompt
	publicIPPrompt := &survey.Input{Message: "Public IP Address",Default:"192.168.0.174"}

	//define docker subnet prompt
	dockerSubnetIPPrompt := &survey.Input{Message: "Docker Subnet IP",Default:"172.17.0.1"}
	
	survey.AskOne(publicIPPrompt, &s.publicIP, nil)
	survey.AskOne(dockerSubnetIPPrompt, &s.dockerSubnetIP, nil)
	
	fmt.Println()
	c.Println(" Define Ports ")
	fmt.Println()

	//define raft port prompt
	raftPortPrompt := &survey.Input{Message: "Raft Port",Default:"23000"}
	
	//define rpc port prompt
	rpcPortPrompt := &survey.Input{Message: "RPC Port",Default:"22000"}

	//define constellation port prompt
	constellationPortPrompt := &survey.Input{Message: "Constellation Port",Default:"9000"}

	//define geth port prompt
	gethPortPrompt := &survey.Input{Message: "Geth Port",Default:"21000"}

	survey.AskOne(raftPortPrompt, &s.raftPort, nil)
	survey.AskOne(rpcPortPrompt, &s.rpcPort, nil)
	survey.AskOne(constellationPortPrompt, &s.constellationPort, nil)
	survey.AskOne(gethPortPrompt, &s.gethPort, nil)

	//fmt.Printf("%+v", s)
}

func checkPortInUse(port string, host string) {
	out,err:= exec.Command("nc","-zv",host,port).Output()
	if err==nil {
		fmt.Printf("Port not available %s", out)
	}  else {
		fmt.Println("Port available")
	}
}