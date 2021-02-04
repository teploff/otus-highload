package auth

import (
	"context"
	"net"

	"github.com/go-kit/kit/endpoint"
	"github.com/teploff/antibruteforce/internal/domain/entity"
	"github.com/teploff/antibruteforce/internal/domain/service"
)

// Endpoints for authorization.
type Endpoints struct {
	SignIn endpoint.Endpoint
}

func makeSignInEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignInRequest)
		ok, err := svc.LogIn(entity.Credentials{
			Login:    req.Login,
			Password: req.Password,
		}, net.ParseIP(req.IP))

		if err != nil {
			return SignInResponse{}, err
		}

		return SignInResponse{
			Ok: ok,
		}, nil
	}
}

// MakeAuthEndpoints for authorization.
func MakeAuthEndpoints(svc service.AuthService) Endpoints {
	return Endpoints{
		SignIn: makeSignInEndpoint(svc),
	}
}
