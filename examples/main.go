package main

import "C"
import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)
const storageClassSuffix string = ".storageclass.storage.k8s.io/"
type NameSpaceVals struct {
	NameSpace []struct {
		Name            string      `yaml:"Name"`
		ResourceQuota 	string      `yaml:"ResourceQuota"`
		DefaultQuota    string		`yaml:"DefaultQuota"`
		Labels         map[string]string `yaml:"Labels"`
	} `yaml:"NameSpace"`
}
type ResourceQuotaVals struct {
	ResourceQuota []struct{
		ResourceQuotaName string	`yaml:"QuotaName"`
		RequestsCPU 	  string	`yaml:"RequestsCPU"`
		LimitsCPU		  string	`yaml:"LimitsCPU"`
		RequestsMemory	  string	`yaml:"RequestsMemory"`
		LimitsMemory      string	`yaml:"LimitsMemory"`
		Pods			  string	`yaml:"Pods"`
		RequestsStorage   string	`yaml:"Name"`
		RequestsEphemeralStorage	string	`yaml:"RequestsStorage"`
		LimitsEphemeralStorage		string	`yaml:"LimitsEphemeralStorage"`
		StorageClasses []struct{
			Name string `yaml:"Name"`
			RequestsStorage string `yaml:"RequestsStorage"`
		} `yaml:"StorageClasses"`
		Labels         map[string]string `yaml:"Labels"`
	}	`yaml:"ResourceQuota"`
}
type StorageClassVals struct {
	StorageClasses []struct {
		Name              string                      `yaml:"Name"`
		Provisioner       string                      `yaml:"Provisioner"`
		Parameters        map[string]string `yaml:"Parameters"`
		ReclaimPolicy     string                      `yaml:"ReclaimPolicy"`
		VolumeBindingMode string                      `yaml:"VolumeBindingMode"`
		Labels            map[string]string `yaml:"Labels"`
	} `yaml:"StorageClasses"`
}
type NameSpaceRoleVals struct {
	NameSpaceRoleDetails struct {
		AppendName string `yaml:"AppendName"`
		Labels     map[string]string `yaml:"Labels"`
		PolicyRules []struct{
			APIGroups []string `yaml:"APIGroups"`
			Resources []string `yaml:"Resources"`
			Verbs 	  []string `yaml:"Verbs"`
		} `yaml:"PolicyRules"`
	} `yaml:"NameSpaceRoleDetails"`
}
type DefaultQuotaVals struct {
	DefaultQuota struct {
		Details []struct {
			Name                 v1.LimitType                 `yaml:"Name"`
			Max                  map[v1.ResourceName]resource.Quantity  `yaml:"max"`
			Min                  map[v1.ResourceName]resource.Quantity  `yaml:"min"`
			Default              map[v1.ResourceName]resource.Quantity  `yaml:"default,omitempty"`
			DefaultRequest       map[v1.ResourceName]resource.Quantity  `yaml:"defaultRequest,omitempty"`
			MaxLimitRequestRatio map[v1.ResourceName]resource.Quantity  `yaml:"Details"`
		} `yaml:"Details"`
		Labels  map[string]string `yaml:"Labels"`
	} `yaml:"DefaultQuota"`
}

func main()  {
//https://192.168.56.2:6443

	var resourcequotayaml, namespaceyaml, masterurl, kubeconfig, config, operation string
	flag.StringVar(&masterurl, "u", "https://localhost:6443", "Provide master url")
	flag.StringVar(&kubeconfig, "c", "./.kube/config", "Provide path to kubeconfig")
	flag.StringVar(&config, "f", "./config/", "Provide path to config")
	flag.StringVar(&operation, "o", "all", "Provide the operation that needs to be performed, valid inputs - namespace, storage, resourcequota, defaultquota, serviceaccount")
	flag.StringVar(&namespaceyaml, "n", "./config/NameSpacesList/Namespaces.yaml", "Provide the path to Namespaces.yaml")
	flag.StringVar(&resourcequotayaml, "r", "./config/ResourceQuotaList/ResourceQuota.yaml", "Provide the path to Namespaces.yaml")


	flag.Parse()
	fmt.Printf("masterurl: %v\n", masterurl)
	fmt.Printf("kubeconfig: %v\n", kubeconfig)
	fmt.Printf("Operation: %v\n", operation)
	fmt.Printf("config: %v\n", config)
	fmt.Printf("namespaceyaml: %v\n", namespaceyaml)
	fmt.Printf("resourcequotayaml: %v\n", resourcequotayaml)


	fmt.Println("Setting up Connection")
	connection, _ := SetupConnection(masterurl,kubeconfig)

	if operation == "all"{
		fmt.Println("Executing Create or Update StorageClasses")
		CreateorUpdateStorageClass(config, connection)

		fmt.Println("Executing Create or Update NameSpaces")
		CreateorUpdateNameSpace(namespaceyaml, connection)

		fmt.Println("Executing Create or Update DefaultQuotas")
		CreateorUpdateDefaultQuota(config, namespaceyaml, connection)

		fmt.Println("Executing Create or Update ResourceQuotas")
		CreateorUpdateResourceQuota (resourcequotayaml, namespaceyaml, connection)


		fmt.Println("Executing Create or Update NameSpaceUsers")
		CreateorUpdateNameSpaceUser(config, namespaceyaml, connection)

	} else if operation == "namespace" {

		fmt.Println("Executing Create or Update NameSpaces")
		CreateorUpdateNameSpace(namespaceyaml, connection)

	} else if operation == "storage" {

		fmt.Println("Executing Create or Update StorageClasses")
		CreateorUpdateStorageClass(config, connection)

	} else if operation == "resourcequota" {

		fmt.Println("Executing Create or Update NameSpaces")
		CreateorUpdateNameSpace(namespaceyaml, connection)

		fmt.Println("Executing Create or Update DefaultQuotas")
		CreateorUpdateDefaultQuota(config, namespaceyaml, connection)

		fmt.Println("Executing Create or Update ResourceQuotas")
		CreateorUpdateResourceQuota (resourcequotayaml, namespaceyaml, connection)

	} else if operation == "defaultquota" {

		fmt.Println("Executing Create or Update NameSpaces")
		CreateorUpdateNameSpace(namespaceyaml, connection)

		fmt.Println("Executing Create or Update DefaultQuotas")
		CreateorUpdateDefaultQuota(config, namespaceyaml, connection)

	} else if operation == "serviceaccount" {

		fmt.Println("Executing Create or Update NameSpaces")
		CreateorUpdateNameSpace(namespaceyaml, connection)

		fmt.Println("Executing Create or Update NameSpaceUsers")
		CreateorUpdateNameSpaceUser(config, namespaceyaml, connection)

	} else {
		fmt.Println("Provide Valid input operation")
	}
}

func SetupConnection(url string, kubeconfig string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags(url, kubeconfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return c, err
}
func CreateorUpdateStorageClass(config string, connection *kubernetes.Clientset) error {

	con := connection
	var StorageClassVals StorageClassVals

	fileNameSpace, err := ioutil.ReadFile(config+"StorageClasses.yaml")
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}

	err = yaml.Unmarshal([]byte(fileNameSpace), &StorageClassVals)
	if err != nil {
		panic(err)
	}

	var reclaimPolicy v1.PersistentVolumeReclaimPolicy
	var vbmode 	storagev1.VolumeBindingMode
	var name, Provisioner string
	LenSC := len(StorageClassVals.StorageClasses)

	mapLabels := make(map[string]string)
	mapParams := make(map[string]string)


	for i := 0; i < LenSC; i++ {

		for key, value := range StorageClassVals.StorageClasses[i].Labels {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			mapLabels[strKey] = strValue
		}

		for key, value := range StorageClassVals.StorageClasses[i].Parameters {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			mapParams[strKey] = strValue
		}

		if StorageClassVals.StorageClasses[i].ReclaimPolicy == "" {
			reclaimPolicy = v1.PersistentVolumeReclaimRetain
		} else if StorageClassVals.StorageClasses[i].ReclaimPolicy == "Retain"{
			reclaimPolicy = v1.PersistentVolumeReclaimRetain
		} else if StorageClassVals.StorageClasses[i].ReclaimPolicy == "Recycle" {
			reclaimPolicy = v1.PersistentVolumeReclaimRecycle
		} else if StorageClassVals.StorageClasses[i].ReclaimPolicy == "Delete" {
			reclaimPolicy = v1.PersistentVolumeReclaimDelete
		} else {
			fmt.Println("Reclaim Policy is not correct")
			return err
		}

		//vbmode := storagev1.VolumeBindingImmediate
		if StorageClassVals.StorageClasses[i].VolumeBindingMode == "" {
			vbmode = storagev1.VolumeBindingWaitForFirstConsumer
		} else if StorageClassVals.StorageClasses[i].VolumeBindingMode == "Immediate"{
			vbmode = storagev1.VolumeBindingImmediate
		} else if StorageClassVals.StorageClasses[i].VolumeBindingMode == "WaitForConsumer" {
			vbmode = storagev1.VolumeBindingWaitForFirstConsumer
		} else {
			fmt.Println("Volume Binding Mode is not correct")
			return err
		}


		if StorageClassVals.StorageClasses[i].Name == "" {
			name = "standard-local"
		} else {
			name = StorageClassVals.StorageClasses[i].Name
		}

		if StorageClassVals.StorageClasses[i].Provisioner == "" {
			Provisioner = "kubernetes.io/no-provisioner"
		} else {
			Provisioner = StorageClassVals.StorageClasses[i].Provisioner
		}

		storageclassjson := storagev1.StorageClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StorageClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: mapLabels,
			},
			Provisioner: Provisioner,
			Parameters:  mapParams,
			ReclaimPolicy: &reclaimPolicy,
			VolumeBindingMode: &vbmode,
		}

		fmt.Println("Storage Class ID: ", i)
		fmt.Println("Storage Class Name: ", name)
		fmt.Println("Storage Class Labels: ", mapLabels)
		fmt.Println("Storage Class Provisioner: ", Provisioner)
		fmt.Println("Storage Class ReclaimPolicy: ", reclaimPolicy)
		fmt.Println("Storage Class VolumeBindingMode: ", vbmode)
		CreateSC, err := con.StorageV1().StorageClasses().Create(context.TODO(), &storageclassjson, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			fmt.Println("Updating StorageClass.....")
			UpdateSC, err := con.StorageV1().StorageClasses().Update(context.TODO(), &storageclassjson, metav1.UpdateOptions{})
			if err != nil{
				fmt.Println(err)
			} else {
				fmt.Println("Updated StorageClass : ", UpdateSC.Name)
			}

			//return nil
		} else {
			fmt.Println("Created StorageClass : ", CreateSC.Name)
		}

	}

	return err
}
func CreateorUpdateNameSpace(namespaceyaml string, connection *kubernetes.Clientset) error {

	var NameSpaceVals NameSpaceVals
	con := connection
	mapLabels := make(map[string]string)

	fileNameSpace, err := ioutil.ReadFile(namespaceyaml)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &NameSpaceVals)
	if err != nil {
		panic(err)
	}

	lenNs := len(NameSpaceVals.NameSpace)

	//println(LenRQ)

 //Create or Update NameSpace

	for i := 0; i < lenNs; i++ {

		for key, value := range NameSpaceVals.NameSpace[i].Labels {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			mapLabels[strKey] = strValue
		}

		// 	NS Details
		fmt.Println("NameSpace ID: ", i)
		fmt.Println("NameSpace Name:	", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("NameSpace Resource Quota:	", NameSpaceVals.NameSpace[i].ResourceQuota)
		fmt.Println("NameSpace Default Quota:	", NameSpaceVals.NameSpace[i].DefaultQuota)
		//fmt.Println("NameSpace Labels: 		 ", NameSpaceVals.NameSpace[i].Labels)

		namespacejson := v1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:   NameSpaceVals.NameSpace[i].Name,
				Labels: mapLabels,
			},
		}

		fmt.Println("NameSpace Labels: ", mapLabels)

		CreateNameSpace, err := con.CoreV1().Namespaces().Create(context.TODO(), &namespacejson, metav1.CreateOptions{})

		if err != nil {
			fmt.Println(err)
			fmt.Println("Updating NameSpace........")
			UpdateNameSpace, _ := con.CoreV1().Namespaces().Update(context.TODO(), &namespacejson, metav1.UpdateOptions{})
			if err != nil{
				fmt.Println(err)
			} else {
				fmt.Println("Updated NameSpace: ", UpdateNameSpace.Name)
			}


		} else {
			println("Created NameSpace : ", CreateNameSpace.Name)
			//Catch the resource details for attaching Resources
		}
	}
	return nil
}
func CreateorUpdateResourceQuota(resourcequotayaml string, namespaceyaml string, connection *kubernetes.Clientset) error {

	var NameSpaceVals NameSpaceVals
	var ResourceQuotaVals ResourceQuotaVals
	con := connection
	var CatchCount int
	var UpdateResourceQuota, CreateResourceQuota *v1.ResourceQuota
	mapLabels := make(map[string]string)
	var StorageClassName string
	var ResourceRequestsCPU, ResourceLimitsCPU, ResourceRequestsMemory, ResourceLimitsMemory, ResourcePods, ResourceRequestsStorage, ResourceRequestsEphemeralStorage, ResourceLimitsEphemeralStorage, StorageClassVolume resource.Quantity

	fileNameSpace, err := ioutil.ReadFile(namespaceyaml)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &NameSpaceVals)
	if err != nil {
		panic(err)
	}

	fileResourceQuota, err := ioutil.ReadFile(resourcequotayaml)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileResourceQuota), &ResourceQuotaVals)
	if err != nil {
		panic(err)
	}

	lenNs := len(NameSpaceVals.NameSpace)
	LenRQ := len(ResourceQuotaVals.ResourceQuota)

	//Create or Update ResourceQuota

	for i := 0; i < lenNs; i++ {

		for key, value := range NameSpaceVals.NameSpace[i].Labels {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			mapLabels[strKey] = strValue
		}

		// NS Details
		fmt.Println("NameSpace Selected - ID: ", i)
		fmt.Println("NameSpace Selected - Name: ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("NameSpace Selected - Labels: ", mapLabels)


		if LenRQ == 0 {
			fmt.Println("Resource Quotas list is not provided")
		} else {
			// Find the matching Resource Quota
			for k := 0; k < LenRQ; k++ {
				if ResourceQuotaVals.ResourceQuota[k].ResourceQuotaName == NameSpaceVals.NameSpace[i].ResourceQuota {
					CatchCount = k
				}
			}
			for key, value := range ResourceQuotaVals.ResourceQuota[CatchCount].Labels {
				strKey := fmt.Sprintf("%v", key)
				strValue := fmt.Sprintf("%v", value)
				mapLabels[strKey] = strValue
			}

			// Count Storage Classes Defined in Resource Quota

			CountStclass := len(ResourceQuotaVals.ResourceQuota[CatchCount].StorageClasses)

			TotalLen := CountStclass + 8

			arroptkey := make([]v1.ResourceName, int(TotalLen))
			arroptValue := make([]resource.Quantity, int(TotalLen))
			arrayresult3 := make(map[v1.ResourceName]resource.Quantity)
			//lists :=  make(map[v1.ResourceName]resource.Quantity)
			arrkey := [8]v1.ResourceName{v1.ResourceRequestsCPU, v1.ResourceLimitsCPU, v1.ResourceRequestsMemory, v1.ResourceLimitsMemory, v1.ResourcePods, v1.ResourceRequestsStorage, v1.ResourceRequestsEphemeralStorage, v1.ResourceLimitsEphemeralStorage}

			for i := 0; i < 8; i++ {
				arroptkey[i] = arrkey[i]
			}

			if ResourceQuotaVals.ResourceQuota[CatchCount].RequestsCPU == "" {
				ResourceRequestsCPU = resource.MustParse("1")
			} else {
				ResourceRequestsCPU = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].RequestsCPU)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].LimitsCPU == "" {
				ResourceLimitsCPU = resource.MustParse("2")
			} else {
				ResourceLimitsCPU = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].LimitsCPU)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].RequestsMemory == "" {
				ResourceRequestsMemory = resource.MustParse("10Mi")
			} else {
				ResourceRequestsMemory = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].RequestsMemory)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].LimitsMemory == "" {
				ResourceLimitsMemory = resource.MustParse("10Mi")
			} else {
				ResourceLimitsMemory = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].LimitsMemory)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].Pods == "" {
				ResourcePods = resource.MustParse("100")
			} else {
				ResourcePods = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].Pods)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].RequestsStorage == "" {
				ResourceRequestsStorage = resource.MustParse("10M")
			} else {
				ResourceRequestsStorage = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].RequestsStorage)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].RequestsEphemeralStorage == "" {
				ResourceRequestsEphemeralStorage = resource.MustParse("10M")
			} else {
				ResourceRequestsEphemeralStorage = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].RequestsEphemeralStorage)
			}
			if ResourceQuotaVals.ResourceQuota[CatchCount].LimitsEphemeralStorage == "" {
				ResourceLimitsEphemeralStorage = resource.MustParse("10M")
			} else {
				ResourceLimitsEphemeralStorage = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].LimitsEphemeralStorage)
			}

			arrVal := [8]resource.Quantity{ResourceRequestsCPU, ResourceLimitsCPU, ResourceRequestsMemory, ResourceLimitsMemory, ResourcePods, ResourceRequestsStorage, ResourceRequestsEphemeralStorage, ResourceLimitsEphemeralStorage}

			for i := 0; i < 8; i++ {
				arroptValue[i] = arrVal[i]
			}

			for i := 0; i < 8; i++ {

				strKey := arroptkey[i]
				strValue := arroptValue[i]
				arrayresult3[strKey] = strValue
			}

			if CountStclass == 0 {

			} else {
				for j := 0; j < CountStclass; j++ {
					if len(ResourceQuotaVals.ResourceQuota[CatchCount].StorageClasses[j].Name) == 0 {
						StorageClassName = "standard-local"
					} else {
						StorageClassName = ResourceQuotaVals.ResourceQuota[CatchCount].StorageClasses[j].Name
					}
					if len(ResourceQuotaVals.ResourceQuota[CatchCount].StorageClasses[j].RequestsStorage) == 0 {
						StorageClassVolume = resource.MustParse("10M")
					} else {
						StorageClassVolume = resource.MustParse(ResourceQuotaVals.ResourceQuota[CatchCount].StorageClasses[j].RequestsStorage)
					}

					fmt.Println("Adding Strorage Class: ", StorageClassName, StorageClassVolume.String())

					arroptkey[8+j] = V1ResourceByStorageClass(StorageClassName, v1.ResourceRequestsStorage)

					arroptValue[8+j] = StorageClassVolume

					strKey := arroptkey[8+j]
					strValue := arroptValue[8+j]
					arrayresult3[strKey] = strValue
				}

			}
			fmt.Println(arrayresult3)

			resourcequotajson := v1.ResourceQuota{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ResourceQuota",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      NameSpaceVals.NameSpace[i].Name + "-resquota",
					Namespace: NameSpaceVals.NameSpace[i].Name,
					Labels:    mapLabels,
				},
				Spec: v1.ResourceQuotaSpec{
					Hard: arrayresult3,
				},
			}

			fmt.Println("Resource ID: ", i)
			fmt.Println("Resource Name: ", NameSpaceVals.NameSpace[i].Name + "-resquota")
			fmt.Println("Resource NameSpace: ", NameSpaceVals.NameSpace[i].Name)
			fmt.Println("Resource Labels: ", mapLabels)
			fmt.Println("Resource Hard Limits: ", arrayresult3)


			CreateResourceQuota, err = con.CoreV1().ResourceQuotas(NameSpaceVals.NameSpace[i].Name).Create(context.TODO(), &resourcequotajson, metav1.CreateOptions{})

			if err != nil {
				fmt.Println(err)
				fmt.Println("Updating ResourceQuota........")
				UpdateResourceQuota, err = con.CoreV1().ResourceQuotas(NameSpaceVals.NameSpace[i].Name).Update(context.TODO(), &resourcequotajson, metav1.UpdateOptions{})
				if err != nil{
					fmt.Println(err)
				} else {
					fmt.Println("Updated ResourceQuota: ", UpdateResourceQuota.Name)
				}
			} else {
				fmt.Println("Created ResourceQuota: ", CreateResourceQuota.Name)
			}
		}
	}
	return nil
}
func CreateorUpdateDefaultQuota(config string, namespaceyaml string, connection *kubernetes.Clientset) {

	con := connection
	var DefaultQuotaVals DefaultQuotaVals
	var NameSpaceVals NameSpaceVals

	fileNameSpace, err := ioutil.ReadFile(namespaceyaml)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &NameSpaceVals)
	if err != nil {
		panic(err)
	}

	fileNameSpace, err = ioutil.ReadFile(config+"DefaultQuota.yaml")
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &DefaultQuotaVals)
	if err != nil {
		panic(err)
	}

	LenNS := len(NameSpaceVals.NameSpace)
	LenLR := len(DefaultQuotaVals.DefaultQuota.Details)
	//fmt.Println(LenLR)
	mapLabels := make(map[string]string)

	LimitRangeItem := make([]v1.LimitRangeItem, LenLR)

	for j := 0; j < LenLR; j++ {
		LimitRangeItem[j] = v1.LimitRangeItem{
			Type:                 DefaultQuotaVals.DefaultQuota.Details[j].Name,
			Max:                  DefaultQuotaVals.DefaultQuota.Details[j].Max,
			Min:                  DefaultQuotaVals.DefaultQuota.Details[j].Min,
			Default:              DefaultQuotaVals.DefaultQuota.Details[j].Default,
			DefaultRequest:       DefaultQuotaVals.DefaultQuota.Details[j].DefaultRequest,
			MaxLimitRequestRatio: DefaultQuotaVals.DefaultQuota.Details[j].MaxLimitRequestRatio,
		}
	}

	for key, value := range DefaultQuotaVals.DefaultQuota.Labels {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapLabels[strKey] = strValue
	}

	for i := 0; i < LenNS; i++ {

		fmt.Println("NameSpace Selected - ID: ", i)
		fmt.Println("NameSpace Selected - Name: ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("NameSpace Selected - Labels: ", mapLabels)

		defaultquotajson := v1.LimitRange{
			TypeMeta: metav1.TypeMeta{
				Kind:       "LimitRange",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      NameSpaceVals.NameSpace[i].Name + "-defaultquota",
				Labels:    mapLabels,
				Namespace: NameSpaceVals.NameSpace[i].Name,
			},
			Spec: v1.LimitRangeSpec{
				Limits: LimitRangeItem,
			},
		}


		fmt.Println("DefaultQuota Name ", NameSpaceVals.NameSpace[i].Name + "-defaultquota")
		fmt.Println("DefaultQuota NameSpace: ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("DefaultQuota Labels: ", mapLabels)
		fmt.Println("DefaultQuota LimitRanges: ", LimitRangeItem)

		CreateDefaultQuota, err := con.CoreV1().LimitRanges(NameSpaceVals.NameSpace[i].Name).Create(context.TODO(), &defaultquotajson, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			fmt.Println("Updating ServiceAccount.....")
			UpdateSerAcc, _ := con.CoreV1().LimitRanges(NameSpaceVals.NameSpace[i].Name).Update(context.TODO(), &defaultquotajson, metav1.UpdateOptions{})
			if err != nil{
				fmt.Println(err)
			} else {
				fmt.Println("Updated ServiceAccount: ", UpdateSerAcc.Name)
			}
		} else {
			fmt.Println("Created ServiceAccount: ", CreateDefaultQuota.Name)
		}
	}
}
func CreateorUpdateNameSpaceUser(config string, namespaceyaml string, connection *kubernetes.Clientset) error {

	var NameSpaceVals NameSpaceVals
	var NameSpaceRoleVals NameSpaceRoleVals
	con := connection
	mapLabels := make(map[string]string)

	fileNameSpace, err := ioutil.ReadFile(namespaceyaml)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &NameSpaceVals)
	if err != nil {
		panic(err)
	}

	fileNameSpace, err = ioutil.ReadFile(config+"DefaultNameSpaceRole.yaml")
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	err = yaml.Unmarshal([]byte(fileNameSpace), &NameSpaceRoleVals)
	if err != nil {
		panic(err)
	}

	lenNs := len(NameSpaceVals.NameSpace)
	for key, value := range NameSpaceRoleVals.NameSpaceRoleDetails.Labels {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapLabels[strKey] = strValue
	}

	lenPL := len(NameSpaceRoleVals.NameSpaceRoleDetails.PolicyRules)
	PolicyRule := make([]rbacv1.PolicyRule, int(lenPL))
	for j := 0; j < lenPL; j++ {
		PolicyRule[j] = rbacv1.PolicyRule{
			Verbs: NameSpaceRoleVals.NameSpaceRoleDetails.PolicyRules[j].Verbs,
			APIGroups: NameSpaceRoleVals.NameSpaceRoleDetails.PolicyRules[j].APIGroups,
			Resources: NameSpaceRoleVals.NameSpaceRoleDetails.PolicyRules[j].Resources,
		}
	}

	//fmt.Println(PolicyRule)
	//Create or Update ServiceAccount
	for i := 0; i < lenNs; i++ {

		// NS Details

		fmt.Println("NameSpace Selected - ID: ", i)
		fmt.Println("NameSpace Selected - Name: ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("NameSpace Selected - Labels: ", mapLabels)

		//create sa
		sajson := v1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount",
				Labels: mapLabels,
				Namespace: NameSpaceVals.NameSpace[i].Name,
			},
		}

		// NS Details
		fmt.Println("ServiceAccount Name: ", NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount")
		fmt.Println("ServiceAccount Labels: ", mapLabels)
		fmt.Println("ServiceAccount NameSpace: ", NameSpaceVals.NameSpace[i].Name)

		CreateSerAcc, err := con.CoreV1().ServiceAccounts(NameSpaceVals.NameSpace[i].Name).Create(context.TODO(), &sajson, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			fmt.Println("Updating ServiceAccount.......")
			UpdateSerAcc, _ := con.CoreV1().ServiceAccounts(NameSpaceVals.NameSpace[i].Name).Update(context.TODO(), &sajson, metav1.UpdateOptions{})
			if err != nil{
				fmt.Println(err)
			} else {
				fmt.Println("Updated ServiceAccount: ", UpdateSerAcc.Name)
			}
		} else {
		fmt.Println("Created ServiceAccount: ", CreateSerAcc.Name)
		}

		// Attaching Role

		rolejson := rbacv1.Role{
				TypeMeta: metav1.TypeMeta{
				Kind:       "Role",
				APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
			ObjectMeta: metav1.ObjectMeta{
			Name:      NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-role",
			Labels: mapLabels,
			Namespace: NameSpaceVals.NameSpace[i].Name,
				},
			Rules: PolicyRule,
		}

		fmt.Println("Role Name: ", NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-role")
		fmt.Println("Role Labels: ", mapLabels)
		fmt.Println("Role NameSpace: ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("Role PolicyRule: ", PolicyRule)

			CreateRole, err := con.RbacV1().Roles(NameSpaceVals.NameSpace[i].Name).Create(context.TODO(), &rolejson, metav1.CreateOptions{})
			if err != nil {
				fmt.Println(err)
				fmt.Println("Updating Role........")
				UpateRole, _ := con.RbacV1().Roles(NameSpaceVals.NameSpace[i].Name).Update(context.TODO(), &rolejson, metav1.UpdateOptions{})
				if err != nil{
					fmt.Println(err)
				}
				fmt.Println("Updated Role: ", UpateRole.Name )
			} else {
				fmt.Println("Created Role: ", CreateRole.Name)
			}

			rolebdjson := rbacv1.RoleBinding{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-rolebinding",
					Labels: mapLabels,
					Namespace: NameSpaceVals.NameSpace[i].Name,
				},
				Subjects: []rbacv1.Subject{
					rbacv1.Subject{
						Kind:      "ServiceAccount",
						Namespace: NameSpaceVals.NameSpace[i].Name,
						Name:      NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount",
					},
				},
				RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind: "Role",
						Name: NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount",
				},
			}

		fmt.Println("RoleBinding Name: ", NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-rolebinding")
		fmt.Println("RoleBinding Labels: ", mapLabels)
		fmt.Println("RoleBinding NameSpace ", NameSpaceVals.NameSpace[i].Name)
		fmt.Println("RoleBinding ServiceAccount: ", NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount")
		fmt.Println("RoleBinding Role: ", NameSpaceVals.NameSpace[i].Name+"-"+NameSpaceRoleVals.NameSpaceRoleDetails.AppendName+"-sericeaccount")

			CreateRoleBinding, err := con.RbacV1().RoleBindings(NameSpaceVals.NameSpace[i].Name).Create(context.TODO(), &rolebdjson, metav1.CreateOptions{})
			if err != nil {
				fmt.Println(err)
				fmt.Println("Updating RoleBinding........")
				UpdateRoleBinding, _ := con.RbacV1().RoleBindings(NameSpaceVals.NameSpace[i].Name).Update(context.TODO(), &rolebdjson, metav1.UpdateOptions{})
				if err != nil{
					fmt.Println(err)
				} else {
					fmt.Println("Updated RoleBinding: ", UpdateRoleBinding.Name)
				}
				//return nil
			} else {
				fmt.Println("Created RoleBinding: ", CreateRoleBinding.Name)
			}
		}

	return err
}
func V1ResourceByStorageClass(storageClass string, resourceName v1.ResourceName) v1.ResourceName {
	return v1.ResourceName(string(storageClass + storageClassSuffix + string(resourceName)))
}