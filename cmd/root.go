package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/aldi-f/kube-disk-stats/internal/analyzer"
	"github.com/aldi-f/kube-disk-stats/internal/display"
	"github.com/aldi-f/kube-disk-stats/internal/k8s"
	"github.com/aldi-f/kube-disk-stats/internal/models"
	kubesort "github.com/aldi-f/kube-disk-stats/pkg/sort"
)

var Version = "1.0.0"

var (
	contextFlag   string
	nodeFlag      string
	outputFlag    string
	topFlag       int
	watchFlag     bool
	intervalFlag  time.Duration
	breakdownFlag bool
)

type NodeStatsFetcher interface {
	GetNodeStatsSummary(ctx context.Context, nodeName string) (*models.StatsSummary, error)
}

type NodeLister interface {
	ListNodes(ctx context.Context) ([]string, error)
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-disk-stats",
		Short: "Query Kubernetes node and pod disk usage",
		Long:  `kube-disk-stats is a CLI tool for querying Kubernetes node and pod disk usage statistics.`,
	}

	cmd.PersistentFlags().StringVarP(&contextFlag, "context", "c", "", "Kubernetes context to use")
	cmd.PersistentFlags().StringVarP(&nodeFlag, "node", "n", "", "Query specific node (default: all nodes)")

	cmd.AddCommand(newAllCmd())
	cmd.AddCommand(newPodsCmd())
	cmd.AddCommand(newNodesCmd())
	cmd.AddCommand(newContainersCmd())
	cmd.AddCommand(newImagesCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "table", "Output format: table or json")
	cmd.Flags().IntVarP(&topFlag, "top", "t", 0, "Show top N results (0 = all)")
	cmd.Flags().BoolVar(&breakdownFlag, "breakdown", false, "Show rootfs/logs/images breakdown")
	cmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch mode: continuously refresh")
	cmd.Flags().DurationVarP(&intervalFlag, "interval", "i", 5*time.Second, "Refresh interval for watch mode")
}

func newAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Display all storage usage (nodes and pods)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), true, true, false)
		},
	}
	addOutputFlags(cmd)
	return cmd
}

func newPodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pods",
		Short: "Display pod storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), true, false, false)
		},
	}
	addOutputFlags(cmd)
	return cmd
}

func newNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Display node storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), false, true, false)
		},
	}
	addOutputFlags(cmd)
	return cmd
}

func newContainersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "containers",
		Short: "Display container storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), false, false, true)
		},
	}
	addOutputFlags(cmd)
	return cmd
}

func newImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Display Docker images on nodes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if watchFlag && breakdownFlag {
				return fmt.Errorf("--watch and --breakdown are mutually exclusive")
			}
			return nil
		},
		RunE: runImages,
	}
	addOutputFlags(cmd)
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kube-disk-stats version %s\n", Version)
		},
	}
}

func run(ctx context.Context, showPods, showNodes, showContainers bool) error {
	client, err := k8s.NewClient(contextFlag)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	if watchFlag {
		return watchLoop(ctx, client, showPods, showNodes, showContainers)
	}

	return executeQuery(ctx, client, showPods, showNodes, showContainers)
}

func runImages(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient(contextFlag)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	if watchFlag {
		return watchImagesLoop(context.Background(), client)
	}

	return executeImagesQuery(context.Background(), client)
}

func watchImagesLoop(ctx context.Context, client *k8s.Client) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(intervalFlag)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Print("\033[2J\033[H")
			fmt.Printf("Last updated: %s\n\n", time.Now().Format(time.RFC3339))

			if err := executeImagesQuery(ctx, client); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}

		case <-sigChan:
			fmt.Println("\nStopping watch mode...")
			return nil
		}
	}
}

func executeImagesQuery(ctx context.Context, client *k8s.Client) error {
	var nodeNames []string
	var err error

	if nodeFlag != "" {
		nodeNames = []string{nodeFlag}
	} else {
		nodeNames, err = client.ListNodes(ctx)
		if err != nil {
			return fmt.Errorf("failed to list nodes: %w", err)
		}
	}

	if len(nodeNames) > 1 {
		fmt.Printf("Found %d nodes, querying...\n", len(nodeNames))
	}

	nodeImages := make(map[string][]models.NodeImage)
	summaries := make([]models.NodeImageSummary, 0, len(nodeNames))

	for i, nodeName := range nodeNames {
		if len(nodeNames) > 1 {
			fmt.Printf("\rQuerying node %d/%d: %s", i+1, len(nodeNames), nodeName)
		}

		images, err := client.GetNodeImages(ctx, nodeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get images for node %s: %v\n", nodeName, err)
			continue
		}

		nodeImages[nodeName] = images

		var totalSize int64
		for _, img := range images {
			totalSize += img.SizeBytes
		}

		summaries = append(summaries, models.NodeImageSummary{
			NodeName:   nodeName,
			ImageCount: len(images),
			TotalSize:  totalSize,
		})
	}

	if len(nodeNames) > 1 {
		fmt.Println()
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalSize > summaries[j].TotalSize
	})

	if topFlag > 0 && topFlag < len(summaries) {
		summaries = summaries[:topFlag]
	}

	if outputFlag == "json" {
		return displayImagesJSON(summaries, nodeImages)
	}

	display.DisplayImagesSummaryTable(summaries)

	if breakdownFlag {
		fmt.Println()
		for _, summary := range summaries {
			images := nodeImages[summary.NodeName]
			if len(images) > 0 {
				fmt.Printf("\nNode: %s\n", summary.NodeName)
				display.DisplayImagesTable(summary.NodeName, images, false)
			}
		}
	}

	return nil
}

func watchLoop(ctx context.Context, client *k8s.Client, showPods, showNodes, showContainers bool) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(intervalFlag)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Print("\033[2J\033[H")
			fmt.Printf("Last updated: %s\n\n", time.Now().Format(time.RFC3339))

			if err := executeQuery(ctx, client, showPods, showNodes, showContainers); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}

		case <-sigChan:
			fmt.Println("\nStopping watch mode...")
			return nil
		}
	}
}

func executeQuery(ctx context.Context, client *k8s.Client, showPods, showNodes, showContainers bool) error {
	var nodeNames []string
	var err error

	if nodeFlag != "" {
		nodeNames = []string{nodeFlag}
	} else {
		nodeNames, err = client.ListNodes(ctx)
		if err != nil {
			return fmt.Errorf("failed to list nodes: %w", err)
		}
	}

	if !watchFlag && len(nodeNames) > 1 {
		fmt.Printf("Found %d nodes, querying...\n", len(nodeNames))
	}

	nodes := make([]*models.NodeStorage, 0, len(nodeNames))
	allContainers := make([]models.Container, 0)

	for i, nodeName := range nodeNames {
		if !watchFlag && len(nodeNames) > 1 {
			fmt.Printf("\rQuerying node %d/%d: %s", i+1, len(nodeNames), nodeName)
		}

		summary, err := client.GetNodeStatsSummary(ctx, nodeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get stats for node %s: %v\n", nodeName, err)
			continue
		}

		images, err := client.GetNodeImages(ctx, nodeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get images for node %s: %v\n", nodeName, err)
			images = nil
		}

		var imageBytes int64
		for _, img := range images {
			imageBytes += img.SizeBytes
		}

		totalBytes := int64(50 * 1024 * 1024 * 1024)
		nodeStorage := analyzer.CalculateNodeStorage(summary, nodeName, totalBytes, imageBytes)
		nodes = append(nodes, nodeStorage)
		allContainers = append(allContainers, nodeStorage.Containers...)
	}

	if !watchFlag && len(nodeNames) > 1 {
		fmt.Println()
	}

	if outputFlag == "json" {
		return displayJSON(nodes, allContainers, showPods, showNodes, showContainers)
	}

	return displayTable(nodes, allContainers, showPods, showNodes, showContainers)
}

func displayTable(nodes []*models.NodeStorage, containers []models.Container, showPods, showNodes, showContainers bool) error {
	if showNodes {
		sorter := kubesort.NodeSorter{Nodes: nodes, Limit: topFlag}
		sortedNodes := sorter.SortByPercentage()
		display.DisplayNodesTable(sortedNodes, breakdownFlag)
		fmt.Println()
	}

	if showPods {
		nodeCapacities := make(map[string]int64)
		for _, node := range nodes {
			nodeCapacities[node.Name] = node.TotalBytes
		}

		pods := analyzer.GroupByPod(containers, nodeCapacities)
		sorter := kubesort.PodSorter{Pods: pods, Limit: topFlag}
		sortedPods := sorter.SortByUsedBytes()

		display.DisplayPodsTable(sortedPods, breakdownFlag)
		fmt.Println()
	}

	if showContainers {
		sortedContainers := containers
		if topFlag > 0 && topFlag < len(sortedContainers) {
			sortedContainers = sortedContainers[:topFlag]
		}
		display.DisplayContainersTable(sortedContainers)
	}

	return nil
}

func displayJSON(nodes []*models.NodeStorage, containers []models.Container, showPods, showNodes, showContainers bool) error {
	if showNodes {
		if err := display.DisplayNodesJSON(nodes); err != nil {
			return err
		}
	}

	if showPods {
		nodeCapacities := make(map[string]int64)
		for _, node := range nodes {
			nodeCapacities[node.Name] = node.TotalBytes
		}

		pods := analyzer.GroupByPod(containers, nodeCapacities)
		if err := display.DisplayPodsJSON(pods); err != nil {
			return err
		}
	}

	if showContainers {
		if err := display.DisplayContainersJSON(containers); err != nil {
			return err
		}
	}

	return nil
}

func displayImagesJSON(summaries []models.NodeImageSummary, nodeImages map[string][]models.NodeImage) error {
	if breakdownFlag {
		return display.DisplayImagesDetailJSON(summaries, nodeImages)
	}
	return display.DisplayImagesSummaryJSON(summaries)
}
