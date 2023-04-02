package exporter

import (
	"net/http"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/http_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	http_opcclient *opcua.Client
)

func InitHTTPServer() {
	if EnabledExporters.Rest {
		http.HandleFunc("/triggerread", ReadFromOPC)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)
}

func SetOPCUAClient(mainClient *opcua.Client) {
	http_opcclient = mainClient
}

func ReadFromOPC(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if http_opcclient == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			nodes := TriggerBulkRead()

			for key, val := range nodes {
				nodeObj, _ := findNodeDetails(key)
				dataType, _ := InferDataType(val)
				http_exporter.PostLoggedData(key, nodeObj.NodeName, val, time.Now(), setup.PubConfig.LoggerConfig.Name, setup.PubConfig.ClientConfig.Url, dataType)
			}
			return
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func ReadNodes(nodeId string) (interface{}, error) {

	id, _ := ua.ParseNodeID(nodeId)

	obj := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
	}
	resp, err := http_opcclient.Read(obj)

	if err != nil {
		logging.LogError(err, "Error while reading"+nodeId, "exporter")
		return nil, err
	}

	return resp.Results[0].Value.Value(), nil

}

func TriggerBulkRead() map[string]interface{} {

	idMap := make(map[string]interface{})

	for _, node := range setup.PubConfig.Nodes {

		val, err := ReadNodes(node.NodeId)

		if err != nil {
			continue
		}

		idMap[node.NodeId] = val

	}

	return idMap
}
