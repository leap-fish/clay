package resources

import (
	"embed"
	"github.com/leap-fish/clay/pkg/clay"
	"github.com/leap-fish/clay/pkg/events"
	defaults "github.com/leap-fish/clay/pkg/plugins/resources/defaults"
	"github.com/leap-fish/clay/pkg/resource"
	log "github.com/sirupsen/logrus"
	"math"
)

type DefaultResourcesPlugin struct {
	FileSystem embed.FS
	Path       string
}

func NewDefaultResourcesPlugin(path string, fs embed.FS) *DefaultResourcesPlugin {
	return &DefaultResourcesPlugin{
		Path:       path,
		FileSystem: fs,
	}
}

func (r *DefaultResourcesPlugin) Order() int {
	return math.MinInt32
}

func (r *DefaultResourcesPlugin) Build(core *clay.Core) {
	log.Info("Registering default handlers")
	resource.RegisterHandler("shader", ".kage", &defaults.KageDefaultHandler{})
	resource.RegisterHandler("shader", ".kage.go", &defaults.KageDefaultHandler{})
	resource.RegisterHandler("image", ".png", &defaults.PngDefaultHandler{})
	resource.RegisterHandler("font", ".ttf", &defaults.TtfDefaultHandler{})
	resource.RegisterHandler("sfx", ".ogg", &defaults.OggDefaultHandler{})
}

func (r *DefaultResourcesPlugin) Ready(core *clay.Core) {
	resourceErrs := resource.LoadFromEmbedFolder(r.Path, r.FileSystem)
	if len(resourceErrs) > 0 {
		for _, err := range resourceErrs {
			log.WithError(err).Error("File load failed")
		}
		log.
			WithField("path", r.Path).
			WithField("fs", r.FileSystem).
			Warnf("Unable to load %d files from embedded file system", len(resourceErrs))
	}

	events.ResourcePluginLoaded.Publish(core.World, 0)
}
