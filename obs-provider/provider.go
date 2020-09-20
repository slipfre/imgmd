package provider

type Provider interface {
	CreateBucket() (err error)
	PutObjectFromFile() (err error)
}