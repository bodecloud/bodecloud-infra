package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/client"

	"cluster/infra/failover"
)

func main() {
	logger := log.New(os.Stdout, "[failover-agent] ", log.LstdFlags|log.Lmsgprefix)

	cfg := loadConfig()
	if cfg.Domain == "" {
		logger.Fatal("DOMAIN is required")
	}
	if cfg.LocalNode == "" {
		hostname, _ := os.Hostname()
		cfg.LocalNode = strings.Split(hostname, ".")[0]
		logger.Printf("TS_HOSTNAME unset — using hostname %s", cfg.LocalNode)
	}

	dockerHost := envOr("DOCKER_HOST", "unix:///var/run/docker.sock")
	cli, err := client.NewClientWithOpts(
		client.WithHost(dockerHost),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		logger.Fatalf("docker client: %v", err)
	}
	defer cli.Close()

	agent, err := failover.NewAgent(cfg, cli, logger)
	if err != nil {
		logger.Fatalf("agent: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := cli.Ping(r.Context()); err != nil {
			http.Error(w, "docker unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	addr := cfg.HealthListenAddr
	if addr == "" {
		addr = ":8082"
	}
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		logger.Printf("health listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("health server: %v", err)
		}
	}()

	logger.Printf("starting failover-agent domain=%s node=%s main=%s peers=%v enabled=%v replicaEnsure=%v",
		cfg.Domain, cfg.LocalNode, cfg.MainHost, cfg.PeerHosts, cfg.Enabled, cfg.ReplicaEnsure)
	if len(cfg.PeerHosts) == 0 {
		logger.Printf("WARN: FAILOVER_PEER_HOSTS is empty — no peer URLs or inventory; set peers for failover")
	}

	err = agent.Run(ctx)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	if err != nil && err != context.Canceled {
		logger.Fatalf("run: %v", err)
	}
	logger.Printf("shutdown complete")
}

func loadConfig() failover.AgentConfig {
	configPath := envOr("CONFIG_PATH", "./volumes")
	registryPath := envOr("PLACEMENT_REGISTRY", filepath.Join(configPath, "placement", "services.yaml"))
	traefikOut := envOr("FAILOVER_TRAEFIK_OUTPUT", filepath.Join(configPath, "traefik", "dynamic", "failover-fallbacks.yaml"))

	peers := splitCSV(os.Getenv("FAILOVER_PEER_HOSTS"))
	port, _ := strconv.Atoi(envOr("FAILOVER_REMOTE_DOCKER_PORT", "2375"))
	maxRestarts, _ := strconv.Atoi(envOr("FAILOVER_MAX_LOCAL_RESTARTS", "3"))
	reconcileSec, _ := strconv.Atoi(envOr("FAILOVER_RECONCILE_SECONDS", "30"))

	return failover.AgentConfig{
		Enabled:                   envBool("FAILOVER_ENABLED", true),
		Domain:                    os.Getenv("DOMAIN"),
		LocalNode:                 os.Getenv("TS_HOSTNAME"),
		MainHost:                  envOr("FAILOVER_MAIN_HOST", "micklethefickle"),
		PeerHosts:                 peers,
		RegistryPath:              registryPath,
		TraefikOutputPath:         traefikOut,
		RemoteDockerPort:          port,
		RemoteDockerTLS:           envBool("FAILOVER_REMOTE_DOCKER_TLS", false),
		MaxLocalRestarts:          maxRestarts,
		ReconcileInterval:         time.Duration(reconcileSec) * time.Second,
		StopReplicasOnIntentional: envBool("FAILOVER_STOP_REPLICAS_ON_INTENTIONAL_STOP", false),
		ReplicaEnsure:             envBool("FAILOVER_REPLICA_ENSURE", false),
		HealthListenAddr:          envOr("FAILOVER_HEALTH_ADDR", ":8082"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
	}
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
