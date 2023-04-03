package exporter

import (
	"errors"
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
	username       string
	password       string
	authEnabled    bool
)

func InitHTTPServer() {
	if EnabledExporters.Rest {
		http.HandleFunc("/triggerread", ReadFromOPC)

		if setup.PubConfig.ExporterConfig.Rest.AuthType == "Basic" {
			authEnabled = true
			username = setup.PubConfig.ExporterConfig.Rest.Username
			password = setup.PubConfig.ExporterConfig.Rest.Password
		}

	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)
}

func SetOPCUAClient(mainClient *opcua.Client) {
	http_opcclient = mainClient
}

func ReadFromOPC(w http.ResponseWriter, r *http.Request) {

	if http_opcclient == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {

		if authEnabled {
			u, p, ok := r.BasicAuth()

			if !ok {
				logging.LogGeneric("warning", "error parsing basic auth credentials", "exporter")
				w.WriteHeader(401)
				return
			}

			err := validateAuth(u, p)

			if err != nil {
				logging.LogError(err, "unauthorized", "exporter")
				w.WriteHeader(401)
				return
			}

		}

		w.WriteHeader(http.StatusOK)
		nodes := TriggerBulkRead()

		for key, val := range nodes {
			nodeObj, err := findNodeDetails(key)

			if err != nil {
				logging.LogError(err, "node not found while reading from opc server", "exporter")
				continue
			}

			dataType, _ := InferDataType(val)
			http_exporter.PostLoggedData(key, nodeObj.NodeName, val, time.Now(), setup.PubConfig.LoggerConfig.Name, setup.PubConfig.ClientConfig.Url, dataType)
		}
		return

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func validateAuth(user string, pw string) error {
	if user != username {
		return errors.New("wrong username")
	}
	if pw != password {
		return errors.New("wrong password")
	}
	return nil

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
