package view

import (
	"context"
	"os/exec"

	"github.com/kzmake/mmock-spawner/proto"
)

// Process は process の view の interface です。
type Process interface {
	Spawn(context.Context, *proto.SpawnResponse, *exec.Cmd) error
	List(context.Context, *proto.ListResponse, []*exec.Cmd) error
	Kill(context.Context, *proto.KillResponse, *exec.Cmd) error
	KillAll(context.Context, *proto.KillAllResponse, []*exec.Cmd) error
}

// Process は process の view に関する定義です。
type renderer struct {
}

// New は process に関する view を生成します。
func New() Process { return &renderer{} }

func (r *renderer) Spawn(_ context.Context, res *proto.SpawnResponse, p *exec.Cmd) error {
	res.Result = &proto.Process{
		Pid:     int32(p.Process.Pid),
		Command: p.String(),
	}

	return nil
}

func (r *renderer) List(_ context.Context, res *proto.ListResponse, ps []*exec.Cmd) error {
	results := []*proto.Process{}
	for _, p := range ps {
		p := p

		process := &proto.Process{
			Pid:     int32(p.Process.Pid),
			Command: p.String(),
		}

		results = append(results, process)
	}

	res.Results = results

	return nil
}

func (r *renderer) Kill(_ context.Context, res *proto.KillResponse, p *exec.Cmd) error {
	res.Result = &proto.Process{
		Pid:     int32(p.Process.Pid),
		Command: p.String(),
	}

	return nil
}

func (r *renderer) KillAll(_ context.Context, res *proto.KillAllResponse, ps []*exec.Cmd) error {
	results := []*proto.Process{}
	for _, p := range ps {
		p := p

		process := &proto.Process{
			Pid:     int32(p.Process.Pid),
			Command: p.String(),
		}

		results = append(results, process)
	}

	res.Results = results

	return nil
}
