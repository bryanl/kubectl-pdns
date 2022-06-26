package pdns

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/multierr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeleteConfig is configuration for the delete command.
type DeleteConfig struct {
	Zone      string
	Name      string
	Type      string
	Namespace string
}

// Validate validates the configuration.
func (c *DeleteConfig) Validate() error {
	var err error

	if c.Zone == "" {
		err = multierr.Append(err, fmt.Errorf("zone is blank"))
	}

	if c.Name == "" {
		err = multierr.Append(err, fmt.Errorf("name is blank"))
	}

	if c.Type == "" {
		err = multierr.Append(err, fmt.Errorf("type is blank"))
	}

	if c.Namespace == "" {
		c.Namespace = "default"
	}

	return err
}

// ObjectName returns the name for the record based on the configuration.
func (c *DeleteConfig) ObjectName() string {
	return fmt.Sprintf("%s-%s", c.Name, strings.ToLower(c.Type))
}

// Delete deletes a record.
func Delete(config DeleteConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("validate configuration: %w", err)
	}

	client, err := buildClient()
	if err != nil {
		return fmt.Errorf("build client: %w", err)
	}

	res := schema.GroupVersionResource{Group: "pdns.bryanl.dev", Version: "v1alpha1", Resource: "records"}

	ctx := context.Background()

	recordClient := client.Resource(res).Namespace(config.Namespace)

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	if err := recordClient.Delete(ctx, config.ObjectName(), deleteOptions); err != nil {
		return fmt.Errorf("delete record: %w", err)
	}

	return nil
}
