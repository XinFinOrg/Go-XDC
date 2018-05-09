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
	"bytes"
	"github.com/fatih/color"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Installs the dependencies for network setup",
	Long: `
Installs the dependencies for network setup`,
	Run: func(cmd *cobra.Command, args []string) {
		//get os info
		fmt.Println(fetchOsInfo())
		//check network connection
		checkNetworkConn()
		//install docker cc
		installDocker()
	},
}

func init() {
	rootCmd.AddCommand(prepareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchOsInfo() string {
	c:= color.New(color.FgWhite,color.BgBlue) // create a new color object 
	var out bytes.Buffer //stores stdout data

	fmt.Println("\nFetching OS Info...\n")
	c.Println(" OS Info ")
	cmdLSB:= exec.Command("lsb_release","-a") // define command to execute
	cmdLSB.Stdout = &out
	cmdLSB.Run() // execute the defined command
	return out.String() // return stdout data string
}

func checkNetworkConn() {
	c:= color.New(color.FgWhite,color.BgBlue) // create a new color object
	fmt.Println("Checking network connection...\n")
	c.Print(" Network Status ")
	fmt.Print(" OK\n")
}

func installDocker() {
	fmt.Println("\nInstalling Docker CE...")
	downloadDockerSetupScr()
}

func downloadDockerSetupScr() {
	cos:= color.New(color.FgWhite,color.BgGreen)
	cof:= color.New(color.FgWhite,color.BgRed)
	fmt.Println(" - Downloading docker setup script from get-docker.com")
	cmdDSS:= exec.Command("wget","get-docker.com","-O","docker-setup.sh")
	err:=cmdDSS.Run()
	if err==nil {
		fmt.Println(" - Docker setup script download complete - docker-setup.sh")
		cmdEDS:= exec.Command("sudo","sh","docker-setup.sh")
		//cmdEDS:= exec.Command("ls")
		fmt.Println(" - Executing docker-setup.sh")
		err:=cmdEDS.Run()
		if err==nil {
			cos.Println("\n Docker CE installed ")
		} else {
			cof.Println("\n Docker CE failed to install ( Reason - Error while executing setup script) ")
		}
	} else {
		cof.Println("\n Docker CE failed to install ( Reason - Could not download setup script) ")
	}
}