package provider

// Provider OBS 服务供应商接口
type Provider interface {
	IsBucketExist(bucketName string) (isExist bool, err error)
	CreateBucket(bucketName string, acl ACL) (err error)
	PutObjectFromFile() (err error)
}

// ACL Bucket 的读写访问权限
type ACL string

const (
	// PublicReadWrite 开放读写
	PublicReadWrite ACL = "public-read-write"
	// PublicRead 开放读
	PublicRead ACL = "public_read"
	// Private 私有
	Private ACL = "private"
)

// Storage 存储类型
type Storage string

const (
	// Standard 标准存储
	Standard Storage = "standard"
	// InfrequentAccess 低频存储
	InfrequentAccess Storage = "infrequent-access"
	// Archive 归档存储
	Archive Storage = "archive"
	// ColdArchive 冷归档存储
	ColdArchive Storage = "cold-archive"
)
