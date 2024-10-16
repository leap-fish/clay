package plugins

import (
	"embed"
	"github.com/leap-fish/clay/pkg/clay"
	"github.com/leap-fish/clay/pkg/plugins/audio"
	"github.com/leap-fish/clay/pkg/plugins/render"
	"github.com/leap-fish/clay/pkg/plugins/resources"
)

// DefaultPlugins adds core engine plugins that are required for rendering.
// You can override this by simply adding them manually to your core.Plugin() command.
// DefaultPlugins includes the resource loading system, camera system and rendering + audio systems.
// This provides some basics, so you don't need to reinvent these unless your needs exceed that of the default
// implementation.
// Set `defaultDir` to "." if you want to use the current directory of the embed.FS.
func DefaultPlugins(fs embed.FS, defaultDir string) []clay.Plugin {
	return []clay.Plugin{
		resources.NewDefaultResourcesPlugin(defaultDir, fs),
		&DefaultCameraPlugin{},
		&audio.DefaultAudioPlugin{},
		&render.DefaultRendererPlugin{},
	}
}
