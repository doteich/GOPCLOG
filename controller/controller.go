package controller

import (
	"log"

	"github.com/doteich/OPC-UA-Logger/global"
	"github.com/gopcua/opcua/ua"
)

type Controller struct {
	Rest bool
}

var EnabledControllers Controller

func InitControllers() {
	EnabledControllers.Rest = true
}

func WriteNode(nodeId string, value any) (*ua.WriteResponse, error) {
	id, err := ua.ParseNodeID(nodeId)
	if err != nil {
		return nil, err
	}

	v, err := ua.NewVariant(value)
	if err != nil {
		log.Fatalf("invalid value: %v", err)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        v,
				},
			},
		},
	}

	return global.HttpOpcClient.Write(req)
}
