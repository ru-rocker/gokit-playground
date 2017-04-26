package auth

import (
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/log"
	"os"
	"github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd"
	"math/rand"
	"strconv"
	"time"
	"net"
)

func Register(consulAddress string,
	consulPort string,
	advertiseAddress string,
	advertisePort string) (registar sd.Registrar) {

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	client := consulsd.NewClient(ConsulClient(consulAddress, consulPort, logger))
	rand.Seed(time.Now().UTC().UnixNano())
	check := api.AgentServiceCheck{
		HTTP:     "http://" + advertiseAddress + ":" + advertisePort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Basic health checks",
	}

	port, _ := strconv.Atoi(advertisePort)
	num := rand.Intn(100) // to make service ID unique
	asr := api.AgentServiceRegistration{
		ID:      "auth" + strconv.Itoa(num), //unique service ID
		Name:    "auth",
		Address: advertiseAddress,
		Port:    port,
		Tags:    []string{"auth", "ru-rocker"},
		Check:   &check,
	}
	registar = consulsd.NewRegistrar(client, &asr, logger)
	return
}

//retrieve consul api client for make consulsd client or KV
func ConsulClient(consulAddress string, consulPort string, logger log.Logger) *api.Client {
	// Service discovery domain. In this example we use Consul.
	consulConfig := api.DefaultConfig()
	consulConfig.Address = net.JoinHostPort(consulAddress, consulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	return consulClient
}
