package titan

import (
	"encoding/json"
	"github.com/Filecoin-Titan/titan-sdk-go-s3/config"
	"github.com/Filecoin-Titan/titan-sdk-go-s3/types"
	"github.com/pkg/errors"
	"net"
	"net/http"
)

type Service struct {
	cfg        config.Config
	httpClient *http.Client
}

func New(options config.Config) (*Service, error) {
	if options.TitanAddress == "" {
		return nil, errors.Errorf("the address cannot be empty")
	}

	if options.CandidateID == "" {
		return nil, errors.Errorf("the L1 node ID cannot be empty")
	}

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}

	s := &Service{
		cfg:        options,
		httpClient: defaultHttpClient(conn, options.Timeout),
	}

	return s, nil
}

type params []interface{}

// GetScheduler get the scheduler URL from the specified candidate ID
func (s *Service) GetScheduler() (string, error) {
	serializedParams, err := json.Marshal(params{s.cfg.CandidateID})
	if err != nil {
		return "", errors.Errorf("marshaling params failed: %v", err)
	}

	req := Request{
		Jsonrpc: "2.0",
		ID:      "1",
		Method:  "titan.GetSchedulerWithNode",
		Params:  serializedParams,
	}

	resp, err := handleRequest(s.httpClient, s.cfg.TitanAddress, req)
	if err != nil {
		return "", err
	}

	out, ok := resp.Result.(string)
	if !ok {
		return "", errors.Errorf("type assign string failed")
	}

	return out, nil
}

// GetLocalMinioEndpoint get local minio endpoint
func (s *Service) GetLocalMinioEndpoint(schedulerURL string) (*types.MinioConfig, error) {
	serializedParams, err := json.Marshal(params{s.cfg.CandidateID})
	if err != nil {
		return nil, errors.Errorf("marshaling params failed: %v", err)
	}

	req := Request{
		Jsonrpc: "2.0",
		ID:      "1",
		Method:  "titan.GetMinioConfigFromCandidate",
		Params:  serializedParams,
	}

	resp, err := handleRequest(s.httpClient, schedulerURL, req)
	if err != nil {
		return nil, err
	}

	var minioCfg types.MinioConfig
	err = decodeResult(resp.Result, &minioCfg)
	if err != nil {
		return nil, err
	}

	return &minioCfg, nil
}
