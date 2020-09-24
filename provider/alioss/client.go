package alioss

// Client Ali OSS Client 服务集成类
type Client struct {
}

// IsBucketExist 判断 Bucket 是否存在
func (client *Client) IsBucketExist(bucketName string) (isExist bool, err error) {
	return
}

// CreateBucket 创建 Bucket
func (client *Client) CreateBucket(bucketName string, options ...ObjectOption) (err error) {
	return
}

// GetOrCreateBucket 获取 Bucket，如果不存在则创建该 Bucket
func (client *Client) GetOrCreateBucket(bucketName string, options ...ObjectOption) (bucket *Bucket, err error) {
	return
}
