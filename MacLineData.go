package main

type MacLineData struct {
	vlan  string
	mac   string
	iface string
}

type HostMacLineData struct {
	HostName string
	mld      []MacLineData
}

func NewHostMacLineData(hostname string) *HostMacLineData {
	return &HostMacLineData{
		HostName: hostname,
		mld:      []MacLineData{},
	}
}

func NewMacLineData(
	vlan string,
	mac string,
	iface string,
) *MacLineData {

	return &MacLineData{
		vlan:  vlan,
		mac:   mac,
		iface: iface,
	}
}

/*
func (m *MacLineData) PrintData() {
	fmt.Printf("Vlan: %v\n", m.vlan)
	fmt.Printf("MAC: %v\n", m.mac)
	fmt.Printf("Port: %v\n", m.iface)
}
*/
