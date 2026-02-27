package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/kube-disk-stats/internal/analyzer"
	"github.com/yourusername/kube-disk-stats/internal/display"
	"github.com/yourusername/kube-disk-stats/internal/k8s"
	"github.com/yourusername/kube-disk-stats/internal/models"
	kubesort "github.com/yourusername/kube-disk-stats/pkg/sort"
)

var Version = "dev"

var (
	contextFlag  string
	nodeFlag     string
	outputFlag   string
	topFlag      int
	watchFlag    bool
	intervalFlag time.Duration
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
		RunE:  runStats,
	}

	cmd.Flags().StringVarP(&contextFlag, "context", "c", "", "Kubernetes context to use")
	cmd.Flags().StringVarP(&nodeFlag, "node", "n", "", "Query specific node (default: all nodes)")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "table", "Output format: table or json")
	cmd.Flags().IntVarP(&topFlag, "top", "t", 0, "Show top N results (0 = all)")
	cmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch mode: continuously refresh")
	cmd.Flags().DurationVarP(&intervalFlag, "interval", "i", 5*time.Second, "Refresh interval for watch mode")

	cmd.AddCommand(newPodsCmd())
	cmd.AddCommand(newNodesCmd())
	cmd.AddCommand(newContainersCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

func runStats(cmd *cobra.Command, args []string) error {
	return run(context.Background(), true, true, false)
}

func newPodsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pods",
		Short: "Display pod storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), true, false, false)
		},
	}
}

func newNodesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "nodes",
		Short: "Display node storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), false, true, false)
		},
	}
}

func newContainersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "containers",
		Short: "Display container storage usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), false, false, true)
		},
	}
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

	nodes := make([]*models.NodeStorage, 0, len(nodeNames))
	allContainers := make([]models.Container, 0)

	for _, nodeName := range nodeNames {
		summary, err := client.GetNodeStatsSummary(ctx, nodeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get stats for node %s: %v\n", nodeName, err)
			continue
		}

		totalBytes := int64(50 * 1024 * 1024 * 1024)
		nodeStorage := analyzer.CalculateNodeStorage(summary, nodeName, totalBytes)
		nodes = append(nodes, nodeStorage)
		allContainers = append(allContainers, nodeStorage.Containers...)
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
		display.DisplayNodesTable(sortedNodes)
		fmt.Println()
	}

	if showPods {
		pods := analyzer.GroupByPod(containers)
		sorter := kubesort.PodSorter{Pods: pods, Limit: topFlag}
		sortedPods := sorter.SortByUsedBytes()

		var totalNodeBytes int64
		for _, node := range nodes {
			totalNodeBytes += node.TotalBytes
		}

		display.DisplayPodsTable(sortedPods, totalNodeBytes)
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
		pods := analyzer.GroupByPod(containers)
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
