package aks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	yaml2 "github.com/go-yaml/yaml"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"log"
	"os"
	"time"
	//"github.com/Azure/go-autorest/autorest/azure/Authorizer"
	//_ "encoding/json"
	//"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/util"
	//"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	//"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	//_ "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	//_ "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	//"github.com/Azure/go-autorest/autorest/to"
	//_ "gopkg.in/yaml.v2"
)

var filename = "cf-aws-fmt.yml"

//AKS Vars
var (
	ctx        = context.Background()
	clientData clientInfo
	authorizer autorest.Authorizer
	authoriz   auth.ClientCredentialsConfig
)

type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
	filename := "cf-fmt.yaml-azure" // AKS
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}

	var ServicePrinciple, ResourceGroup, Region string

	Cloud, err := yaml.Get("Cloud").Get("Name").String()
	if Cloud == "Azure" {
		ServicePrinciple, _ = yaml.Get("Cloud").Get("ServicePrinciple").String()
		ResourceGroup, _ = yaml.Get("Cloud").Get("ResourceGroup").String()
		Region, _ = yaml.Get("Cloud").Get("Region").String()
	}

	if Cloud == "Azure" {
		fmt.Printf("Cloud: %#v\n", Cloud)
		fmt.Printf("ServicePrinciple: %#v\n", ServicePrinciple)
		fmt.Printf("ResourceGroup: %#v\n", ResourceGroup)
		fmt.Printf("Region: %#v\n", Region)
		fmt.Printf("Creating sessions")
	}
	if Cloud == "Azure" {
		os.Setenv("AZURE_AUTH_LOCATION", string(ServicePrinciple))
		//var err error
		//authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		////if err != nil {
		////log.Fatalf("Failed to get OAuth config: %v", err)
		////}
		////authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
		////if err != nil {
		////	log.Fatalf("Failed to read JSON: %+v", err)
		////}

		authoriz = auth.NewClientCredentialsConfig("481ff77d-be69-4e96-9f4c-f55b17d048d2", "lAJTW0zdCMP475WzfRArfG0G_ombo30_Hz", "d897562d-b399-4e76-b1d0-449fbf6c2c1f")
		//autorest.Authorizer = authrize.Authorizer()
		//if err != nil {
		//log.Fatalf("Failed to get OAuth config: %v", err)
		//}
		////authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
		////if err != nil {
		////	log.Fatalf("Failed to read JSON: %+v", err)
		////}

		////clientData.SubscriptionID = (*authInfo)["subscriptionId"].(string)
		//clientData.VMPassword = (*authInfo)["clientSecret"].(string)
		////fmt.Printf("subscriptionId: %#v\n", clientData.SubscriptionID)
	}

	fmt.Printf("Session created ")

	VPCClient := resources.NewDeploymentsClient("d3d5c541-15a8-46f3-9aae-c8948f3067e6")
	VPCClient.Authorizer, err = authoriz.Authorizer()
	if err != nil {
		log.Fatalf("Failed to read JSON: %+v", err)
	}

	template, err := ioutil.ReadFile("aks-template.json")

	fmt.Println(string(template))
	if err != nil {
		return
	}

	params, err := ioutil.ReadFile("aks-params.json")
	if err != nil {
		return
	}
	fmt.Println(string(params))

	yamlparams, err := ioutil.ReadFile("aks-params.yaml")
	var body interface{}
	if err := yaml2.Unmarshal([]byte(yamlparams), &body); err != nil {
		panic(err)
	}

	body = convert(body)

	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Output: %s\n", b)
	}

	contentstemp := make(map[string]interface{})
	contentsparms := make(map[string]interface{})

	json.Unmarshal(template, &contentstemp)
	//json.Unmarshal(params, &contentsparms)
	json.Unmarshal(b, &contentsparms)
	Properties := &resources.DeploymentProperties{
		Template:   &contentstemp,
		Parameters: &contentsparms,
		Mode:       resources.Complete,
	}

	Deployment := resources.Deployment{
		Properties: Properties,
	}

	Validation, err := VPCClient.Validate(ctx, ResourceGroup, "testdeploy", Deployment)
	if err != nil {
		fmt.Printf("ERR: %#v\n", err)
	}
	fmt.Printf("Validation: %#v\n", Validation.Status)

	//util.PrintAndLog("validated VM template deployment")

	deploy, err := VPCClient.CreateOrUpdate(ctx, ResourceGroup, "testdeploy", Deployment)

	//deploy, err := VPCClient.CreateOrUpdate(ctx, ResourceGroup, "testdeploy", resources.Deployment{
	//	Properties: &resources.DeploymentProperties{
	//		Template:   contentstemp,
	//		Parameters: contentsparms,
	//		Mode:       resources.Complete,
	//	}})

	fmt.Printf("Deployment: %#v\n", deploy.Response())
	fmt.Println(err)
	err = deploy.Future.WaitForCompletionRef(ctx, VPCClient.BaseClient.Client)
	if err != nil {
		return
	}
	test, _ := deploy.Result(VPCClient)
	fmt.Println(test.Status)

	list, err := VPCClient.Get(ctx, ResourceGroup, "testdeploy")

	//fmt.Println(list.Body)
	fmt.Printf("List: %#v\n", list.Properties.Parameters)

	//delete, _ := VPCClient.Delete(ctx, ResourceGroup,"testdeploy")
	//err = delete.Future.WaitForCompletionRef(ctx, VPCClient.BaseClient.Client)
	//if err != nil {
	//	return
	//}
	// fmt.Printf("Del: %#v\n", delete)
	//fmt.Print(deploy)
	//test2, err := VPCClient.CreateOrUpdate(context.Background(),
	//	"AKS-test-Resourse-Group",
	//	"Test-VNet-2",
	//	network.VirtualNetwork{
	//		Location: to.StringPtr("westus"),
	//		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
	//			AddressSpace: &network.AddressSpace{
	//				AddressPrefixes: &[]string{"10.0.0.0/8"},
	//			},
	//			Subnets: &[]network.Subnet{
	//				{
	//					Name: to.StringPtr("public"),
	//					SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
	//						AddressPrefix: to.StringPtr("10.0.0.0/16"),
	//					},
	//				},
	//				{
	//					Name: to.StringPtr("private"),
	//					SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
	//						AddressPrefix: to.StringPtr("10.1.0.0/16"),
	//					},
	//				},
	//			},
	//		},
	//	})
	//if err != nil {
	//	fmt.Printf("Networks: %#v\n", test2.Response())
	//}
	//fmt.Printf("Networks: %#v\n",test2.)
	//test, err := VPCClient.List(context.Background(), "AKS-test-Resourse-Group")
	//if err != nil {
	//	fmt.Printf("Networks Err List: %#v\n", err)
	//}
	//result := strings.Join(test.Values(), ",")
	//fmt.Printf("%s\n",string(json.MarshalIndent(&test.Values()[0], "", " ")))
	//jstr, err := json.MarshalIndent(&test.Values()[0], "", "  ")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%s\n", string(jstr))

	//test.Values()

	//CreateOrUpdate(context.Background(),"AKS-test-Resourse-Group")
	//listvpc, _ := network.NewVirtualNetworksClient(clientData.SubscriptionID).List( context.Background(), "Test-resource-group")
	//List(ctx.B,
	//"Test-resource-group")
	//fmt.Printf("Networks List: %#v\n", listvpc)
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

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
