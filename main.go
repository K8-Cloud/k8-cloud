package main

import (
	"flag"
	"fmt"
	"github.com/K8-Cloud/k8-cloud/SetupCluster"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	var config HelmConfig

	var operation, configFile, context string

	flag.StringVar(&operation, "o", "Cluster", "Provide whether operation needed to be performed - Cluster/Addons")
	flag.StringVar(&configFile, "c", "cf-fmt.yaml", "Provide path to Config yaml")
	flag.StringVar(&context, "context", "minikube", "Provide kubernetes context for addon")
	flag.Parse()
	fmt.Printf("Operation: %v\n", operation)
	fmt.Printf("Config File: %v\n", configFile)
	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	if operation == "addons" {
	err = yaml.Unmarshal(yamlFile, &configFile)
		if err != nil {
			panic(err)
		}
	}else if operation == "setup_addon" {
		helmInit(context)
		helmAddRepositories(config)
		helmInstallReleases(config, context)
	//}else if operation == "setup_addon3" {
	//	helm3setup()
	//	helm3AddRepositories(config)
	} else if operation == "cluster" {
		SetupCluster.CheckCluster(yamlFile)
	} else {
		fmt.Print("Operation Not Supported")
	}
}