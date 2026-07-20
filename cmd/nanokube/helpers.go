package main

import (
	"fmt"
	"os"
	"strings"
)

// defaultNodeName matches kubeadm/kubelet: lowercased OS hostname.
func defaultNodeName() (string, error) {
	h, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("get hostname: %w", err)
	}
	return strings.ToLower(h), nil
}
