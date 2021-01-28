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
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
		helmInit(context)
		helmAddRepositories(config)
		fmt.Print(config)
		helmInstallReleases(config, context)
	}else if operation == "cluster" {
		getFileFromURL("vpc-1.yaml","https://k8s-cloud-templates.s3.amazonaws.com/vpc-1.yaml")
		getFileFromURL("0005-eks-cluster.yaml","https://k8s-cloud-templates.s3.amazonaws.com/0005-eks-cluster.yaml")
		getFileFromURL("0007-esk-managed-node-group.yaml","https://k8s-cloud-templates.s3.amazonaws.com/0007-esk-managed-node-group.yaml")
		SetupCluster.CheckCluster(yamlFile)
	} else {
		fmt.Print("Operation Not Supported")
	}
}