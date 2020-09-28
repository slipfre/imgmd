package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/slipfre/imgmd/provider"
)

// Client Ali OSS Client 服务集成类
type Client struct {
	aliOSSClient *oss.Client
}

// NewClient 创建 Client 对象
func NewClient(endpoint, accessKeyID, accessKeySecret string) (client *Client, err error) {
	aliOssClient, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return
	}
	client = &Client{
		aliOSSClient: aliOssClient,
	}
	return
}

// IsBucketExist 判断 Bucket 是否存在
func (client *Client) IsBucketExist(bucketName string) (isExist bool, err error) {
	isExist, err = client.aliOSSClient.IsBucketExist(bucketName)
	return
}

// CreateBucket 创建 Bucket
func (client *Client) CreateBucket(bucketName string, options ...provider.ObjectOption) (err error) {
	objectConfig := provider.DefaultOptionConfig()
	for _, option := range options {
		option(objectConfig)
	}
	aliACL, err := toAliACL(objectConfig.ACL)
	if err != nil {
		return
	}
	aliStorageClass, err := toAliStorageClass(objectConfig.Storage)
	if err != nil {
		return
	}
	aliRedundancyType, err := toAliRedundancyType(objectConfig.RedundancyType)
	if err != nil {
		return
	}
	err = client.aliOSSClient.CreateBucket(
		bucketName,
		oss.ObjectACL(aliACL),
		oss.StorageClass(aliStorageClass),
		oss.RedundancyType(aliRedundancyType),
	)
	return
}

// GetOrCreateBucket 获取 Bucket，如果不存在则创建该 Bucket
func (client *Client) GetOrCreateBucket(bucketName string, options ...provider.ObjectOption) (providerBucket provider.Bucket, err error) {
	isExist, err := client.IsBucketExist(bucketName)
	if err != nil {
		return
	}
	if !isExist {
		err = client.CreateBucket(bucketName, options...)
	}
	aliBucket, err := client.aliOSSClient.Bucket(bucketName)
	if err != nil {
		return
	}
	providerBucket = NewBucket(aliBucket)
	return
}

// DeleteBucket 删除 Bucket, 只有 Bucket 为 empty 时可以删除
func (client *Client) DeleteBucket(bucketName string) (err error) {
	err = client.aliOSSClient.DeleteBucket(bucketName)
	return
}
