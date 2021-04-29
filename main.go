package main

import (
	"flag"
	"fmt"
	"github.com/K8-Cloud/k8-cloud/SetupCluster"
	"github.com/K8-Cloud/k8-cloud/manageCluster"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
)

func setupK8sConnection(InitialConfigVals InitialConfigVals) *kubernetes.Clientset {
	fmt.Println("Setting up Connection")
	fmt.Printf("MasterUrl: %v\n", InitialConfigVals.ClusterDetails.MasterUrl)
	fmt.Printf("KubeConfig: %v\n", InitialConfigVals.ClusterDetails.KubeConfig)
	connection, _ := manageCluster.SetupConnection(InitialConfigVals.ClusterDetails.MasterUrl, InitialConfigVals.ClusterDetails.KubeConfig)

	fmt.Printf("ClusterName: %v\n", InitialConfigVals.ClusterDetails.ClusterName)
	fmt.Printf("MasterKey: %v\n", InitialConfigVals.ClusterDetails.MasterKey)
	fmt.Printf("Configs: %v\n", InitialConfigVals.ClusterDetails.Configs)
	fmt.Printf("StorageClasses.yaml: %v\n", InitialConfigVals.ClusterDetails.StorageClassFile)
	fmt.Printf("Namepaces.yaml: %v\n", InitialConfigVals.ClusterDetails.NameSpaceFile)
	fmt.Printf("ResourceQuotas.yaml: %v\n", InitialConfigVals.ClusterDetails.ResourceQuotaFile)

	return connection
}

func main() {
	var config HelmConfig
	var clustername, masterurl, kubeconfig string
	var InitialConfigVals InitialConfigVals
	var operation, configFile, context, name string
	//var version bool

	//flag.StringVar(&operation, "o", "all", "Provide the operation that needs to be performed, valid inputs - namespace, storage, resourcequota, defaultquota, serviceaccount")
	flag.StringVar(&operation, "o", "cluster", "Provide whether operation needed to be performed - Cluster/Addons")
	flag.StringVar(&configFile, "c", "cf-fmt.yaml", "Provide path to Config yaml")
	flag.StringVar(&context, "context", "minikube", "Provide kubernetes context for addon")
	flag.StringVar(&name, "name", "backup", "backup name")
	version := flag.Bool("version", false, "display version")
	flag.StringVar(&clustername, "k", "dev-cluster", "Provide cluster name")
	flag.StringVar(&masterurl, "u", "https://localhost:6443", "Provide master url")
	flag.StringVar(&kubeconfig, "x", "~/.kube/config", "Provide path to kubeconfig")
	flag.Parse()

	if *version {
		fmt.Print("k8-cloud version: 1.0.0\n")
		os.Exit(0)
	}

	if operation == "addons" {
		yamlFile, err := ioutil.ReadFile(configFile)

		makeDir("templates")

		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
		helmInit(context)
		helmAddRepositories(config)
		fmt.Print(config)
		helmInstallReleases(config, context)
	} else if operation == "cluster" {
		yamlFile, err := ioutil.ReadFile(configFile)

		makeDir("templates")

		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

		SetupCluster.CheckCluster(yamlFile)
	} else if operation == "take_backup" {
		takeBackup(name, context)
		fmt.Print("Work In Progress\n")
	} else {
		fmt.Print("Operation Not Supported")
	}

	var manageOperation = StrSlice{"all", "init", "namespace", "storage", "resourcequota", "defaultquota", "serviceaccount"}

	if manageOperation.Has(operation) {
		filePath := "K8Cli" + "/mgmt/" + clustername

		InitClusterConfig, err := ioutil.ReadFile(filePath + "/config.yaml")
		if err != nil {
			fmt.Println(err)
			//panic(err)
		}
		err = yaml.Unmarshal([]byte(InitClusterConfig), &InitialConfigVals)
		if err != nil {
			panic(err)
		}
	}


	if operation == "all" {

		connection := setupK8sConnection(InitialConfigVals)

		fmt.Println("Executing Create or Update StorageClasses")
		manageCluster.CreateorUpdateStorageClass(InitialConfigVals.ClusterDetails.StorageClassFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

		fmt.Println("Executing Create or Update NameSpaces")
		manageCluster.CreateorUpdateNameSpace(InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

		fmt.Println("Executing Create or Update DefaultQuotas")
		manageCluster.CreateorUpdateDefaultQuota(InitialConfigVals.ClusterDetails.Configs, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

		fmt.Println("Executing Create or Update ResourceQuotas")
		manageCluster.CreateorUpdateResourceQuota(InitialConfigVals.ClusterDetails.ResourceQuotaFile, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

		fmt.Println("Executing Create or Update NameSpaceUsers")
		manageCluster.CreateorUpdateNameSpaceUser(InitialConfigVals.ClusterDetails.Configs, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "namespace" {

		connection := setupK8sConnection(InitialConfigVals)
		fmt.Println("Executing Create or Update NameSpaces")
		manageCluster.CreateorUpdateNameSpace(InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "storage" {

		connection := setupK8sConnection(InitialConfigVals)
		fmt.Println("Executing Create or Update StorageClasses")
		manageCluster.CreateorUpdateStorageClass(InitialConfigVals.ClusterDetails.StorageClassFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "resourcequota" {

		connection := setupK8sConnection(InitialConfigVals)
		fmt.Println("Executing Create or Update NameSpaces")
		manageCluster.CreateorUpdateNameSpace(InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)
		fmt.Println("Executing Create or Update DefaultQuotas")
		manageCluster.CreateorUpdateDefaultQuota(InitialConfigVals.ClusterDetails.Configs, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)
		fmt.Println("Executing Create or Update ResourceQuotas")
		manageCluster.CreateorUpdateResourceQuota(InitialConfigVals.ClusterDetails.ResourceQuotaFile, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "defaultquota" {

		connection := setupK8sConnection(InitialConfigVals)
		fmt.Println("Executing Create or Update NameSpaces")
		manageCluster.CreateorUpdateNameSpace(InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)
		fmt.Println("Executing Create or Update DefaultQuotas")
		manageCluster.CreateorUpdateDefaultQuota(InitialConfigVals.ClusterDetails.Configs, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "serviceaccount" {

		connection := setupK8sConnection(InitialConfigVals)
		fmt.Println("Executing Create or Update NameSpaces")
		manageCluster.CreateorUpdateNameSpace(InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)
		fmt.Println("Executing Create or Update NameSpaceUsers")
		manageCluster.CreateorUpdateNameSpaceUser(InitialConfigVals.ClusterDetails.Configs, InitialConfigVals.ClusterDetails.NameSpaceFile, connection, InitialConfigVals.ClusterDetails.MasterKey)

	} else if operation == "init" {

		fmt.Println("Initializing K8Cli")
		fmt.Printf("ClusterName: %v\n", &clustername)
		fmt.Printf("masterurl: %v\n", masterurl)
		fmt.Printf("kubeconfig: %v\n", kubeconfig)
		manageCluster.Init(clustername, masterurl, kubeconfig)

	} else {

		fmt.Printf("MasterUrl: %v\n", InitialConfigVals.ClusterDetails.MasterUrl)
		fmt.Printf("KubeConfig: %v\n", InitialConfigVals.ClusterDetails.KubeConfig)
		fmt.Printf("MasterKey: %v\n", InitialConfigVals.ClusterDetails.MasterKey)
		fmt.Printf("Configs: %v\n", InitialConfigVals.ClusterDetails.Configs)
		fmt.Printf("StorageClasses.yaml: %v\n", InitialConfigVals.ClusterDetails.StorageClassFile)
		fmt.Printf("Namepaces.yaml: %v\n", InitialConfigVals.ClusterDetails.NameSpaceFile)
		fmt.Printf("ResourceQuotas.yaml: %v\n", InitialConfigVals.ClusterDetails.ResourceQuotaFile)
		fmt.Println("Provide Valid input operation")
	}

	deleteDir("templates")
}
