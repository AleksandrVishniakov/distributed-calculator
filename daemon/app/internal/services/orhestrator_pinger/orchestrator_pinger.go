package orhestrator_pinger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"log"
	"net/http"
	"time"
)

type OrchestratorPinger struct {
	host             string
	orchestratorHost string
	requestJSON      []byte
}

func NewOrchestratorPinger(id int, host string, executors int, orchestratorHost string) (*OrchestratorPinger, error) {
	requestBody := &dto.OrchestratorPingDTO{
		Id:        id,
		Url:       host,
		Executors: executors,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	pinger := &OrchestratorPinger{
		host:             host,
		orchestratorHost: orchestratorHost,
		requestJSON:      requestJSON,
	}

	time.AfterFunc(2*time.Second, func() {
		err := pinger.SendPing()
		if err != nil {
			log.Fatal(err)
		}
	})

	return pinger, err
}

func (d *OrchestratorPinger) MustPingOrchestrator(period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := d.SendPing()
				if err != nil {
					log.Fatalf("orchestrator ping error: %s", err.Error())
				}
			}
		}
	}()
}

func (d *OrchestratorPinger) SendPing() error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, d.orchestratorHost+"/api/worker", bytes.NewBuffer(d.requestJSON))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("orchestrator ping failed with status " + resp.Status)
	}

	return nil
}
