package display

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/aldi-f/kube-disk-stats/internal/models"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
	KB = 1024
)

func formatBytes(bytes int64) string {
	if bytes >= GB {
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GB))
	}
	if bytes >= MB {
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MB))
	}
	if bytes >= KB {
		return fmt.Sprintf("%.2f KiB",
			float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}

func DisplayNodesTable(nodes []*models.NodeStorage, breakdown bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"NODE", "TOTAL", "USED", "%", "PODS", "AGE"}
	if breakdown {
		headers = []string{"NODE", "TOTAL", "ROOTFS", "LOGS", "USED", "%", "PODS", "AGE"}
	}

	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, node := range nodes {
		percentage := ColorizePercentageText(node.Percentage)

		row := []string{
			truncateString(node.Name, 35),
			formatBytes(node.TotalBytes),
		}

		if breakdown {
			rootfsBytes := int64(0)
			logsBytes := int64(0)
			for _, container := range node.Containers {
				rootfsBytes += container.RootFSBytes
				logsBytes += container.LogsBytes
			}
			row = append(row, formatBytes(rootfsBytes), formatBytes(logsBytes))
		}

		row = append(row,
			percentage,
			fmt.Sprintf("%d", node.PodCount),
			node.Age,
		)

		table.Append(row)
	}

	table.Render()
}

func DisplayPodsTable(pods []*models.PodStorage, breakdown bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"NAMESPACE", "POD", "NODE", "USAGE", "%", "AGE"}
	if breakdown {
		headers = []string{"NAMESPACE", "POD", "NODE", "ROOTFS", "LOGS", "USED", "%", "AGE"}
	}

	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, pod := range pods {
		var percentage float64
		if pod.NodeTotalBytes > 0 {
			percentage = float64(pod.TotalBytes) / float64(pod.NodeTotalBytes) * 100
		}
		percentageText := ColorizePercentageText(percentage)

		row := []string{
			truncateString(pod.Namespace, 20),
			truncateString(pod.Name, 30),
			truncateString(pod.NodeName, 30),
		}

		if breakdown {
			rootfsBytes := int64(0)
			logsBytes := int64(0)
			for _, container := range pod.Containers {
				rootfsBytes += container.RootFSBytes
				logsBytes += container.LogsBytes
			}
			row = append(row, formatBytes(rootfsBytes), formatBytes(logsBytes), formatBytes(pod.TotalBytes))
		} else {
			row = append(row, formatBytes(pod.TotalBytes))
		}

		row = append(row, percentageText, pod.Age)

		table.Append(row)
	}

	table.Render()
}

func DisplayContainersTable(containers []models.Container) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAMESPACE", "POD", "CONTAINER", "ROOTFS", "LOGS", "TOTAL"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, container := range containers {
		table.Append([]string{
			truncateString(container.Namespace, 20),
			truncateString(container.PodName, 30),
			truncateString(container.Name, 25),
			formatBytes(container.RootFSBytes),
			formatBytes(container.LogsBytes),
			formatBytes(container.TotalBytes),
		})
	}

	table.Render()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func DisplayImagesTable(nodeName string, images []models.NodeImage) {
	var totalBytes int64
	for _, img := range images {
		totalBytes += img.SizeBytes
	}

	fmt.Printf("Node: %s\n", nodeName)
	fmt.Printf("Total images: %d | Total size: %s\n\n", len(images), formatBytes(totalBytes))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IMAGE", "SIZE", "%"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, img := range images {
		primaryName := ""
		if len(img.Names) > 0 {
			primaryName = img.Names[0]
		}

		var percentage float64
		if totalBytes > 0 {
			percentage = float64(img.SizeBytes) / float64(totalBytes) * 100
		}

		table.Append([]string{
			primaryName,
			formatBytes(img.SizeBytes),
			fmt.Sprintf("%.1f%%", percentage),
		})
	}

	table.Render()
	fmt.Println()
}
