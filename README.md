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
./infrastructure -o cluster -c examples/eks-cluster.yml
```
#### Setup Add-Ons
```
./infrastructure -o setup_addon -c examples/addon.yaml --context test-eks2
``` 

#### Examples

