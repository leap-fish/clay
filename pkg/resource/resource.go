package resource

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"io"
)

var loaderInstance = loader{}

// Path is a non-filesystem path for a given resource.
type Path string

type LoadedResource struct {
	Path     Path
	instance any
}

// ResourceHandler represents a type that handles a specific file extension for loading.
type ResourceHandler interface {
	Load(reader io.ReadCloser) (any, error)
}

func RegisterHandler(prefix string, extension string, handler ResourceHandler) {
	loaderInstance.handleResourceType(prefix, extension, handler)
}

func LoadFromEmbedFolder(directory string, fs embed.FS) []error {
	log.WithField("directory", directory).Debug("Loading resources from embedded file system")
	return loaderInstance.loadFromFs(directory, fs)
}

// Get returns a resource, casted to the type of the generic argument given.
func Get[T any](path Path) T {
	var result T
	value, ok := loaderInstance.resources[path]
	if !ok {
		log.Error("Attempted to load resource, which was unavailable: ", path)
		return result
	}
	return value.instance.(T)
}
