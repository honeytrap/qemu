package qemu

import (
	"net"
	"time"

	"github.com/digitalocean/go-qemu/hypervisor"
)

var (
	defaultTimeout = 3 * time.Second
)

//=====================================================================================================

// DriverCType defines the possible driver connection type to be used
// for the Domain.
type DriverCType int

// contains two available DriverCType for connecting with the qemu hypervisor.
const (
	RPCCDriver DriverCType = iota
	SocketCDriver
)

//=====================================================================================================

var (
	// DefaultDomain defines a default domain configuration for use when
	// a configuration is not giving.
	DefaultDomain = DomainConfig{
		Network: "unix",
		Address: "/var/run/libvirt/libvirt-sock",
		Timeout: defaultTimeout,
		DCType:  RPCCDriver,
	}
)

// DomainConfig defines the configuration used to build drive a Domain
// internal settings.
type DomainConfig struct {
	Network   string
	Address   string
	Timeout   time.Duration
	DCType    DriverCType
	Addresses []hypervisor.SocketAddress
}

//=====================================================================================================

// Domain defines a single qemu vm which is used to interface and
// operate with the vm for interaction.
type Domain struct {
	config     DomainConfig
	driver     hypervisor.Driver
	hypervisor *hypervisor.Hypervisor
}

// NewDomain returns a new Domain instance which allows interaction with the
// installed qemu hypervisor.
func NewDomain(config *DomainConfig) *Domain {
	var dm Domain

	if config == nil {
		dm.config = DefaultDomain
	} else {
		dm.config = *config
	}

	switch dm.config.DCType {
	case RPCCDriver:
		dm.driver = hypervisor.NewRPCDriver(func() (net.Conn, error) {
			return net.DialTimeout(dm.config.Network, dm.config.Address, dm.config.Timeout)
		})

		break
	case SocketCDriver:
		addresses := append([]hypervisor.SocketAddress{
			{
				Network: dm.config.Network,
				Address: dm.config.Address,
				Timeout: dm.config.Timeout,
			},
		}, dm.config.Addresses...)

		dm.driver = hypervisor.NewSocketDriver(addresses)
		break
	}

	dm.hypervisor = hypervisor.New(dm.driver)

	return &dm
}
