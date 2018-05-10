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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"log"
	"strconv"
)

var projectName string
var nodes int
var publicIP string
var image="xinfinorg/quorum:v2.0.0"

var RPC_PORT=22000
var GETH_PORT=21000
var RAFT_PORT=23000
var CONSTELLATION_PORT=9000
var qd string
var separator string
var enode_url string

//Test Enode Address, Key, Genesis
var enode_address="26e80451f629db9249cf1f325e1346863532987ec816103b3ef64d193b213786d80837dfebfd5d42ec05ed755c0e520739808fe9134efb350b7bbf9cb8fc5d06"
var keystr = "{\"address\":\"0638e1574728b6d862dd5d3a3e0942c3be47d996\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"d8119d67cb134bc65c53506577cfd633bbbf5acca976cea12dd507de3eb7fd6f\",\"cipherparams\":{\"iv\":\"76e88f3f246d4bf9544448d1a27b06f4\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"6d05ade3ee96191ed73ea019f30c02cceb6fc0502c99f706b7b627158bfc2b0a\"},\"mac\":\"b39c2c56b35958c712225970b49238fb230d7981ef47d7c33c730c363b658d06\"},\"id\":\"00307b43-53a3-4e03-9d0c-4fcbb3da29df\",\"version\":3}"
var genesisstr="{\"alloc\": {    \"0638e1574728b6d862dd5d3a3e0942c3be47d996\": {      \"balance\": \"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\"    }  },  \"coinbase\": \"0x0000000000000000000000000000000000000000\",  \"config\": {    \"homesteadBlock\": 0,    \"chainId\": 1,    \"eip155Block\": null,    \"eip158Block\": null,    \"isQuorum\": true  },  \"difficulty\": \"0x0\",  \"extraData\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",  \"gasLimit\": \"0x47b760\",  \"mixhash\": \"0x00000000000000000000000000000000000000647572616c65787365646c6578\",  \"nonce\": \"0x0\",  \"parentHash\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",  \"timestamp\": \"0x00\"}"

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up your Network")

		//Get Woring Directory
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		//Create Key File
		keyfilepath := filepath.Join(dir, "/key")
		keyfile, err := os.Create(keyfilepath)
		if err != nil {
			log.Fatal("Cannot create key file", err)
		}
		defer keyfile.Close()

		fmt.Fprintf(keyfile, keystr)

		//Create Genesis File
		genesisfilepath := filepath.Join(dir, "/genesis.json")
		genesisfile, err := os.Create(genesisfilepath)
		if err != nil {
			log.Fatal("Cannot create genesis file", err)
		}
		defer genesisfile.Close()

		fmt.Fprintf(genesisfile, genesisstr)

		//Create static-nodes.json file
		staticnodesfilepath := filepath.Join(dir, "/static-nodes.json")
		staticnodesfile, err := os.Create(staticnodesfilepath)
		if err != nil {
			log.Fatal("Cannot create file static-nodes.json", err)
		}
		defer staticnodesfile.Close()

		fmt.Fprintf(staticnodesfile, "[\n")

		//Create dir structure/files add enode addresses to static-nodes.json
		for i := 1; i <= nodes; i++ {
			qd= "qdata_" + strconv.Itoa(i)

			logpath := filepath.Join(dir, "/"+qd+"/logs")
			os.MkdirAll(logpath, os.ModePerm);

			keyspath := filepath.Join(dir, "/"+qd+"/keys")
			os.MkdirAll(keyspath, os.ModePerm);

			gethpath := filepath.Join(dir, "/"+qd+"/dd/geth")
			os.MkdirAll(gethpath, os.ModePerm);

			keystorepath := filepath.Join(dir, "/"+qd+"/dd/keystore")
			os.MkdirAll(keystorepath, os.ModePerm);

			passwordfilepath := filepath.Join(dir, "/"+qd+"/passwords.txt")
			os.Create(passwordfilepath)

			if i < nodes {
				separator = ","
				} else {
				separator = ""
			}

			enode_url="enode://"+ enode_address +"@"+publicIP+":"+strconv.Itoa(GETH_PORT+i)+"?discport=0&raftport="+strconv.Itoa(RAFT_PORT+i)

			fmt.Fprintf(staticnodesfile, strconv.Quote(enode_url)+separator+"\n")

		 }

		fmt.Fprintf(staticnodesfile, "]")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVarP(&projectName,"projectname", "p", "XDC-NW", "Project name for your network")
	setupCmd.Flags().IntVarP(&nodes,"nodes",  "n", 2,"Enter number of inital static nodes")
	setupCmd.Flags().StringVarP(&publicIP,"ipaddress",  "i", "","Enter Public IP of this machine")

	setupCmd.MarkFlagRequired("projectname")
	setupCmd.MarkFlagRequired("nodes")
	setupCmd.MarkFlagRequired("ipaddress")
}
