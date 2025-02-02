package node

import (
	"context"

	"github.com/eskpil/salmon/vm/internal/node/tasks"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
)

type State struct {
	Client *rockferry.Client

	t *tasks.TaskList
}

func New() (*State, error) {
	var err error
	state := new(State)

	client, err := rockferry.New()
	if err != nil {
		return nil, err
	}

	state.t, err = tasks.NewTaskList(client)
	if err != nil {
		return nil, err
	}

	state.Client = client

	return state, err
}

func (s *State) Watch(ctx context.Context) error {
	ctx = context.WithoutCancel(ctx)

	if err := s.watchStorageVolumes(ctx); err != nil {
		return err
	}

	return s.t.Run(ctx)
}

func (s *State) watchStorageVolumes(ctx context.Context) error {
	go func() {
		volumes, err := s.Client.StorageVolumes().List(ctx, "", nil)
		if err != nil {
			return
		}

		for _, vol := range volumes {
			if vol.Status.Phase == resource.PhaseRequested {
				task := new(tasks.CreateVolumeTask)
				task.Volume = vol
				s.t.AppendBound(task)
			}
		}

		stream, err := s.Client.StorageVolumes().Watch(ctx, "", nil)
		if err != nil {
			return
		}

		for {
			vol := <-stream

			if vol.Status.Phase == resource.PhaseRequested {
				task := new(tasks.CreateVolumeTask)
				task.Volume = vol
				s.t.AppendBound(task)
			}

		}
	}()

	return nil
}
