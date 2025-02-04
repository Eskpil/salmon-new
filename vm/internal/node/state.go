package node

import (
	"context"

	"github.com/eskpil/salmon/vm/internal/node/config"
	"github.com/eskpil/salmon/vm/internal/node/tasks"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/status"
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

	if err := s.watchMachines(ctx); err != nil {
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
		task := new(tasks.SyncStorageVolumesTask)
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

		stream, err := s.Client.MachineRequests().Watch(ctx, rockferry.WatchActionPut, "", nil)
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

func (s *State) watchMachines(ctx context.Context) error {
	go func() {
		stream, err := s.Client.Machines().Watch(ctx, rockferry.WatchActionDelete, "", nil)
		if err != nil {
			return
		}

		for {
			machine := <-stream
			task := new(tasks.DeleteVmTask)
			task.Machine = machine
			s.t.AppendUnbound(task)
		}
	}()

	return nil
}

func (s *State) watchStorageVolumes(ctx context.Context) error {
	volumes, err := s.Client.StorageVolumes().List(ctx, "", nil)
	if err != nil && !status.Is(err, status.ErrNoResults) {
		return err
	}

	for _, vol := range volumes {
		if vol.Status.Phase == resource.PhaseRequested {
			task := new(tasks.CreateVolumeTask)
			task.Volume = vol
			s.t.AppendBound(task)
		}
	}

	// TODO: Combine?
	go func() {
		stream, err := s.Client.StorageVolumes().Watch(ctx, rockferry.WatchActionPut, "", nil)
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

	go func() {
		stream, err := s.Client.StorageVolumes().Watch(ctx, rockferry.WatchActionDelete, "", nil)
		if err != nil {
			return
		}

		for {
			vol := <-stream

			task := new(tasks.DeleteVolumeTask)
			task.Volume = vol
			s.t.AppendUnbound(task)

		}
	}()

	return nil
}
