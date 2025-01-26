package main

import (
	"fmt"
	"log"
	"net"

	"antrea.io/libOpenflow/openflow15"
	"antrea.io/ofnet/ofctrl"
)

type Sdnc struct {
	Switch *ofctrl.OFSwitch
}

func (o *Sdnc) PacketRcvd(sw *ofctrl.OFSwitch, packet *ofctrl.PacketIn) {

	log.Printf("App: Received packet: %+v", packet)
}

func (o *Sdnc) SwitchConnected(sw *ofctrl.OFSwitch) {
	log.Printf("App: Switch connected: %v", sw.DPID())

	// Create all tables
	rxVlanTbl, err := sw.NewTable(1)
	macSaTable, err := sw.NewTable(2)
	macDaTable, err := sw.NewTable(3)
	ipTable, err := sw.NewTable(4)
	inpTable := sw.DefaultTable() // table 0. i.e starting table

	_ = err

	// Discard mcast source mac
	dscrdMcastSrc, err := inpTable.NewFlow(ofctrl.FlowMatch{
		&McastSrc:     {0x01, 0, 0, 0, 0, 0},
		&McastSrcMask: {0x01, 0, 0, 0, 0, 0},
	}, 100)
	dscrdMcastSrc.Next(sw.DropAction())

	// All valid packets go to vlan table
	validInputPkt := inpTable.NewFlow(FlowMatch{}, 1)
	validInputPkt.Next(rxVlanTbl)

	// Set access vlan for port 1 and go to mac lookup
	tagPort := rxVlanTbl.NewFlow(FlowMatch{
		InputPort: Port(1),
	}, 100)
	tagPort.SetVlan(10)
	tagPort.Next(macSaTable)

	// Match on IP dest addr and forward to a port
	ipFlow := ipTable.NewFlow(FlowParams{
		Ethertype: 0x0800,
		IpDa:      &net.IPv4("10.10.10.10"),
	}, 100)

	outPort := sw.NewOutputPort(10)
	ipFlow.Next(outPort)
}

func (o *Sdnc) FlowGraphEnabledOnSwitch() bool {
	return true
}

func (o *Sdnc) SwitchDisconnected(sw *ofctrl.OFSwitch) {
	log.Printf("App: Switch disconnected: %v", sw.DPID())
}

func (o *Sdnc) MultipartReply(sw *ofctrl.OFSwitch, rep *openflow15.MultipartReply) {
	fmt.Println("multipart reply")
}

func (o *Sdnc) TLVMapEnabledOnSwitch() bool {
	return true
}

// PortStatusRcvd notifies AppInterface a new PortStatus message is received.
func (o *Sdnc) PortStatusRcvd(status *openflow15.PortStatus) {

	fmt.Println("new port status")
	fmt.Println("\tmac:", status.Desc.HWAddr.String())
	fmt.Println("\tport:", status.Desc.PortNo)
}

func main() {
	var app Sdnc

	ctrler := ofctrl.NewController(&app)

	// Listen for connections
	ctrler.Listen(":6633")
}
