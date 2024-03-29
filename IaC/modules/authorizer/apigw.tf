#-----------------------------------------------------------------------------------------------
#   VPC Link Config
#-----------------------------------------------------------------------------------------------

# Create App Security group
resource "aws_security_group" "vpc_link" {
  name   = "${var.project_name}-vpc-link"
  vpc_id = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Create the VPC Link configured with the ALB subnets. 
resource "aws_apigatewayv2_vpc_link" "api_vpc_link" {
  name               = "${var.project_name}-api_vpc_link"
  security_group_ids = [aws_security_group.vpc_link.id]
  subnet_ids         = var.private_subnet_ids
}




#-----------------------------------------------------------------------------------------------
#   API GW Config
#-----------------------------------------------------------------------------------------------

# Create the API Gateway HTTP endpoint
resource "aws_apigatewayv2_api" "api" {
  name          = "${var.project_name}-apigw"
  protocol_type = "HTTP"
}

# Set a default stage
resource "aws_apigatewayv2_stage" "apigw_stage" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
  depends_on  = [aws_apigatewayv2_api.api]
}


# Create authorizer pointing to Lambda (clear authorizer caching and remove identity resources)
resource "aws_apigatewayv2_authorizer" "api_authorizer" {
  name                              = "${var.project_name}_oAuth2Authorizer"
  api_id                            = aws_apigatewayv2_api.api.id
  authorizer_uri                    = aws_lambda_function.authorizer_lambda.invoke_arn
  identity_sources                  = []
  authorizer_type                   = "REQUEST"
  authorizer_payload_format_version = "2.0"
  authorizer_result_ttl_in_seconds  = 0
  enable_simple_responses           = true
}


# Search for the Load Balancer created by the K8s service for Hackaton micorservice
data "aws_lb" "eks_hackaton" {
  tags = {
    "kubernetes.io/service-name"                = "${var.project_name}/${var.lb_service_name_hackaton}"
    "kubernetes.io/cluster/${var.project_name}" = "owned"
  }
}

# Get the Listener of the Load Balancer created by this Load Balancer
data "aws_lb_listener" "eks_hackaton" {
  load_balancer_arn = data.aws_lb.eks_hackaton.arn
  port              = var.lb_service_port_hackaton
}

# Create the API Gateway HTTP_PROXY integration between the created API and the private load balancer via the VPC Link.
resource "aws_apigatewayv2_integration" "api_integration_hackaton" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "HTTP_PROXY"
  integration_uri        = data.aws_lb_listener.eks_hackaton.arn
  integration_method     = "ANY"
  connection_type        = "VPC_LINK"
  connection_id          = aws_apigatewayv2_vpc_link.api_vpc_link.id
  payload_format_version = "1.0"
  request_parameters = {
    "overwrite:header.userid" : "$context.authorizer.username"
  }
  depends_on = [aws_apigatewayv2_vpc_link.api_vpc_link,
    aws_apigatewayv2_api.api
  ]
}


# API GW route with ANY method
resource "aws_apigatewayv2_route" "apigw_route_hackaton" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "ANY /{proxy+}"
  target             = "integrations/${aws_apigatewayv2_integration.api_integration_hackaton.id}"
  authorization_type = "CUSTOM"
  authorizer_id      = aws_apigatewayv2_authorizer.api_authorizer.id
  depends_on         = [aws_apigatewayv2_integration.api_integration_hackaton]
}

# API GW route for /swagger/* for Hackaton
resource "aws_apigatewayv2_route" "apigw_route_swagger_hackaton" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "ANY /swagger/{proxy+}"
  target             = "integrations/${aws_apigatewayv2_integration.api_integration_hackaton.id}"
  authorization_type = "NONE"
  depends_on         = [aws_apigatewayv2_integration.api_integration_hackaton]
}
