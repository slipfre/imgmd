package provider

// Client OBS 服务供应商的客户端接口，用于提供 bucket 的增删改查操作
type Client interface {
	IsBucketExist(bucketName string) (isExist bool, err error)
	CreateBucket(bucketName string, options ...ObjectOption) (err error)
	GetOrCreateBucket(bucketName string, options ...ObjectOption) (bucket Bucket, err error)
	DeleteBucket(bucketName string) (err error)
}
