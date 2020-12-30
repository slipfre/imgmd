package conf

import (
	"errors"

	"github.com/slipfre/imgmd/provider"
	"github.com/spf13/viper"
)

// GetBucketFromConfigFile 解析配置文件
func GetBucketFromConfigFile(path string) (bucket provider.Bucket, err error) {
	viper.SetConfigFile(path)
	viper.SetConfigFile("yaml")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	obs := viper.GetStringMapString("OBS")
	if obs == nil {
		err = errors.New("'OBS' not found in config file")
		return
	}
	client, err := provider.ObtainClient(provider.Provider(obs["PROVIDER"]), obs["AKID"], obs["AKS"], obs["ENDPOINT"])
	if err != nil {
		return
	}
	bucket, err = client.GetOrCreateBucket(obs["BUCKET"])
	if err != nil {
		return
	}
	repo := viper.GetStringMapString("REPOSITORY")
	if repo == nil {
		err = errors.New("'REPOSITORY' not found in config file")
		return
	}
	return
}

// GetRepoPathFromConfig 获取 Repository 的路径
func GetRepoPathFromConfig(path string) (repoPath string, err error) {
	viper.SetConfigFile(path)
	viper.SetConfigFile("yaml")
	repo := viper.GetStringMapString("REPOSITORY")
	if repo == nil {
		err = errors.New("'OBS' not found in config file")
		return
	}
	repoPath = repo["PATH"]
	if repoPath == "" {
		err = errors.New("'PATH' of repository not config or is empty")
		return
	}
	return
}
