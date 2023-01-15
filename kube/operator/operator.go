package main

import (
	"encoding/json"
	"fmt"
	"gopc_operator/controller"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ChildElements struct {
	Status   DesiredResourceNumbers `json:"status"`
	Children []interface{}          `json:"children"`
}

type DesiredResourceNumbers struct {
	Pods            int8 `json:"pods"`
	ConfigMaps      int8 `json:"configmaps"`
	Services        int8 `json:"services"`
	ServiceMonitors int8 `json:"servicemonitors"`
}

type Request struct {
	Parent   Parent   `json:"parent"`
	Children Children `json:"children"`
}

type Children struct {
	Pods map[string]string `json:"Pod.v1"`
}

type Parent struct {
	Spec RequestSpec `json:"spec"`
}

type RequestSpec struct {
	Data      string `json:"data"`
	SecretRef string `json:"secretRef"`
}

func main() {
	http.HandleFunc("/sync", sendPodData)
	err := http.ListenAndServe(":4900", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func sendPodData(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		var body Request

		json.NewDecoder(r.Body).Decode(&body)

		fmt.Println(body)

		if len(body.Children.Pods) > 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		} else {
			podName := "opcua-datalogger" + setId()

			secret := body.Parent.Spec.SecretRef

			newPod := controller.SpawnPod(podName, secret)
			newConfigmap := controller.SpawnCM(body.Parent.Spec.Data, podName)

			var childs ChildElements
			childs.Children = append(childs.Children, newPod)
			childs.Children = append(childs.Children, newConfigmap)

			childs.Status.Pods = 1
			childs.Status.ConfigMaps = 1
			childs.Status.Pods = 1
			childs.Status.ServiceMonitors = 1

			json, err := json.Marshal(childs)

			if err != nil {
				fmt.Println(err)
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)

	}

}

func setId() string {
	randNum := strconv.Itoa(rand.Intn(1000))
	day := strconv.Itoa(time.Now().Day())
	year := strconv.Itoa(time.Now().Year())

	id := day + year + randNum
	return id
}
