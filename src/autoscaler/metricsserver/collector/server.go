package collector

import (
	"fmt"
	"net/http"
	"os"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/healthendpoint"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/routes"

	"code.cloudfoundry.org/lager/v3"
	"github.com/gorilla/mux"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

type VarsFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string)

func (vh VarsFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vh(w, r, vars)
}

func NewServer(logger lager.Logger, serverConfig *ServerConfig, query MetricQueryFunc, httpStatusCollector healthendpoint.HTTPStatusCollector) (ifrit.Runner, error) {
	mh := NewMetricHandler(logger, serverConfig.NodeIndex, serverConfig.NodeAddrs, query)
	httpStatusCollectMiddleware := healthendpoint.NewHTTPStatusCollectMiddleware(httpStatusCollector)

	r := routes.MetricsCollectorRoutes()
	r.Use(httpStatusCollectMiddleware.Collect)
	r.Get(routes.GetMetricHistoriesRouteName).Handler(VarsFunc(mh.GetMetricHistories))

	var addr string
	if os.Getenv("APP_AUTOSCALER_TEST_RUN") == "true" {
		addr = fmt.Sprintf("localhost:%d", serverConfig.Port)
	} else {
		addr = fmt.Sprintf("0.0.0.0:%d", serverConfig.Port)
	}

	var runner ifrit.Runner
	if (serverConfig.TLS.KeyFile == "") || (serverConfig.TLS.CertFile == "") {
		runner = http_server.New(addr, r)
	} else {
		tlsConfig, err := serverConfig.TLS.CreateServerConfig()
		if err != nil {
			logger.Error("failed-new-server-new-tls-config", err, lager.Data{"tls": serverConfig.TLS})
			return nil, fmt.Errorf("metrics collector tls error: %w", err)
		}
		runner = http_server.NewTLSServer(addr, r, tlsConfig)
	}

	logger.Info("http-server-created", lager.Data{"serverConfig": serverConfig})
	return runner, nil
}
