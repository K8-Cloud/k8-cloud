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

###Commands:
#### Setup EKS Cluster
```
./k8-cloud -o cluster -c examples/eks-cluster.yml
```
#### Setup Add-Ons
```
./k8-cloud -o addons -c examples/addon.yaml --context test-eks5
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

## 23/02/2021
* Azure Support
