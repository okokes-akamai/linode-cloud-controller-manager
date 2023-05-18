package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/linode/linodego"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

const providerIDPrefix = "linode://"

type invalidProviderIDError struct {
	value string
}

func (e invalidProviderIDError) Error() string {
	return fmt.Sprintf("invalid provider ID %q", e.value)
}

func parseProviderID(providerID string) (int, error) {
	if !strings.HasPrefix(providerID, providerIDPrefix) {
		return 0, invalidProviderIDError{providerID}
	}
	id, err := strconv.Atoi(strings.TrimPrefix(providerID, providerIDPrefix))
	if err != nil {
		return 0, invalidProviderIDError{providerID}
	}
	return id, nil
}

func linodeFilterListOptions(targetLabel string) *linodego.ListOptions {
	jsonFilter := fmt.Sprintf(`{"label":%q}`, targetLabel)
	return linodego.NewListOptions(0, jsonFilter)
}

func linodeByName(ctx context.Context, client Client, nodeName types.NodeName) (*linodego.Instance, error) {
	log.Printf("PERF: linodeByName: %v", nodeName)
	linodes, err := client.ListInstances(ctx, linodeFilterListOptions(string(nodeName)))
	if err != nil {
		return nil, err
	}

	if len(linodes) == 0 {
		return nil, cloudprovider.InstanceNotFound
	} else if len(linodes) > 1 {
		return nil, fmt.Errorf("Multiple instances found with name %v", nodeName)
	}

	return &linodes[0], nil
}

func linodeByID(ctx context.Context, client Client, id int) (*linodego.Instance, error) {
	log.Printf("PERF: linodeByID: %v", id)
	instance, err := client.GetInstance(ctx, id)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, fmt.Errorf("linode not found with id %v", id)
	}
	return instance, nil
}
