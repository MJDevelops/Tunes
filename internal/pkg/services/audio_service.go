package services

import (
	"container/list"
	"context"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AudioService struct {
	ctx     context.Context
	queue   *audio.Queue
	queries *db.Queries
	elems   map[int64]*list.Element
}

func NewAudioService(queries *db.Queries) *AudioService {
	return &AudioService{
		queue:   audio.NewQueue(50),
		queries: queries,
	}
}

func (s *AudioService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *AudioService) AddToQueue(trackId int64) error {
	t, err := s.queries.GetTrack(context.TODO(), trackId)
	if err != nil {
		return err
	}

	d, _ := audio.NewDecoder(t.Path)
	f, _ := audio.NewAudioFile(d)

	if _, ok := s.elems[trackId]; !ok {
		s.elems[trackId] = s.queue.Add(&f)
	}

	return nil
}
