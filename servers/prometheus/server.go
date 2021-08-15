package prometheus

import (
	"net/http"

	"go.mrm.dev/venstar-monitor"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Monitor      *monitor.Monitor
	ListenAddr   string
	HTTPServeMux *http.ServeMux
}

type Server struct {
	listenAddr string
	monitor    *monitor.Monitor
	serveMux   *http.ServeMux
	registry   *prometheus.Registry
	collector  *venstarCollector
}

func (s *Server) setup() error {
	s.registry = prometheus.NewRegistry()
	s.collector = newVenstarCollector(s.monitor)
	err := s.registry.Register(s.collector)
	if err != nil {
		return errors.Wrap(err, "registering venstar collector")
	}

	promhandler := promhttp.InstrumentMetricHandler(
		s.registry, promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}),
	)

	s.serveMux.Handle("/metrics", promhandler)
	return nil
}

func (s *Server) Serve() error {
	return http.ListenAndServe(s.listenAddr, s.serveMux)
}

func NewServer(config Config) (*Server, error) {
	listen := ":9872"
	mux := http.DefaultServeMux

	if config.Monitor == nil {
		return nil, errors.New("monitor object expected")
	}
	if config.ListenAddr != "" {
		listen = config.ListenAddr
	}
	if config.HTTPServeMux != nil {
		mux = config.HTTPServeMux
	}

	server := &Server{
		listenAddr: listen,
		monitor:    config.Monitor,
		serveMux:   mux,
	}
	err := server.setup()
	if err != nil {
		return nil, err
	}

	return server, nil
}
