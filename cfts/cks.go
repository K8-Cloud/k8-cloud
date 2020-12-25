package cfts

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/smallfish/simpleyaml"
	_ "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var filename = "cf-aws-fmt.yml"

type Config struct {
	Key   string
	Value string
}

type Cftvpc struct {
	StackName   string
	TemplateURL string
}

//AKS Vars
//var (
//	ctx        = context.Background()
//	clientData clientInfo
//	authorizer autorest.Authorizer
//)

type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

//Setup AKS or EKS Cluster

func setupCluster() {

	////Reading inputs from yaml

	//filename := "cf-aws-fmt.yml" EKS
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}

	////Creating Elements
	//Start EKS Cluster elements
	//ElementsSubnetIDs := make(map[string]string)
	var Acceesskey, Secretkey, Region string
	//var nodelen int
	//var nodelist []interface{}
	//var sess session.Session
	//End EKS Cluster elements

	//passing values for setting up connections

	//Start EKS Cluster session values
	Cloud, err := yaml.Get("Cloud").Get("Name").String()
	if Cloud == "AWS" {
		Acceesskey, _ = yaml.Get("Cloud").Get("AccessKey").String()
		Secretkey, _ = yaml.Get("Cloud").Get("SecretAccKey").String()
		Region, _ = yaml.Get("Cloud").Get("Region").String()
	}
	//End EKS Cluster elements session values

	//Print EKS Cluster elements
	if Cloud == "AWS" {
		fmt.Printf("Cloud: %#v\n", Cloud)
		fmt.Printf("AccessKey: %#v\n", Acceesskey)
		fmt.Printf("SecAccKey: %#v\n", Secretkey)
		fmt.Printf("Region: %#v\n", Region)
		fmt.Printf("Creating sessions")
	}

	////Create Sessions
	//Create session EKS Cluster elements
	if Cloud == "AWS" {
		//sess, err := session.NewSession(&aws.Config{
		//	//aws.Config{throttle.Throttle()}
		//	Region:      aws.String(Region),
		//	Credentials: credentials.NewStaticCredentials(Acceesskey, Secretkey, ""),
		//})
	}
	//End session EKS Cluster elements
	fmt.Printf("Session created ")

	////Checking if VPC is enabled

}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return &contents, nil
}
