# IaC Project

This project uses Terraform to provision and manage infrastructure on AWS. It sets up a VPC, an RDS database, an EKS cluster, a Kubernetes application, and an authorizer module.

## Modules

- `vpc_for_eks`: This module sets up a VPC for the EKS cluster.
- `rds`: This module sets up an RDS database.
- `eks_cluster`: This module sets up an EKS cluster.
- `app`: This module deploys a Kubernetes application.
- `authorizer`: This module sets up an authorizer with Cognito, Lambda, and API Gateway.

## Usage

Before running this project, you need to initialize your Terraform workspace, which will download the provider plugins for AWS:

```bash
# configure the aws cli auth (~/.aws/credentials)
# check the ./aws-setup.sh to configure the bucket and github actions
./aws-setup.sh
# run the last command that will be printed by the script (terraform init)
```

After initializing your workspace, you can run the following commands:

```bash
# check the ./actions.sh to see the available commands to simulate github actions locally
./actions.sh
```

or just run the following command to apply the infrastructure:

```bash
terraform apply
```
