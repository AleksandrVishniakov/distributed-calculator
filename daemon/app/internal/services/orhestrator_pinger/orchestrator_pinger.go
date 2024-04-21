package orhestrator_pinger

import (
	"context"
	"fmt"
	"log"
	"time"

	orchestrator "github.com/AleksandrVishniakov/dc-protos/gen/go/orchestrator/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrchestratorPinger struct {
	id               uint64
	client           orchestrator.OrchestratorClient
	host             string
	executors        int
	orchestratorHost string
	requestJSON      []byte
}

func NewOrchestratorPinger(
	ctx context.Context,
	id uint64,
	host string,
	gRPCHost string,
	executors int,
) (*OrchestratorPinger, error) {
	cc, err := grpc.DialContext(
		ctx,
		gRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	pinger := &OrchestratorPinger{
		client:    orchestrator.NewOrchestratorClient(cc),
		id:        id,
		executors: executors,
		host:      host,
	}

	time.AfterFunc(2*time.Second, func() {
		err := pinger.SendPing(ctx)
		if err != nil {
			log.Fatal(err)
		}
	})

	return pinger, err
}

func (d *OrchestratorPinger) MustPingOrchestrator(
	ctx context.Context,
	period time.Duration,
) {
	ticker := time.NewTicker(period)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := d.SendPing(ctx)
				if err != nil {
					log.Fatalf("orchestrator ping error: %s", err.Error())
				}
			}
		}
	}()
}

func (d *OrchestratorPinger) SendPing(ctx context.Context) error {
	resp, err := d.client.RegisterWorker(ctx, &orchestrator.WorkerRegisterRequest{
		Id:        d.id,
		Url:       d.host,
		Executors: uint32(d.executors),
	})

	if err != nil {
		return fmt.Errorf("grpc client error: %w", err)
	}

	if !resp.Ok {
		return fmt.Errorf("grpc client error: response is not ok")
	}

	return nil
}
