#!/bin/bash

######################## IMPORTANTE ########################################################
#  Antes de executar esse script, atualize os dados abaixo com as informações fornecidas pelo
#  terraform apply
############################################################################################
apigw_endpoint="https://kv2ze9k579.execute-api.us-east-1.amazonaws.com"
cognito_client_id="3gbr0s7qd4iaff6tivvk2kda7m"
cognito_url="https://hackaton-dev.auth.us-east-1.amazoncognito.com"
cognito_userpool_id="us-east-1_x7s3ZPD1w"


cognito_username="11122233300"
cognito_password="F@ap1234"


token=$(aws cognito-idp admin-initiate-auth --user-pool-id $cognito_userpool_id \
    --client-id $cognito_client_id \
    --auth-flow ADMIN_NO_SRP_AUTH \
    --auth-parameters USERNAME=$cognito_username,PASSWORD=$cognito_password \
| jq -r '.AuthenticationResult.AccessToken')


echo -e "\n-------------------------------TOKEN------------------------------------------"
echo "$token"

#test_endpoint="$apigw_endpoint/category/9764bd96-3bcf-11ee-be56-0242ac120002"
test_endpoint="$apigw_endpoint/hello"

test_body='{limit: 10, offset: 0}'

echo -e "\n-------------------------TEST WITHOUT TOKEN------------------------------------"

test_without_token="curl --location $test_endpoint -s \
--header 'Content-Type: application/json' \
--data '$test_body'"

eval "$test_without_token"


echo -e "\n-------------------------TEST WITH TOKEN-------------------------------------"

test_with_token="curl -X GET --location $test_endpoint \
--header 'Authorization: $token'  \
--header 'Content-Type: application/json' \
--data '$test_body'"

eval "$test_with_token"

# Swagger URL
echo -e "\n-------------------------SWAGGER URL-------------------------------------"
echo -e "Hackaton Swagger URL: $apigw_endpoint/swagger/index.html"


