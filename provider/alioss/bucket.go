package alioss

import (
	"errors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/slipfre/imgmd/provider"
)

// Bucket Ali OSS Bucket 服务集成类
type Bucket struct {
	aliBucket *oss.Bucket
}

// NewBucket 创建一个 Bucket 对象
func NewBucket(bucket *oss.Bucket) *Bucket {
	return &Bucket{
		aliBucket: bucket,
	}
}

// PutObjectFromFile 上传本地文件
func (bucket *Bucket) PutObjectFromFile(objectKey, filePath string, options ...provider.ObjectOption) (err error) {
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
	bucket.aliBucket.PutObjectFromFile(
		objectKey,
		filePath,
		oss.ObjectACL(aliACL),
		oss.StorageClass(aliStorageClass),
		oss.RedundancyType(aliRedundancyType),
	)
	return
}

// ToAliACL 把 provider.ACL 转化为 oss.ACLType
func toAliACL(acl provider.ACL) (ossACL oss.ACLType, err error) {
	switch acl {
	case provider.Private:
		ossACL = oss.ACLPrivate
	case provider.PublicRead:
		ossACL = oss.ACLPublicRead
	case provider.PublicReadWrite:
		ossACL = oss.ACLPublicReadWrite
	default:
		err = errors.New("invalid acl type")
	}
	return
}

// ToAliStorageClass 把 provider.Storage 妆化为 oss.StorageClassType
func toAliStorageClass(storage provider.Storage) (ossStorageClass oss.StorageClassType, err error) {
	switch storage {
	case provider.Archive:
		ossStorageClass = oss.StorageArchive
	case provider.ColdArchive:
		ossStorageClass = oss.StorageColdArchive
	case provider.InfrequentAccess:
		ossStorageClass = oss.StorageIA
	case provider.Standard:
		ossStorageClass = oss.StorageStandard
	default:
		err = errors.New("invalid storage type")
	}
	return
}

// ToAliRedundancyType 把 provider.DataRedundancyType 转化为 oss.DataRedundancyType
func toAliRedundancyType(redundancyType provider.DataRedundancyType) (ossRedundancyType oss.DataRedundancyType, err error) {
	switch redundancyType {
	case provider.LRS:
		ossRedundancyType = oss.RedundancyLRS
	case provider.ZRS:
		ossRedundancyType = oss.RedundancyZRS
	default:
		err = errors.New("invalid redundancy type")
	}
	return
}
