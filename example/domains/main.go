package main

import (
	"github.com/honeytrap/qemu"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
)

var (
	events = metrics.New(stdout.Stdout{})
)

func main() {

	domain := qemu.NewDomain(nil)

	vms, err := domain.ListVMS()
	if err != nil {
		events.Emit(stdout.Error("Failed to list Hypervisor vm").WithFields(metrics.Fields{
			"error":  err.Error(),
			"config": qemu.DefaultDomain,
		}))

		panic(err)
	}

	events.Emit(stdout.Info("List of Hypervisor VMs").WithFields(metrics.Fields{
		"vms":    vms,
		"config": qemu.DefaultDomain,
	}))

	if len(vms) == 0 {
		return
	}

	ubuntuContainer, err := domain.NewContainer(vms[0])
	if err != nil {
		events.Emit(stdout.Error("Failed to get vm domain container").WithFields(metrics.Fields{
			"error": err.Error(),
			"vm":    vms[0],
		}))

		panic(err)
	}

	res, err := ubuntuContainer.Network()
	if err != nil {
		events.Emit(stdout.Error("Failed to get vm started").WithFields(metrics.Fields{
			"error": err.Error(),
			"vm":    vms[0],
		}))

		panic(err)
	}

	events.Emit(stdout.Info("VM Network").WithFields(metrics.Fields{
		"response": string(res),
	}))

	if err := ubuntuContainer.Resume(); err != nil {
		events.Emit(stdout.Error("Failed to get vm to resume").WithFields(metrics.Fields{
			"error": err.Error(),
			"vm":    vms[0],
		}))

		panic(err)
	}
}
