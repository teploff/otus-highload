package service

import (
	"net"

	"github.com/teploff/antibruteforce/internal/domain/entity"
)

// AuthService provides interface for authorization business-logic
//
// LogIn - attempting to log in.
type AuthService interface {
	LogIn(credentials entity.Credentials, ip net.IP) (bool, error)
}
