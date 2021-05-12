# k8-cloud
Multi Cloud K8s CLuster Setup


# Enterprise Version Support:
* Helm Backup and Restore -- Needs Testing
* SSO LDAP (Dex) Authentication
* Prime Support Engineer
* Dynamics number of subnets -- TODO


# Supported in Open Source 
* Cluster Creation
* Canary Upgrades
* AWS, AKS and GKE Support
* Encryption at Rest
* Addons Deployment with helm
* Support 3 Private and Public subnets max

##Commands:
### Setup EKS Cluster
```
./k8-cloud --operation cluster --config examples/eks-cluster.yml
```
### get cluster config on local
it will create kubeconfig under ~/.kube/config or if there is existing file it will do a safe merge with existing contexts
```
aws eks update-kubeconfig --name <cluster_name> --alias <alias_name>
```

### Setup Add-Ons
```
./k8-cloud --operation addons --config examples/addon.yaml --context test-eks5

``` 

### Init Cluster Management
```
./k8-cloud --operation init --context test-eks9
```

### Setup namespace
```
./k8-cloud --operation namespace --context test-eks9
```

### Setup Resource Quota
```
./k8-cloud --operation resourcequota --context test-eks9
```






#### Examples

#### TODO
1. restrict control plane with cidr
1. Lable nodes
2. Namespace Quota
   * resource limits
   * pods 50
   * storage 50Gi
3. a chart
    * create namespace
    * user with admin access
    * user with readonly access
4. k8s netwok policies
5. Compare CFT checksum before apply
    




### 17/02/2021
* version command -- DONE
* backup take option with config file -- Planning how to implement  
* subnets support 4, 6, 8 with fixed CFT Samples
* eks cluster creation yaml doc

## 03/04/2021
* add namespace management in the operations
* documentations








## 23/02/2021
* Azure Support
