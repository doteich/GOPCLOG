package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
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

type Spec struct {
	RestartPolicy string      `json:"restartPolicy"`
	Containers    []Container `json:"containers"`
	Volumes       []Volumes   `json:"volumes"`
}

type Volumes struct {
	Name      string    `json:"name"`
	ConfigMap ConfigMap `json:"configMap"`
}

type ConfigMap struct {
	Name string `json:"name"`
}

type Pod struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type ChildElements struct {
	Status   DesiredPodNumber `json:"status"`
	Children []Pod            `json:"children"`
}

type DesiredPodNumber struct {
	Pods int8 `json:"pods"`
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

		var body interface{}

		json.NewDecoder(r.Body).Decode(&body)

		fmt.Println(body)

		metadata := Metadata{Name: "opcua-datalogger", Namespace: "default"}
		container := Container{Name: "opcua-datalogger", Image: "cinderstries/opcua-logger", VolumeMounts: []VolumeMounts{{Name: "config-volume", Mountpath: "/etc/config"}}}
		volumes := Volumes{Name: "config-volume", ConfigMap: ConfigMap{Name: "log-config"}}
		spec := Spec{RestartPolicy: "OnFailure", Containers: []Container{container}, Volumes: []Volumes{volumes}}

		newPod := Pod{ApiVersion: "v1", Kind: "Pod", Metadata: metadata, Spec: spec}

		var childs ChildElements
		childs.Children = append(childs.Children, newPod)
		childs.Status.Pods = 1

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
