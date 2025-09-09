package audio

import (
	"time"
)

type TrackMeta interface {
	Title() string
	Artist() string
	Duration() time.Duration
	Album() string
	Genre() string
}
