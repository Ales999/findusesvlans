package main

import "fmt"

type MacLineData struct {
	// Vlan    Mac Address       Type        Ports
	vlan     string
	mac      string
	types    string
	iface    string
	hostname string
}

func NewMacLineData(
	vlan string,
	mac string,
	types string,
	iface string,
	hostname string,
) *MacLineData {

	return &MacLineData{
		vlan:     vlan,
		mac:      mac,
		types:    types,
		iface:    iface,
		hostname: hostname,
	}

}

func (m *MacLineData) PrintData() {
	fmt.Printf("Vlan: %v\n", m.vlan)
	fmt.Printf("MAC: %v\n", m.mac)
	fmt.Printf("Port: %v\n", m.iface)
}
