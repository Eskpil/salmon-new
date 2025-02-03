package tasks

import (
	"context"
	"fmt"

	"github.com/eskpil/salmon/vm/internal/node/queries"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
)

type Executor struct {
	Libvirt   *queries.Client
	Rockferry *rockferry.Client

	NodeId string
}

type Task interface {
	Execute(context.Context, *Executor) error
}

type BoundTask interface {
	Execute(context.Context, *Executor) error
	Resource() *resource.Resource[any]
}

type TaskList struct {
	e            *Executor
	boundTasks   chan BoundTask
	unboundTasks chan Task
}

func NewTaskList(client *rockferry.Client, nodeId string) (*TaskList, error) {
	var err error
	list := new(TaskList)
	list.unboundTasks = make(chan Task, 100)
	list.boundTasks = make(chan BoundTask, 100)
	list.e = new(Executor)

	list.e.Libvirt, err = queries.NewClient()
	list.e.Rockferry = client
	list.e.NodeId = nodeId

	return list, err
}

func (t *TaskList) AppendBound(task BoundTask) {
	t.boundTasks <- task
}

func (t *TaskList) AppendUnbound(task Task) {
	t.unboundTasks <- task
}

func (t *TaskList) setResourcePhase(ctx context.Context, res *resource.Resource[any], phase resource.Phase, error string) error {
	generic := t.e.Rockferry.Generic(resource.ResourceKindAll)

	copy := new(resource.Resource[any])
	*copy = *res

	copy.Status.Phase = phase
	if error != "" && phase == resource.PhaseErrored {
		copy.Status.Error = new(string)
		*copy.Status.Error = error
	}

	err := generic.Patch(ctx, res, copy)
	return err
}

func (t *TaskList) Run(ctx context.Context) error {
	for {
		select {
		case task := <-t.unboundTasks:
			{
				if err := task.Execute(ctx, t.e); err != nil {
					fmt.Println("failed to execute task", err)
				}
			}
		case task := <-t.boundTasks:
			{

				if err := t.setResourcePhase(ctx, task.Resource(), resource.PhaseCreating, ""); err != nil {
					fmt.Println("could not set resource phase", err)
					continue
				}

				if err := task.Execute(ctx, t.e); err != nil {
					if err := t.setResourcePhase(ctx, task.Resource(), resource.PhaseErrored, err.Error()); err != nil {
						fmt.Println("could not set resource phase", err)
						continue
					}
				}

				if err := t.setResourcePhase(ctx, task.Resource(), resource.PhaseCreated, ""); err != nil {
					fmt.Println("could not set resource phase", err)
					continue
				}
			}
		}
	}
}
