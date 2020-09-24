package provider

// Bucket OBS 服务供应商的 bucket 接口，用于提供对象存储功能
type Bucket interface {
	PutObjectFromFile(objectKey, filePath string, options ...ObjectOption) (err error)
}

// ObjectOption Bucket 相关的可选参数
type ObjectOption func(config *OptionConfig)

// WithACL Object 的 ACL 可选参数
func WithACL(acl ACL) ObjectOption {
	return func(config *OptionConfig) {
		config.acl = acl
	}
}

// WithStorage Object 的 Storage 可选参数
func WithStorage(storage Storage) ObjectOption {
	return func(config *OptionConfig) {
		config.storage = storage
	}
}

// WithRedundancyType Object 的 DataRedundancyTypoe 可选参数
func WithRedundancyType(redundancyType DataRedundancyType) ObjectOption {
	return func(config *OptionConfig) {
		config.redundancyType = redundancyType
	}
}

// OptionConfig Object 相关的配置参数
type OptionConfig struct {
	acl            ACL
	storage        Storage
	redundancyType DataRedundancyType
}

// ACL Bucket 的读写访问权限
type ACL string

const (
	// PublicReadWrite 开放读写
	PublicReadWrite ACL = "public-read-write"
	// PublicRead 开放读，默认值
	PublicRead ACL = "public_read"
	// Private 私有
	Private ACL = "private"
)

// Storage 存储类型
type Storage string

const (
	// Standard 标准存储，默认值
	Standard Storage = "standard"
	// InfrequentAccess 低频存储
	InfrequentAccess Storage = "infrequent-access"
	// Archive 归档存储
	Archive Storage = "archive"
	// ColdArchive 冷归档存储
	ColdArchive Storage = "cold-archive"
)

// DataRedundancyType 数据容灾类型
type DataRedundancyType string

const (
	// LRS 同可用区多存储设备容灾，默认值
	LRS DataRedundancyType = "LRS"
	// ZRS 多可用区容灾
	ZRS DataRedundancyType = "ZRS"
)
