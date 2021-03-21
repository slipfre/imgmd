package conf

import (
	"errors"

	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/provider/factory"
	"github.com/spf13/viper"
)

// GetBucketFromConfigFile 解析配置文件
func GetBucketFromConfigFile(path string) (bucket provider.Bucket, err error) {
	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	obs := viper.GetStringMapString("OBS")
	if obs == nil {
		err = errors.New("'OBS' not found in config file")
		return
	}
	client, err := factory.ObtainClient(factory.Provider(obs["provider"]), obs["akid"], obs["aks"], obs["endpoint"])
	if err != nil {
		return
	}
	bucket, err = client.GetOrCreateBucket(obs["bucket"])
	if err != nil {
		return
	}
	repo := viper.GetStringMapString("repository")
	if repo == nil {
		err = errors.New("'REPOSITORY' not found in config file")
		return
	}
	return
}

// GetRepoPathFromConfig 获取 Repository 的路径
func GetRepoPathFromConfig(path string) (repoPath string, err error) {
	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	repo := viper.GetStringMapString("REPOSITORY")
	if repo == nil {
		err = errors.New("'REPOSITORY' not found in config file")
		return
	}
	repoPath = repo["path"]
	if repoPath == "" {
		err = errors.New("'PATH' of repository not config or is empty")
		return
	}
	return
}
