package controller

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

func SpawnPod(podId string, secret string) Pod {

	var container Container

	if len(secret) > 0 {
		var userEnv EnvVar
		userEnv.Name = "OPCUA_USERNAME"
		userEnv.ValueFrom.SecretKeyRef.Name = secret
		userEnv.ValueFrom.SecretKeyRef.Key = "username"
		userEnv.ValueFrom.SecretKeyRef.Optional = false

		var pwEnv EnvVar
		pwEnv.Name = "OPCUA_PASSWORD"
		pwEnv.ValueFrom.SecretKeyRef.Name = secret
		pwEnv.ValueFrom.SecretKeyRef.Key = "password"
		pwEnv.ValueFrom.SecretKeyRef.Optional = false

		container = Container{Name: podId, Image: "cinderstries/opcua-logger", VolumeMounts: []VolumeMounts{{Name: "config-volume", Mountpath: "/etc/config"}}, Ports: []ContainerPorts{{ContainerPort: 4444}}, EnvVars: []EnvVar{userEnv, pwEnv}}
	} else {

		container = Container{Name: podId, Image: "cinderstries/opcua-logger", VolumeMounts: []VolumeMounts{{Name: "config-volume", Mountpath: "/etc/config"}}, Ports: []ContainerPorts{{ContainerPort: 4444}}}
	}

	metadata := Metadata{Name: podId, Namespace: "default", Labels: Labels{App: "opcua-datalogger"}}
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
