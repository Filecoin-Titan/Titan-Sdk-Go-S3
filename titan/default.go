package titan

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"net"
	"net/http"
	"time"
)

func defaultTLSConf() *tls.Config {
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
		NextProtos:         []string{http3.NextProtoH3},
	}
}

func defaultQUICConfig() *quic.Config {
	return &quic.Config{
		KeepAlivePeriod: 30 * time.Second,
		MaxIdleTimeout:  30 * time.Second,
		Allow0RTT:       func(net.Addr) bool { return true },
	}
}

func defaultHttpClient(conn net.PacketConn, timeout time.Duration) *http.Client {
	return &http.Client{Transport: &http3.RoundTripper{
		TLSClientConfig: defaultTLSConf(),
		QuicConfig:      defaultQUICConfig(),
		Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
			address, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				return nil, err
			}
			return quic.DialEarlyContext(ctx, conn, address, "localhost", tlsCfg, cfg)
		},
	}, Timeout: timeout}
}
