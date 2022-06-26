package pdns

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/retry"
)

// SetConfig is configuration for the set command.
type SetConfig struct {
	Zone        string
	Name        string
	Type        string
	RawContents string
	Namespace   string
}

// Validate validates the configuration.
func (c *SetConfig) Validate() error {
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

	if c.RawContents == "" {
		err = multierr.Append(err, fmt.Errorf("contents are blank"))
	}

	if c.Namespace == "" {
		c.Namespace = "default"
	}

	return err
}

// ToSpec converts the configurations to a record spec.
func (c *SetConfig) ToSpec() map[string]interface{} {
	m := map[string]interface{}{
		"zone":  c.Zone,
		"name":  c.Name,
		"type":  c.Type,
		"value": strings.Split(c.RawContents, ","),
	}

	return m
}

// ObjectName returns the name for the record based on the configuration.
func (c *SetConfig) ObjectName() string {
	return fmt.Sprintf("%s-%s", c.Name, strings.ToLower(c.Type))
}

// Set creates or updates a record.
func Set(config SetConfig) error {
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
	_, err = recordClient.Get(ctx, config.ObjectName(), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("check for existing record: %w", err)
		}

		if err := createRecord(ctx, recordClient, config); err != nil {
			return fmt.Errorf("create record: %w", err)
		}

		return nil
	}

	if err := updateRecord(ctx, recordClient, config); err != nil {
		return fmt.Errorf("update record: %w", err)
	}

	return nil
}

func createRecord(ctx context.Context, client dynamic.ResourceInterface, config SetConfig) error {
	record := initRecord(config)
	if _, err := client.Create(ctx, record, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create record on cluster: %w", err)
	}

	return nil
}

func updateRecord(
	ctx context.Context,
	client dynamic.ResourceInterface,
	config SetConfig) error {
	retryErr := retry.RetryOnConflict(
		retry.DefaultRetry, func() error {
			result, getErr := client.Get(ctx, config.ObjectName(), metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}

			if err := unstructured.SetNestedMap(result.Object, config.ToSpec(), "spec"); err != nil {
				return err
			}

			if _, err := client.Update(ctx, result, metav1.UpdateOptions{}); err != nil {
				return err
			}

			return nil
		})

	return retryErr
}

func initRecord(config SetConfig) *unstructured.Unstructured {
	record := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "pdns.bryanl.dev/v1alpha1",
			"kind":       "Record",
			"spec":       config.ToSpec(),
		},
	}

	record.SetName(config.ObjectName())
	record.SetNamespace(config.Namespace)

	return record
}
