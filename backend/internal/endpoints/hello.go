package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type (
	HelloEndpoints struct {
		GetHelloEndpoint endpoint.Endpoint
	}
)

func MakeHelloEndpoints() HelloEndpoints {
	return HelloEndpoints{
		GetHelloEndpoint: makeGetHelloEndpoint(),
	}
}

func makeGetHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "Hello, World!", nil
	}
}
