package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/portshift/agent/pkg/plugin_common"
	"io/ioutil"
	"os"
)

type ansibleResponse struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

func exitJson(responseBody ansibleResponse) {
	returnResponse(responseBody)
}

func failJson(responseBody ansibleResponse) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody ansibleResponse) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(ansibleResponse{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func main() {
	var response ansibleResponse

	helpStr := "" +
		"USAGE:\n" +
		"	Use portshift_ansible_plugin as a task in your ansible playbook, and provide the parameters listed below:\n" +
		"OPTIONS:\n" +
		"	url              (string)   Url of the Portshift management (default: \"console.portshift.io\")\n" +
		"	access-key       (string)   Access key of the service. (Required)\n" +
		"	secret-key       (string)   Secret key of the service. (Required)\n" +
		"	executable-path  (string)   Path to the executable.\n" +
		"	process-name     (string)   Name of the process.\n" +
		"	args             (string)   Command line arguments of the executable. (ex. args: \" arg1 arg2 arg3 \")\n" +
		"	cwd              (string)   Working directory of the executable.\n" +
		"	label            ([]string) Application label in the format: key=value. Should be specified as list of strings (ex. label: [ \"key1=label1\", \"key2=label2\" ] )\n" +
		"	name             (string)   A unique name of the App (will appear in Portshift UI). (Required)\n" +
		"	type             (string)   Type of the app. (Required)\n" +
		"	value            (string)   Name of the executable. (Required)\n\n"

	if len(os.Args) != 2 {
		fmt.Printf("%s", helpStr)
		response.Msg = "No argument file provided"
		failJson(response)
	}

	if os.Args[1] == "--help" {
		fmt.Printf("%s", helpStr)
		os.Exit(0)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Msg = fmt.Sprintf("Could not read configuration file %s: %s", argsFile, err)
		failJson(response)
	}

	var moduleArgs plugin_common.CreateAppParams

	err = json.Unmarshal(text, &moduleArgs)
	if err != nil {
		response.Msg = fmt.Sprintf("Configuration file %s not valid JSON: %s", argsFile, err)
		failJson(response)
	}

	plugin_common.CreateApp(moduleArgs)
	exitJson(response)
}
