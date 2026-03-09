package analyzer

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/aldi-f/kube-disk-stats/internal/models"
)

func CalculateNodeStorage(summary *models.StatsSummary, nodeName string, totalBytes int64, imageBytes int64, images []models.NodeImage, pods []*corev1.Pod, includeImage bool) *models.NodeStorage {
	node := &models.NodeStorage{
		Name:       nodeName,
		Age:        calculateAge(summary.Node.StartTime),
		TotalBytes: totalBytes,
		ImageBytes: imageBytes,
	}

	imageMap := make(map[string]int64)
	for _, img := range images {
		for _, name := range img.Names {
			imageMap[name] = img.SizeBytes
		}
	}

	containerImageMap := make(map[string]string)
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			key := fmt.Sprintf("%s/%s/%s", pod.Namespace, pod.Name, container.Name)
			containerImageMap[key] = container.Image
		}
	}

	var totalUsed int64
	containers := make([]models.Container, 0)

	for _, pod := range summary.Pods {
		podTotal := int64(0)

		for _, container := range pod.Containers {
			rootfs := getSafeInt(container.RootFS.UsedBytes)
			logs := getSafeInt(container.Logs.UsedBytes)

			var containerImageBytes int64
			if includeImage {
				key := fmt.Sprintf("%s/%s/%s", pod.PodRef.Namespace, pod.PodRef.Name, container.Name)
				if imageName, exists := containerImageMap[key]; exists {
					containerImageBytes = imageMap[imageName]
				}
			}

			totalBytesContainer := rootfs + logs
			if includeImage {
				totalBytesContainer += containerImageBytes
			}

			containers = append(containers, models.Container{
				Name:        container.Name,
				PodName:     pod.PodRef.Name,
				Namespace:   pod.PodRef.Namespace,
				PodAge:      calculateAge(pod.StartTime),
				RootFSBytes: rootfs,
				LogsBytes:   logs,
				ImageBytes:  containerImageBytes,
				TotalBytes:  totalBytesContainer,
				NodeName:    nodeName,
				NodeAge:     node.Age,
			})

			podTotal += totalBytesContainer
		}

		totalUsed += podTotal
	}

	if !includeImage {
		totalUsed += imageBytes
	}
	node.UsedBytes = totalUsed
	if totalBytes > 0 {
		node.Percentage = float64(totalUsed) / float64(totalBytes) * 100
	}
	node.Containers = containers
	node.PodCount = len(summary.Pods)

	return node
}

func getSafeInt(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func calculateAge(startTime string) string {
	if startTime == "" {
		return "unknown"
	}

	t, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return "unknown"
	}

	age := time.Since(t)

	days := int(age.Hours() / 24)
	hours := int(age.Hours()) % 24
	minutes := int(age.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func GroupByPod(containers []models.Container, nodeCapacities map[string]int64) []*models.PodStorage {
	podMap := make(map[string]*models.PodStorage)

	for _, container := range containers {
		key := container.Namespace + "/" + container.PodName
		if pod, exists := podMap[key]; exists {
			pod.Containers = append(pod.Containers, container)
			pod.TotalBytes += container.TotalBytes
			pod.ImageBytes += container.ImageBytes
		} else {
			nodeTotalBytes := int64(0)
			if cap, ok := nodeCapacities[container.NodeName]; ok {
				nodeTotalBytes = cap
			}
			podMap[key] = &models.PodStorage{
				Name:           container.PodName,
				Namespace:      container.Namespace,
				NodeName:       container.NodeName,
				Age:            container.PodAge,
				TotalBytes:     container.TotalBytes,
				NodeTotalBytes: nodeTotalBytes,
				ImageBytes:     container.ImageBytes,
				Containers:     []models.Container{container},
			}
		}
	}

	result := make([]*models.PodStorage, 0, len(podMap))
	for _, pod := range podMap {
		result = append(result, pod)
	}

	return result
}
