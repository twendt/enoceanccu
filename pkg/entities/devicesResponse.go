package entities

type DevicesResponse struct {
	InterfaceName string      `json:"interfaceName"`
	Devices       []DeviceDef `json:"devices"`
}

type DeviceDef struct {
	Type    string `json:"TYPE"`
	Address string `json:"ADDRESS"`
}
