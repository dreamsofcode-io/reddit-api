module github.com/dreamsofcode-io/reddit-api/api

go 1.20

replace github.com/dreamsofcode-io/reddit-api/shared => ../shared

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.14.0
	github.com/dreamsofcode-io/reddit-api/shared v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.7.4
)

require (
	github.com/aws/aws-sdk-go-v2 v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.33 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.32 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.36 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.4.63 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.2 // indirect
	github.com/aws/smithy-go v1.14.1 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
