package models

import "time"

type NodeStorage struct {
	Name       string      `json:"name"`
	Age        string      `json:"age"`
	TotalBytes int64       `json:"total_bytes"`
	UsedBytes  int64       `json:"used_bytes"`
	ImageBytes int64       `json:"image_bytes"`
	Percentage float64     `json:"percentage"`
	PodCount   int         `json:"pod_count"`
	Containers []Container `json:"containers"`
}

type Container struct {
	Name        string `json:"name"`
	PodName     string `json:"pod_name"`
	Namespace   string `json:"namespace"`
	PodAge      string `json:"pod_age"`
	RootFSBytes int64  `json:"rootfs_bytes"`
	LogsBytes   int64  `json:"logs_bytes"`
	TotalBytes  int64  `json:"total_bytes"`
	NodeName    string `json:"node_name"`
	NodeAge     string `json:"node_age"`
}

type PodStorage struct {
	Name           string      `json:"name"`
	Namespace      string      `json:"namespace"`
	NodeName       string      `json:"node_name"`
	Age            string      `json:"age"`
	TotalBytes     int64       `json:"total_bytes"`
	NodeTotalBytes int64       `json:"node_total_bytes"`
	Containers     []Container `json:"containers"`
}

type NodeImage struct {
	Names     []string `json:"names"`
	SizeBytes int64    `json:"size_bytes"`
}

type NodeImageSummary struct {
	NodeName   string `json:"node_name"`
	ImageCount int    `json:"image_count"`
	TotalSize  int64  `json:"total_size"`
}

type Config struct {
	Context    string
	NodeFilter string
	Output     string
	Top        int
	Watch      bool
	Interval   time.Duration
	Breakdown  bool
	ShowImages bool
}

type StatsSummary struct {
	Node NodeStats  `json:"node"`
	Pods []PodStats `json:"pods"`
}

type NodeStats struct {
	StartTime string `json:"startTime"`
}

type PodStats struct {
	PodRef     PodReference     `json:"podRef"`
	StartTime  string           `json:"startTime"`
	Containers []ContainerStats `json:"containers"`
}

type PodReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ContainerStats struct {
	Name   string      `json:"name"`
	RootFS StorageInfo `json:"rootfs"`
	Logs   StorageInfo `json:"logs"`
}

type StorageInfo struct {
	UsedBytes *int64 `json:"usedBytes"`
}
