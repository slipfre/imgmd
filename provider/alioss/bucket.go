package alioss

// Bucket Ali OSS Bucket 服务集成类
type Bucket struct {
}

// PutObjectFromFile 上传本地文件
func (bucket *Bucket) PutObjectFromFile(objectKey, filePath string, options ...ObjectOption) (err error) {
	return
}
