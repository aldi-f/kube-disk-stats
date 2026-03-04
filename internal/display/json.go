package display

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aldi-f/kube-disk-stats/internal/models"
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

func DisplayImagesJSON(nodeName string, images []models.NodeImage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	output := struct {
		Node   string             `json:"node"`
		Images []models.NodeImage `json:"images"`
	}{
		Node:   nodeName,
		Images: images,
	}

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode images: %w", err)
	}
	return nil
}

func DisplayImagesSummaryJSON(summaries []models.NodeImageSummary) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(summaries); err != nil {
		return fmt.Errorf("failed to encode image summaries: %w", err)
	}
	return nil
}

func DisplayImagesDetailJSON(summaries []models.NodeImageSummary, nodeImages map[string][]models.NodeImage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	output := struct {
		Summaries []models.NodeImageSummary     `json:"summaries"`
		Details   map[string][]models.NodeImage `json:"details"`
	}{
		Summaries: summaries,
		Details:   nodeImages,
	}

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode image details: %w", err)
	}
	return nil
}
