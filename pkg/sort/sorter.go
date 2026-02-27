package sort

import (
	"sort"

	"github.com/yourusername/kube-disk-stats/internal/models"
)

type NodeSorter struct {
	Nodes []*models.NodeStorage
	Limit int
}

func (ns *NodeSorter) SortByUsedBytes() []*models.NodeStorage {
	sorted := make([]*models.NodeStorage, len(ns.Nodes))
	copy(sorted, ns.Nodes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].UsedBytes > sorted[j].UsedBytes
	})

	if ns.Limit > 0 && ns.Limit < len(sorted) {
		return sorted[:ns.Limit]
	}
	return sorted
}

func (ns *NodeSorter) SortByPercentage() []*models.NodeStorage {
	sorted := make([]*models.NodeStorage, len(ns.Nodes))
	copy(sorted, ns.Nodes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Percentage > sorted[j].Percentage
	})

	if ns.Limit > 0 && ns.Limit < len(sorted) {
		return sorted[:ns.Limit]
	}
	return sorted
}

type PodSorter struct {
	Pods  []*models.PodStorage
	Limit int
}

func (ps *PodSorter) SortByUsedBytes() []*models.PodStorage {
	sorted := make([]*models.PodStorage, len(ps.Pods))
	copy(sorted, ps.Pods)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].TotalBytes > sorted[j].TotalBytes
	})

	if ps.Limit > 0 && ps.Limit < len(sorted) {
		return sorted[:ps.Limit]
	}
	return sorted
}
