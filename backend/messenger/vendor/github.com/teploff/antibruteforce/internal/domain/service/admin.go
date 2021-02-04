package service

import "net"

// AdminService provides interface for admin business-logic
//
// ResetBucketByLogin - flush bucket with passed login
//
// ResetBucketByPassword - flush bucket with passed password
//
// ResetBucketByIP - flush bucket with passed ip
//
// AddInBlacklist - adding subnet to the blacklist
//
// RemoveFromBlacklist - flushing subnet from the blacklist
//
// AddInWhitelist - adding subnet to the whitelist
//
// RemoveFromWhitelist - flushing subnet from the whitelist.
type AdminService interface {
	ResetBucketByLogin(login string) error
	ResetBucketByPassword(password string) error
	ResetBucketByIP(ip net.IP) error
	AddInBlacklist(ipNet *net.IPNet) error
	RemoveFromBlacklist(ipNet *net.IPNet) error
	AddInWhitelist(ipNet *net.IPNet) error
	RemoveFromWhitelist(ipNet *net.IPNet) error
}
