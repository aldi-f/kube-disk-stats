package display

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/kube-disk-stats/internal/models"
)

func DisplayNodesJSON(nodes []*models.NodeStorage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(nodes); err != nil {
		return fmt.Errorf("failed to encode nodes: %w", err)
	}
	return nil
}

func DisplayPodsJSON(pods []*models.PodStorage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(pods); err != nil {
		return fmt.Errorf("failed to encode pods: %w", err)
	}
	return nil
}

func DisplayContainersJSON(containers []models.Container) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(containers); err != nil {
		return fmt.Errorf("failed to encode containers: %w", err)
	}
	return nil
}
