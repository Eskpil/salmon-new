package node

import (
	"context"

	"github.com/eskpil/salmon/vm/internal/node/config"
	"github.com/eskpil/salmon/vm/internal/node/tasks"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
)

type State struct {
	Client *rockferry.Client

	t *tasks.TaskList
}

func New(c *config.Config) (*State, error) {
	var err error
	state := new(State)

	client, err := rockferry.New(c.Url)
	if err != nil {
		return nil, err
	}

	state.t, err = tasks.NewTaskList(client, c.Id)
	if err != nil {
		return nil, err
	}

	state.Client = client

	return state, err
}

func (s *State) Watch(ctx context.Context) error {
	ctx = context.WithoutCancel(ctx)

	if err := s.startupTasks(); err != nil {
		return err
	}

	if err := s.watchStorageVolumes(ctx); err != nil {
		return err
	}

	if err := s.watchMachineRequests(ctx); err != nil {
		return err
	}

	return s.t.Run(ctx)
}

func (s *State) startupTasks() error {
	{
		task := new(tasks.SyncNodeTask)
		s.t.AppendUnbound(task)
	}

	{
		task := new(tasks.SyncStoragePoolsTask)
		s.t.AppendUnbound(task)
	}

	{
		task := new(tasks.SyncNetworksTask)
		s.t.AppendUnbound(task)
	}

	return nil
}

func (s *State) watchMachineRequests(ctx context.Context) error {
	go func() {
		requests, err := s.Client.MachineRequests().List(ctx, "", nil)
		if err != nil {
			return
		}

		for _, req := range requests {
			if req.Status.Phase == resource.PhaseRequested {
				task := new(tasks.CreateVirtualMachineTask)
				task.Request = req
				s.t.AppendBound(task)
			}
		}

		stream, err := s.Client.MachineRequests().Watch(ctx, "", nil)
		if err != nil {
			return
		}

		for {
			req := <-stream

			if req.Status.Phase == resource.PhaseRequested {
				task := new(tasks.CreateVirtualMachineTask)
				task.Request = req
				s.t.AppendBound(task)
			}

		}
	}()

	return nil
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
