#!/bin/bash

# Set the GitHub Actions inputs
github_event_inputs_TF_STATE_BUCKET="default-with-account-id"
github_event_inputs_svc_hackaton_image_tag="branch"
github_event_inputs_DEPLOY_ENVIRONMENT="branch"
github_event_name=""
# get the current head ref and ref
github_head_ref=$(git rev-parse --abbrev-ref HEAD)
github_ref=$(git rev-parse --abbrev-ref HEAD)


# Set the environment variables
if [ "${github_event_inputs_TF_STATE_BUCKET}" == "default-with-account-id" ]; then
  # get the AWS account ID
  AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
  export TF_STATE_BUCKET="s1sg-tc-terraform-state-bucket-${AWS_ACCOUNT_ID}"
else
  export TF_STATE_BUCKET="${github_event_inputs_TF_STATE_BUCKET}"
fi
if [ "${github_event_inputs_svc_hackaton_image_tag}" == "branch" ]; then
  if [ "${github_event_name}" == "pull_request" ]; then
    BRANCH_NAME=$(echo "${github_head_ref}" | awk -F"/" '{print $NF}' | sed 's/[^a-zA-Z0-9]/_/g')
    export svc_hackaton_image_tag="${BRANCH_NAME}"
  else
    BRANCH_NAME=$(echo "${github_ref}" | awk -F"/" '{print $NF}' | sed 's/[^a-zA-Z0-9]/_/g')
    export svc_hackaton_image_tag="${BRANCH_NAME}"
  fi
else
  export svc_hackaton_image_tag="${github_event_inputs_svc_hackaton_image_tag}"
fi
if [ "${github_event_inputs_DEPLOY_ENVIRONMENT}" == "branch" ]; then
  export DEPLOY_ENVIRONMENT="${BRANCH_NAME}"
else
  export DEPLOY_ENVIRONMENT="${github_event_inputs_DEPLOY_ENVIRONMENT}"
fi

# Print the environment variables
echo "TF_STATE_BUCKET=${TF_STATE_BUCKET}"
echo "svc_hackaton_image_tag=${svc_hackaton_image_tag}"
echo "DEPLOY_ENVIRONMENT=${DEPLOY_ENVIRONMENT}"

# Initialize Terraform
terraform init -upgrade -backend-config="bucket=${TF_STATE_BUCKET}" -backend-config="key=terraform.tfstate" -backend-config="region=us-east-1"

# Set the redeploy annotation
export TF_VAR_redeploy_annotation=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Apply the Terraform configuration
#export TF_VAR_environment="${DEPLOY_ENVIRONMENT}"
export TF_VAR_environment="develop"
export TF_VAR_image_registry="ghcr.io/soat1stackgolang"
#export TF_VAR_svc_hackaton_image_tag="${svc_hackaton_image_tag}"
export TF_VAR_svc_hackaton_image_tag="develop"
#terraform plan
#terraform apply -auto-approve
#terraform refresh
terraform destroy -auto-approve