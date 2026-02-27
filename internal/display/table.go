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
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}

func DisplayNodesTable(nodes []*models.NodeStorage) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NODE", "TOTAL", "USED", "%", "PODS", "AGE"})
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
		table.Append([]string{
			truncateString(node.Name, 35),
			formatBytes(node.TotalBytes),
			formatBytes(node.UsedBytes),
			percentage,
			fmt.Sprintf("%d", node.PodCount),
			node.Age,
		})
	}

	table.Render()
}

func DisplayPodsTable(pods []*models.PodStorage, nodeTotalBytes int64) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAMESPACE", "POD", "NODE", "USAGE", "%", "AGE"})
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
		if nodeTotalBytes > 0 {
			percentage = float64(pod.TotalBytes) / float64(nodeTotalBytes) * 100
		}
		percentageText := ColorizePercentageText(percentage)

		table.Append([]string{
			truncateString(pod.Namespace, 20),
			truncateString(pod.Name, 30),
			truncateString(pod.NodeName, 30),
			formatBytes(pod.TotalBytes),
			percentageText,
			pod.Age,
		})
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
