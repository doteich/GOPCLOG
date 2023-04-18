package controller

import (
	"encoding/json"
	"strings"
)

type Service struct {
	ApiVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ServiceMetadata `json:"metadata"`
	Spec       ServiceSpec     `json:"spec"`
}

type ServiceMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"Namespace"`
	Labels    struct {
		App string `json:"app"`
		Id  string `json:"id"`
	} `json:"labels"`
}

type ServiceSpec struct {
	Ports    []Port `json:"ports"`
	Type     string `json:"type"`
	Selector struct {
		App string `json:"app"`
		Id  string `json:"id"`
	} `json:"selector"`
}

type Port struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
}

func SpawnService(podId string, data string) Service {
	var newContent DataContent
	json.Unmarshal([]byte(data), &newContent)

	id := strings.ReplaceAll(newContent.MethodConfig.Name, " ", "")

	socketPort := Port{
		Name:       "socket",
		Protocol:   "TCP",
		Port:       80,
		TargetPort: 8080,
	}

	miscPort := Port{
		Name:       "misc",
		Protocol:   "TCP",
		Port:       4444,
		TargetPort: 4444,
	}

	spec := ServiceSpec{Ports: []Port{miscPort, socketPort}, Type: "ClusterIP"}
	spec.Selector.App = "opcua-datalogger"
	spec.Selector.Id = id

	metadata := ServiceMetadata{Name: podId, Namespace: "default", Labels: Labels{App: "opcua-datalogger", Id: id}}

	newService := Service{ApiVersion: "v1", Kind: "Service", Metadata: metadata, Spec: spec}

	return newService
}
