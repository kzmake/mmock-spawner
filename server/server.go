package server

import (
	"context"
	"sync"

	klogger "github.com/kzmake/micro-kit/pkg/wrapper/logger"
	cli "github.com/micro/cli/v2"
	micro "github.com/micro/go-micro/v2"
	mserver "github.com/micro/go-micro/v2/server"

	"github.com/kzmake/mmock-spawner/proto"
)

func waitgroup(waitGroup *sync.WaitGroup) mserver.HandlerWrapper {
	return func(h mserver.HandlerFunc) mserver.HandlerFunc {
		return func(ctx context.Context, req mserver.Request, rsp interface{}) error {
			waitGroup.Add(1)
			defer waitGroup.Done()
			return h(ctx, req, rsp)
		}
	}
}

// New はサーバーを生成します。
func New(
	name, version string,
	controller proto.ProcessServiceHandler,
) (micro.Service, error) {
	wg := new(sync.WaitGroup)
	service := micro.NewService(
		micro.Name(name),
		micro.Version(version),
		micro.Address("0.0.0.0:5000"),

		micro.WrapHandler(
			waitgroup(wg),
			klogger.NewHandlerWrapper(),
		),

		micro.AfterStop(func() error {
			wg.Wait()
			return controller.KillAll(context.Background(), &proto.KillAllRequest{}, &proto.KillAllResponse{})
		}),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			return nil
		}),
	)

	if err := proto.RegisterProcessServiceHandler(service.Server(), controller); err != nil {
		return nil, err
	}

	return service, nil
}
