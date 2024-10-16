package audio

import (
	eaudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/leap-fish/clay/pkg/components/audio"
	"github.com/leap-fish/clay/pkg/resource"
	log "github.com/sirupsen/logrus"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"time"
)

const bytesPerAudioSample = 4
const sampleRate = 48_000

type SoundBytes []byte

type DefaultAudioEffectSystem struct {
	audioEffectQuery *donburi.Query
	ctx              *eaudio.Context
}

func (d *DefaultAudioEffectSystem) Init(w donburi.World) {
	d.ctx = eaudio.NewContext(sampleRate)
}

func NewDefaultAudioEffectSystem() *DefaultAudioEffectSystem {
	return &DefaultAudioEffectSystem{
		audioEffectQuery: donburi.NewQuery(
			filter.Contains(audio.Component)),
	}
}

func (d *DefaultAudioEffectSystem) Update(w donburi.World, dt time.Duration) {
	d.audioEffectQuery.Each(w, func(entry *donburi.Entry) {
		sfx := audio.Component.Get(entry)
		streamBytes := resource.Get[SoundBytes](sfx.Path)
		byteLen := len(streamBytes)
		total := time.Second * time.Duration(byteLen) / bytesPerAudioSample / sampleRate

		// Load the player
		if sfx.Player == nil {
			sfx.Player = d.ctx.NewPlayerFromBytes(streamBytes)
		}

		if sfx.Player.Volume() != sfx.Volume {
			sfx.Player.SetVolume(sfx.Volume)
		}

		// Remove the effect when the sound has finished
		if sfx.Player.Position() >= total {
			w.Remove(entry.Entity())
		}

		// Start playing if it isn't.
		if !sfx.Player.IsPlaying() && sfx.Player.Position() == 0 {
			_ = sfx.Player.SetPosition(0)
			log.
				WithField("streamBytesLen", len(streamBytes)).
				WithField("volume", sfx.Player.Volume()).
				WithField("totalLength", total).
				Tracef("Playing audio at path %s", sfx.Path)
			sfx.Player.Play()
		}

	})
}
