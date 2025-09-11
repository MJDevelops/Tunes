package audio

import (
	"time"
)

type TrackMeta struct {
	Title    string
	Artist   string
	Duration time.Duration
	Album    string
	Genre    string
}
