package entities

type Device struct {
	RcvEEPs     []string `json:"rcv-eeps"`
	SendEEPs    []string `json:"send-eeps"`
	HMType      string   `json:"hm-type"`
	HMAddress   string   `json:"hm-address"`
	LocalSendID string   `json:"local-send-id"`
	ID          string   `json:"id"`
}
