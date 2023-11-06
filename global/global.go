package global

import "github.com/gopcua/opcua"

var HttpOpcClient *opcua.Client

func SetOPCUAClient(mainClient *opcua.Client) {
	HttpOpcClient = mainClient
}