# Using Terraform

Terraform is a declarative way to create cloud infrastructure. These
simple scripts demonstrate how to use Terraform to create a BigQuery
dataset and table. In turn the services use this database to store and retrieve
information.

## Install Terraform



## Initialize


## Plan
```shell
terraform plan -out=plan.tf
```

## Apply
```shell
terraform apply plan.tf
```