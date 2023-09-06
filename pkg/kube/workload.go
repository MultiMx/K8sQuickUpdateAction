package kube

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Workload struct {
	kube      Kube
	Namespace string
	Workload  string
}

func (a Workload) DeploymentUrl() string {
	return fmt.Sprintf(
		"%s/k8s/clusters/%s/apis/apps/v1/namespaces/%s/deployments/%s",
		a.kube.Conf.Backend,
		a.kube.Conf.Cluster,
		a.Namespace,
		a.Workload,
	)
}

func (a Workload) SetImage(image string) error {
	res, err := a.kube.Request("PATCH", &Request{
		Url:  a.DeploymentUrl(),
		Body: bytes.NewBuffer([]byte(fmt.Sprintf(`[{"op": "replace", "path": "/spec/template/spec/containers/%d/image", "value": "%s"}]`, a.kube.Conf.Container, image))),
	})
	if err != nil {
		return err
	}
	_ = res.Body.Close()
	return nil
}

func (a Workload) FullAvailable() (bool, error) {
	res, err := a.kube.Request("GET", &Request{
		Url: a.DeploymentUrl(),
	})
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var data struct {
		Status struct {
			Replicas          uint `json:"replicas"`
			AvailableReplicas uint `json:"availableReplicas"`
		} `json:"status"`
	}
	return data.Status.Replicas == data.Status.AvailableReplicas, json.NewDecoder(res.Body).Decode(&data)
}
