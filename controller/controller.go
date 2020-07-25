package controller

import (
	"context"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/micro/go-micro/v2/errors"

	"github.com/kzmake/mmock-spawner/proto"
	"github.com/kzmake/mmock-spawner/view"
)

var (
	ps = map[int]*exec.Cmd{}
	mu = &sync.RWMutex{}
)

type spawner struct {
	entrypoint  string
	defaultArgs []string
	renderer    view.Process
}

// interfaces
var _ proto.ProcessServiceHandler = (*spawner)(nil)

// New は process に関する controller を生成します。
func New(
	entrypoint string,
	defaultArgs []string,
	renderer view.Process,
) proto.ProcessServiceHandler {
	return &spawner{
		entrypoint:  entrypoint,
		defaultArgs: defaultArgs,
		renderer:    renderer,
	}
}

func (c *spawner) Spawn(ctx context.Context, req *proto.SpawnRequest, res *proto.SpawnResponse) error {
	var p *exec.Cmd
	if len(req.Args) != 0 {
		p = exec.Command(c.entrypoint, req.Args...)
	} else {
		p = exec.Command(c.entrypoint, c.defaultArgs...)
	}

	// exec command
	if err := p.Start(); err != nil {
		return errors.InternalServerError("Unexpected", "cannot exec process(%s)", p.String())
	}

	// zombie process 対策
	go func() { _ = p.Wait() }()
	time.Sleep(500 * time.Millisecond)

	// kill -0 pid
	if err := p.Process.Signal(syscall.Signal(0)); err != nil {
		out, _ := p.Output()
		return errors.BadRequest("CannotExecuteProcess", "cannot exec process(%s): %s: %+v", p.String(), string(out), err)
	}

	mu.Lock()
	defer mu.Unlock()
	ps[p.Process.Pid] = p

	return c.renderer.Spawn(ctx, res, p)
}

func (c *spawner) List(ctx context.Context, _ *proto.ListRequest, res *proto.ListResponse) error {
	mu.RLock()
	defer mu.RUnlock()

	processes := []*exec.Cmd{}
	for _, p := range ps {
		p := p

		processes = append(processes, p)
	}

	return c.renderer.List(ctx, res, processes)
}

func (c *spawner) Kill(ctx context.Context, req *proto.KillRequest, res *proto.KillResponse) error {
	mu.Lock()
	defer mu.Unlock()

	p, ok := ps[int(req.Pid)]
	if !ok {
		return errors.NotFound("NotFoundProcess", "process(pid: %d) is not found", req.Pid)
	}

	if err := p.Process.Kill(); err != nil {
		return errors.InternalServerError("InternalServerError", "cannot kill process(pid: %d)", req.Pid)
	}

	delete(ps, p.Process.Pid)

	return c.renderer.Kill(ctx, res, p)
}

func (c *spawner) KillAll(ctx context.Context, _ *proto.KillAllRequest, res *proto.KillAllResponse) error {
	mu.Lock()
	defer mu.Unlock()

	processes := []*exec.Cmd{}
	for _, p := range ps {
		p := p

		_ = p.Process.Kill()
		processes = append(processes, p)
	}
	ps = map[int]*exec.Cmd{}

	return c.renderer.KillAll(ctx, res, processes)
}
