package controller

type Service struct {
	ApiVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   Metadata    `json:"metadata"`
	Spec       ServiceSpec `json:"spec"`
}

type ServiceSpec struct {
	Selector string
}
