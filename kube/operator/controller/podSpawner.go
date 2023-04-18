package controller

import (
	"encoding/json"
	"strings"
)

type OPCConfig struct {
	SecretRef string `json:"secretRef"`
}

type MethodConfig struct {
	Name string `json:"name"`
}

type MongoDB struct {
	SecretRef string `json:"secretRef"`
}

type Rest struct {
	SecretRef string `json:"secretRef"`
}

type Exporters struct {
	MongoDB MongoDB `json:"mongodb"`
	Rest    Rest    `json:"rest"`
}

type DataContent struct {
	OPCConfig    OPCConfig    `json:"opcConfig"`
	Exporters    Exporters    `json:"exporters"`
	MethodConfig MethodConfig `json:"methodConfig"`
}

type Pod struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Labels    Labels `json:"labels"`
}

type Labels struct {
	App string `json:"app"`
	Id  string `json:"id"`
}

type Spec struct {
	RestartPolicy string      `json:"restartPolicy"`
	Containers    []Container `json:"containers"`
	Volumes       []Volumes   `json:"volumes"`
}

type Container struct {
	Name         string           `json:"name"`
	Image        string           `json:"image"`
	VolumeMounts []VolumeMounts   `json:"volumeMounts"`
	Ports        []ContainerPorts `json:"ports"`
	EnvVars      []EnvVar         `json:"env"`
}

type VolumeMounts struct {
	Name      string `json:"name"`
	Mountpath string `json:"mountPath"`
}

type ContainerPorts struct {
	ContainerPort int `json:"containerPort"`
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

type EnvVar struct {
	Name      string `json:"name"`
	ValueFrom struct {
		SecretKeyRef struct {
			Name     string `json:"name"`
			Key      string `json:"key"`
			Optional bool   `json:"optional"`
		} `json:"secretKeyRef"`
	} `json:"valueFrom"`
}

func SpawnPod(podId string, data string) Pod {

	envs, id := setVarPodData(data)

	container := Container{Name: podId, Image: "cinderstries/opcua-logger", VolumeMounts: []VolumeMounts{{Name: "config-volume", Mountpath: "/etc/config"}}, Ports: []ContainerPorts{{ContainerPort: 4444}}, EnvVars: envs}

	metadata := Metadata{Name: podId, Namespace: "default", Labels: Labels{App: "opcua-datalogger", Id: id}}
	volumes := Volumes{Name: "config-volume", ConfigMap: PodConfigMap{Name: podId + "-cm"}}
	spec := Spec{RestartPolicy: "OnFailure", Containers: []Container{container}, Volumes: []Volumes{volumes}}

	newPod := Pod{ApiVersion: "v1", Kind: "Pod", Metadata: metadata, Spec: spec}

	return newPod
}

func SpawnCM(config string, podId string) ConfigMap {
	data := make(map[string]string)
	data["config.json"] = config
	newConfigmap := ConfigMap{ApiVersion: "v1", Kind: "ConfigMap", Metadata: Metadata{Name: podId + "-cm", Namespace: "default"}, Data: data}
	return newConfigmap
}

func setVarPodData(data string) ([]EnvVar, string) {

	var newContent DataContent

	json.Unmarshal([]byte(data), &newContent)

	envs := make([]EnvVar, 0)
	var env EnvVar

	if newContent.OPCConfig.SecretRef != "" {

		env.Name = "OPCUA_USERNAME"
		env.ValueFrom.SecretKeyRef.Name = newContent.OPCConfig.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "username"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)

		env.Name = "OPCUA_PASSWORD"
		env.ValueFrom.SecretKeyRef.Name = newContent.OPCConfig.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "password"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)
	}

	if newContent.Exporters.MongoDB.SecretRef != "" {

		env.Name = "MONGODB_USERNAME"
		env.ValueFrom.SecretKeyRef.Name = newContent.Exporters.MongoDB.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "username"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)

		env.Name = "MONGODB_PASSWORD"
		env.ValueFrom.SecretKeyRef.Name = newContent.Exporters.MongoDB.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "password"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)
	}

	if newContent.Exporters.Rest.SecretRef != "" {
		env.Name = "REST_USERNAME"
		env.ValueFrom.SecretKeyRef.Name = newContent.Exporters.Rest.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "username"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)

		env.Name = "REST_PASSWORD"
		env.ValueFrom.SecretKeyRef.Name = newContent.Exporters.Rest.SecretRef
		env.ValueFrom.SecretKeyRef.Key = "password"
		env.ValueFrom.SecretKeyRef.Optional = false

		envs = append(envs, env)
	}

	return envs, strings.ReplaceAll(newContent.MethodConfig.Name, " ", "")
}
