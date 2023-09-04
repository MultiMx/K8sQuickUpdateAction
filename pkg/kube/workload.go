package kube

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (a Kube) SetImage(image string) error {
	res, err := a.Request("PATCH", &Request{
		Url:  a.Conf.DeploymentUrl(),
		Body: bytes.NewBuffer([]byte(fmt.Sprintf(`[{"op": "replace", "path": "/spec/template/spec/containers/%d/image", "value": "%s"}]`, a.Conf.Container, image))),
	})
	if err != nil {
		return err
	}
	_ = res.Body.Close()
	return nil
}

func (a Kube) DeploymentFullAvailable() (bool, error) {
	res, err := a.Request("GET", &Request{
		Url: a.Conf.DeploymentUrl(),
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
