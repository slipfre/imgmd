package alioss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

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
func (bucket *Bucket) PutObjectFromFile(objectKey, filePath string, options ...provider.ObjectOption) (url string, err error) {
	if reader, err := os.Open(filePath); err == nil {
		defer reader.Close()
		return bucket.PutObject(objectKey, reader, options...)
	}
	return
}

// PutObjectFromBytes 上传 byte 数组
func (bucket *Bucket) PutObjectFromBytes(objectKey string, data []byte, options ...provider.ObjectOption) (url string, err error) {
	return bucket.PutObject(objectKey, bytes.NewReader(data), options...)
}

// PutObject 上传 Object
func (bucket *Bucket) PutObject(objectKey string, reader io.Reader, options ...provider.ObjectOption) (url string, err error) {
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
	err = bucket.aliBucket.PutObject(
		objectKey,
		reader,
		oss.ObjectACL(aliACL),
		oss.StorageClass(aliStorageClass),
		oss.RedundancyType(aliRedundancyType),
	)
	if err != nil {
		return
	}
	url = fmt.Sprintf(
		"http://%s.%s/%s",
		bucket.aliBucket.BucketName,
		bucket.aliBucket.GetConfig().Endpoint,
		objectKey,
	)
	return
}

// DeleteObject 删除 Object
func (bucket *Bucket) DeleteObject(objectKey string) (err error) {
	err = bucket.aliBucket.DeleteObject(objectKey)
	return
}

// GetObjectLastModified 获取 Object 最后一次修改的时间
func (bucket *Bucket) GetObjectLastModified(objectKey string) (*time.Time, error) {
	headers, err := bucket.aliBucket.GetObjectMeta(objectKey)
	if err != nil {
		return nil, err
	}
	lastModified := headers.Get("Last-Modified")
	lastModifiedTime, err := time.ParseInLocation(time.RFC1123, lastModified, time.UTC)
	if err != nil {
		return nil, err
	}
	return &lastModifiedTime, err
}

// IsObjectExist 判断 Object 是否存在
func (bucket *Bucket) IsObjectExist(objectKey string) (isExist bool, err error) {
	isExist, err = bucket.aliBucket.IsObjectExist(objectKey)
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

// ToAliStorageClass 把 provider.Storage 转化为 oss.StorageClassType
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
