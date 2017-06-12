package qemu

import (
	"net"
	"time"

	"github.com/digitalocean/go-qemu/hypervisor"
	"github.com/digitalocean/go-qemu/qemu"
	"github.com/digitalocean/go-qemu/qmp"
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

// ListVMS returns the list of all VMS available in the hypervisor.
func (dm *Domain) ListVMS() ([]string, error) {
	return dm.hypervisor.DomainNames()
}

// NewContainer returns a new DomainContainer for a giving hypervisor vm.
func (dm *Domain) NewContainer(domainName string) (*DomainContainer, error) {
	domain, err := dm.hypervisor.Domain(domainName)
	if err != nil {
		return nil, err
	}

	return &DomainContainer{
		name:   domainName,
		domain: domain,
	}, nil
}

// DomainContainer defines a structure which interacts with a given VM
// domain object.
type DomainContainer struct {
	name   string
	domain *qemu.Domain
}

// Domain returns the underline qemu.Domain object for this container.
func (dmc *DomainContainer) Domain() *qemu.Domain {
	return dmc.domain
}

// Name returns the giving name associated with the domain vm.
func (dmc *DomainContainer) Name() string {
	return dmc.name
}

// Devices returns the giving list of available devices for the domain vm.
func (dmc *DomainContainer) Devices() ([]qemu.BlockDevice, []qemu.PCIDevice, error) {
	pcis, err := dmc.domain.PCIDevices()
	if err != nil {
		return nil, nil, err
	}

	blocks, err := dmc.domain.BlockDevices()
	if err != nil {
		return nil, nil, err
	}

	return blocks, pcis, nil
}

// Status returns the domain vm current status.
func (dmc *DomainContainer) Status() (qemu.Status, error) {
	return dmc.domain.Status()
}

// Network returns network details related to the vm.
func (dmc *DomainContainer) Network() ([]byte, error) {
	// if status, err := dmc.domain.Status
	res, err := dmc.domain.Run(qmp.Command{
		Execute: "query-vnc",
	})

	return res, err
}

// Wakeup activtes the domain if wakeup and starts the domains vm for
// interaction.
func (dmc *DomainContainer) Wakeup() error {
	// if status, err := dmc.domain.Status
	_, err := dmc.domain.Run(qmp.Command{
		Execute: "system_wakeup",
	})

	return err
}

// Resume activtes the domain if paused and resumes the domains vm for
// interaction.
func (dmc *DomainContainer) Resume() error {
	// if status, err := dmc.domain.Status
	_, err := dmc.domain.Run(qmp.Command{
		Execute: "cont",
	})

	return err
}

// Start activtes the domain if not running and starts the domains vm for
// interaction.
func (dmc *DomainContainer) Start() error {
	// if status, err := dmc.domain.Status
	_, err := dmc.domain.Run(qmp.Command{
		Execute: "system_wakeup",
	})

	return err
}

// Reset resets the domain and its state.
func (dmc *DomainContainer) Reset() error {
	return dmc.domain.SystemReset()
}

// Stop Poweroff the giving domain.
func (dmc *DomainContainer) Stop() error {
	return dmc.domain.SystemPowerdown()
}

// NetworkAddress returns the associated network address for the vm.
func (dmc *DomainContainer) NetworkAddress() string {
	var addr string

	return addr
}
