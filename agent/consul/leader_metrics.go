package consul

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/prometheus"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/agent/connect/ca"
	"github.com/hashicorp/consul/logging"
	"github.com/hashicorp/consul/tlsutil"
)

var metricsKeyMeshRootCAExpiry = []string{"mesh", "active-root-ca", "expiry"}
var metricsKeyMeshActiveSigningCAExpiry = []string{"mesh", "active-signing-ca", "expiry"}

var CertExpirationGauges = []prometheus.GaugeDefinition{
	{
		Name: metricsKeyMeshRootCAExpiry,
		Help: "Seconds until the service mesh root certificate expires. Updated every hour",
	},
	{
		Name: metricsKeyMeshActiveSigningCAExpiry,
		Help: "Seconds until the service mesh signing certificate expires. Updated every hour",
	},
	{
		Name: metricsKeyAgentTLSCertExpiry,
		Help: "Seconds until the agent tls certificate expires. Updated every hour",
	},
}

func rootCAExpiryMonitor(s *Server) CertExpirationMonitor {
	return CertExpirationMonitor{
		Key: metricsKeyMeshRootCAExpiry,
		Labels: []metrics.Label{
			{Name: "datacenter", Value: s.config.Datacenter},
		},
		Logger: s.logger.Named(logging.Connect),
		Query: func() (time.Duration, error) {
			return getRootCAExpiry(s)
		},
	}
}

func getRootCAExpiry(s *Server) (time.Duration, error) {
	state := s.fsm.State()
	_, root, err := state.CARootActive(nil)
	switch {
	case err != nil:
		return 0, fmt.Errorf("failed to retrieve root CA: %w", err)
	case root == nil:
		return 0, fmt.Errorf("no active root CA")
	}

	return time.Until(root.NotAfter), nil
}

func signingCAExpiryMonitor(s *Server) CertExpirationMonitor {
	isPrimary := s.config.Datacenter == s.config.PrimaryDatacenter
	if isPrimary {
		return CertExpirationMonitor{
			Key: metricsKeyMeshActiveSigningCAExpiry,
			Labels: []metrics.Label{
				{Name: "datacenter", Value: s.config.Datacenter},
			},
			Logger: s.logger.Named(logging.Connect),
			Query: func() (time.Duration, error) {
				provider, _ := s.caManager.getCAProvider()

				if _, ok := provider.(ca.PrimaryUsesIntermediate); !ok {
					return getActiveIntermediateExpiry(s)
				}

				return getRootCAExpiry(s)

			},
		}
	}

	return CertExpirationMonitor{
		Key: metricsKeyMeshActiveSigningCAExpiry,
		Labels: []metrics.Label{
			{Name: "datacenter", Value: s.config.Datacenter},
		},
		Logger: s.logger.Named(logging.Connect),
		Query: func() (time.Duration, error) {
			return getActiveIntermediateExpiry(s)
		},
	}
}

func getActiveIntermediateExpiry(s *Server) (time.Duration, error) {
	state := s.fsm.State()
	_, root, err := state.CARootActive(nil)
	if err != nil {
		return 0, err
	}

	// the CA used in a secondary DC is the active intermediate,
	// which is the last in the IntermediateCerts stack
	if len(root.IntermediateCerts) == 0 {
		return 0, errors.New("no intermediate available")
	}
	cert, err := connect.ParseCert(root.IntermediateCerts[len(root.IntermediateCerts)-1])
	if err != nil {
		return 0, err
	}
	return time.Until(cert.NotAfter), nil
}

type CertExpirationMonitor struct {
	Key    []string
	Labels []metrics.Label
	Logger hclog.Logger
	// Query is called at each interval. It should return the duration until the
	// certificate expires, or an error if the query failed.
	Query func() (time.Duration, error)
}

const certExpirationMonitorInterval = time.Hour

func (m CertExpirationMonitor) Monitor(ctx context.Context) error {
	ticker := time.NewTicker(certExpirationMonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			d, err := m.Query()
			if err != nil {
				m.Logger.Warn("failed to emit certificate expiry metric", "error", err)
			}
			expiry := d / time.Second
			metrics.SetGaugeWithLabels(m.Key, float32(expiry), m.Labels)
		}
	}
}

var metricsKeyAgentTLSCertExpiry = []string{"agent", "tls", "cert", "expiry"}

// AgentTLSCertExpirationMonitor returns a CertExpirationMonitor which will
// monitor the expiration of the certificate used for agent TLS.
func AgentTLSCertExpirationMonitor(c *tlsutil.Configurator, logger hclog.Logger, dc string) CertExpirationMonitor {
	return CertExpirationMonitor{
		Key: metricsKeyAgentTLSCertExpiry,
		Labels: []metrics.Label{
			{Name: "node", Value: c.Base().NodeName},
			{Name: "datacenter", Value: dc},
		},
		Logger: logger,
		Query: func() (time.Duration, error) {
			raw := c.Cert()
			if raw == nil {
				return 0, fmt.Errorf("tls not enabled")
			}

			cert, err := x509.ParseCertificate(raw.Certificate[0])
			if err != nil {
				return 0, fmt.Errorf("failed to parse agent tls cert: %w", err)
			}
			return time.Until(cert.NotAfter), nil
		},
	}
}
