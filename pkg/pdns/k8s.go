package pdns

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func buildClient() (dynamic.Interface, error) {
	p, err := getKubeConfigPath()
	if err != nil {
		return nil, fmt.Errorf("get kubeconfig path: %w", err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", p)
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes client: %w", err)
	}

	return client, nil
}

func getKubeConfigPath() (string, error) {
	if p := os.Getenv("KUBECONFIG"); p != "" {
		return p, nil
	}

	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config"), nil
	}

	return "", fmt.Errorf("unable to get kubeconfig path")
}
