package provider

import (
	"fmt"

	"github.com/slipfre/imgmd/provider/alioss"
)

// Provider Represent name of OBS Service
type Provider string

const (
	// ALI Represent ali obs
	ALI Provider = "ali"
)

// ObtainClient 获取 Client
func ObtainClient(provider Provider, akid, aks, endpoint string) (Client, error) {
	switch provider {
	case ALI:
		return alioss.NewClient(endpoint, akid, aks), nil
	}
	return nil, fmt.Errorf("unsupported provider: '%s'", provider)
}
