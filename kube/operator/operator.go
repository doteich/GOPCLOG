package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pod struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Spec struct {
	RestartPolicy string      `json:"restartPolicy"`
	Containers    []Container `json:"containers"`
	Volumes       []Volumes   `json:"volumes"`
}

type Container struct {
	Name         string         `json:"name"`
	Image        string         `json:"image"`
	VolumeMounts []VolumeMounts `json:"volumeMounts"`
}

type VolumeMounts struct {
	Name      string `json:"name"`
	Mountpath string `json:"mountPath"`
}

type Volumes struct {
	Name      string       `json:"name"`
	ConfigMap PodConfigMap `json:"configMap"`
}

type PodConfigMap struct {
	Name string `json:"name"`
}

type ConfigMap struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Data       map[string]string `json:"data"`
}

type ChildElements struct {
	Status   DesiredResourceNumbers `json:"status"`
	Children []interface{}          `json:"children"`
}

type DesiredResourceNumbers struct {
	Pods       int8 `json:"pods"`
	ConfigMaps int8 `json:"configmaps"`
}

type Request struct {
	Parent Parent `json:"parent"`
}

type Parent struct {
	Spec RequestSpec `json:"spec"`
}

type RequestSpec struct {
	Data string `json:"data"`
}

func main() {
	http.HandleFunc("/sync", sendPodData)
	err := http.ListenAndServe(":4900", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func sendPodData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("CONNECTION" + r.Method)
	if r.Method == "POST" {

		var body Request

		json.NewDecoder(r.Body).Decode(&body)
		fmt.Println(body)

		metadata := Metadata{Name: "opcua-datalogger", Namespace: "default"}
		container := Container{Name: "opcua-datalogger", Image: "cinderstries/opcua-logger", VolumeMounts: []VolumeMounts{{Name: "config-volume", Mountpath: "/etc/config"}}}
		volumes := Volumes{Name: "config-volume", ConfigMap: PodConfigMap{Name: "opcua-datalogger-cm"}}
		spec := Spec{RestartPolicy: "OnFailure", Containers: []Container{container}, Volumes: []Volumes{volumes}}

		newPod := Pod{ApiVersion: "v1", Kind: "Pod", Metadata: metadata, Spec: spec}

		data := make(map[string]string)
		data["config.json"] = body.Parent.Spec.Data

		newConfigmap := ConfigMap{ApiVersion: "v1", Kind: "ConfigMap", Metadata: Metadata{Name: "opcua-datalogger-cm", Namespace: "default"}, Data: data}

		var childs ChildElements
		childs.Children = append(childs.Children, newPod)
		childs.Children = append(childs.Children, newConfigmap)
		childs.Status.Pods = 1
		childs.Status.ConfigMaps = 1

		json, err := json.Marshal(childs)

		if err != nil {
			fmt.Println(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusBadRequest)

	}

}
