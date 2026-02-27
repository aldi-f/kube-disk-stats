package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourusername/kube-disk-stats/internal/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNodeStatsSummary(ctx context.Context, nodeName string) (*models.StatsSummary, error) {
	req := c.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(nodeName).
		SubResource("proxy").
		Suffix("stats/summary")

	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats summary for node %s: %w", nodeName, err)
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to read stats summary for node %s: %w", nodeName, err)
	}

	var summary models.StatsSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse stats summary for node %s: %w", nodeName, err)
	}

	return &summary, nil
}

func (c *Client) GetNodeCapacity(ctx context.Context, nodeName string) (int64, error) {
	node, err := c.Clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	for _, addr := range node.Status.Addresses {
		if addr.Type == "InternalIP" {
			continue
		}
	}

	return 50 * 1024 * 1024 * 1024, nil
}
