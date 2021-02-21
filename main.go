package main

import (
	"flag"
	"fmt"
	"github.com/K8-Cloud/k8-cloud/SetupCluster"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var config HelmConfig

	var operation, configFile, context string
	//var version bool

	flag.StringVar(&operation, "o", "cluster", "Provide whether operation needed to be performed - Cluster/Addons")
	flag.StringVar(&configFile, "c", "cf-fmt.yaml", "Provide path to Config yaml")
	flag.StringVar(&context, "context", "minikube", "Provide kubernetes context for addon")
	version := flag.Bool("version", false, "display version")
	flag.Parse()

	if *version {
		fmt.Print("k8-cloud version: 0.6.0\n")
		os.Exit(0)
	}

	yamlFile, err := ioutil.ReadFile(configFile)

	makeDir("templates")

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
		getFileFromURL("templates/vpc-1.yaml","https://k8s-cloud-templates.s3.amazonaws.com/vpc-1.yaml")
		getFileFromURL("templates/0005-eks-cluster.yaml","https://k8s-cloud-templates.s3.amazonaws.com/0005-eks-cluster.yaml")
		getFileFromURL("templates/0007-esk-managed-node-group.yaml","https://k8s-cloud-templates.s3.amazonaws.com/0007-esk-managed-node-group.yaml")
		SetupCluster.CheckCluster(yamlFile)
	} else if operation == "enable_backup" {
		fmt.Print("Work In Progress\n")
	} else {
		fmt.Print("Operation Not Supported")
	}

	deleteDir("templates")
}