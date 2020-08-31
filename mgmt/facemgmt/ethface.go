package facemgmt

import "net"

package facemgmt

import (
"errors"

"github.com/rackn/gohai/plugins/net"
"github.com/usnistgov/ndn-dpdk/dpdk/eal"
"github.com/usnistgov/ndn-dpdk/dpdk/ethdev"
"github.com/usnistgov/ndn-dpdk/iface/ethface"
)

type EthFaceMgmt struct{}

func (EthFaceMgmt) ListPorts(args struct{}, reply *[]PortInfo) error {
	result := make([]PortInfo, 0)
	for _, dev := range ethdev.List() {
		result = append(result, makePortInfo(dev))
	}
	*reply = result
	return nil
}

func (EthFaceMgmt) ListPortFaces(args PortArg, reply *[]BasicInfo) error {
	dev := ethdev.Find(args.Port)
	if !dev.Valid() {
		return errors.New("EthDev not found")
	}

	result := make([]BasicInfo, 0)
	if port := ethface.FindPort(dev); port != nil {
		for _, face := range port.Faces() {
			result = append(result, makeBasicInfo(face))
		}
	}
	*reply = result
	return nil
}

func (EthFaceMgmt) ReadPortStats(args PortStatsArg, reply *ethdev.Stats) error {
	dev := ethdev.Find(args.Port)
	if !dev.Valid() {
		return errors.New("EthDev not found")
	}

	*reply = dev.Stats()
	if args.Reset {
		dev.ResetStats()
	}
	return nil
}

type PortArg struct {
	Port string
}

type PortInfo struct {
	Name       string           // port name
	NumaSocket eal.NumaSocket   // NUMA socket
	MacAddr    net.HardwareAddr // MAC address
	Active     bool             // whether port is active
	ImplName   string           // internal implementation name
}

func makePortInfo(dev ethdev.EthDev) (info PortInfo) {
	info.Name = dev.Name()
	info.NumaSocket = dev.NumaSocket()
	info.MacAddr = net.HardwareAddr(dev.MacAddr())
	port := ethface.FindPort(dev)
	if port != nil {
		info.Active = true
		info.ImplName = port.ImplName()
	}
	return info
}

type PortStatsArg struct {
	PortArg
	Reset bool
}
