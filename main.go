package main

import (
	"os"

	"golang.org/x/xerrors"

	"github.com/kzmake/micro-kit/pkg/logger/technical"

	"github.com/kzmake/mmock-spawner/controller"
	"github.com/kzmake/mmock-spawner/server"
	"github.com/kzmake/mmock-spawner/view"
)

var (
	name    = "kzmake.mmockspawner.service.v1"
	version = "v0.1.0"

	entrypoint   = "mmock"
	defatultArgs = []string{"-console-ip", "0.0.0.0", "-server-ip", "0.0.0.0", "-config-path", "./config"}
)

func main() {
	v := view.New()
	c := controller.New(entrypoint, defatultArgs, v)
	s, err := server.New(name, version, c)
	if err != nil {
		technical.Errorf("%+v", xerrors.Errorf("serverの生成に失敗しました: %w", err))
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		technical.Errorf("%+v", xerrors.Errorf("serverの起動に失敗しました: %w", err))
		os.Exit(1)
	}
}
