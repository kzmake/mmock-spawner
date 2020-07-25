package example

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"

	"github.com/kzmake/mmock-spawner/proto"
)

func Example() {
	ctx := context.Background()

	service := proto.NewProcessService("", grpc.NewClient())
	options := []client.CallOption{
		client.WithAddress("127.0.0.1:5000"),
	}

	spawnRes, err := service.Spawn(ctx, &proto.SpawnRequest{}, options...)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}
	fmt.Println(spawnRes.String())

	listRes, err := service.List(ctx, &proto.ListRequest{}, options...)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}
	fmt.Println(listRes.String())

	killAllRes, err := service.KillAll(ctx, &proto.KillAllRequest{}, options...)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}
	fmt.Println(killAllRes.String())

	// Output:
	// result:{pid:15 command:"/bin/mmock -console-ip 0.0.0.0 -server-ip 0.0.0.0 -config-path ./config"}
	// results:{pid:15 command:"/bin/mmock -console-ip 0.0.0.0 -server-ip 0.0.0.0 -config-path ./config"}
	// results:{pid:15 command:"/bin/mmock -console-ip 0.0.0.0 -server-ip 0.0.0.0 -config-path ./config"}
}
