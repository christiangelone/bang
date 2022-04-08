package resource

type Resource interface {
	IsDir() bool
	Name() string
	Path() string
	FullPath() string
	Read() ([]byte, error)
}
