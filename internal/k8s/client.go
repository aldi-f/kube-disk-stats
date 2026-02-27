package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aldi-f/kube-disk-stats/internal/models"
)

type Client struct {
	*kubernetes.Clientset
}

func NewClient(context string) (*Client, error) {
	var config *rest.Config
	var err error

	if context != "" {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{
				CurrentContext: context,
			},
		).ClientConfig()
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				clientcmd.NewDefaultClientConfigLoadingRules(),
				&clientcmd.ConfigOverrides{},
			).ClientConfig()
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{clientset}, nil
}

func (c *Client) ListNodes(ctx context.Context) ([]string, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	nodeNames := make([]string, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	return nodeNames, nil
}

func (c *Client) GetContextName() (string, error) {
	rawConfig, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}
	return rawConfig.CurrentContext, nil
}

func (c *Client) GetNodeImages(ctx context.Context, nodeName string) ([]models.NodeImage, error) {
	node, err := c.Clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	images := make([]models.NodeImage, 0, len(node.Status.Images))
	for _, img := range node.Status.Images {
		images = append(images, models.NodeImage{
			Names:     img.Names,
			SizeBytes: img.SizeBytes,
		})
	}

	return images, nil
}
