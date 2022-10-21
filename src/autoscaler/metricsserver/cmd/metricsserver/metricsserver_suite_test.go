package main_test

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "code.cloudfoundry.org/app-autoscaler/src/autoscaler/testhelpers"

	"code.cloudfoundry.org/cfhttp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	yaml "gopkg.in/yaml.v2"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/db"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/metricsserver/config"
)

var (
	msPath           string
	cfg              config.Config
	msPort           int
	wsPort           int
	healthport       int
	configFile       *os.File
	httpClient       *http.Client
	healthHttpClient *http.Client
)

const testCertDir = "../../../../../test-certs"

func TestMetricsServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MetricsServer Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	ms, err := gexec.Build("code.cloudfoundry.org/app-autoscaler/src/autoscaler/metricsserver/cmd/metricsserver", "-race")
	Expect(err).NotTo(HaveOccurred())

	msDB := OpenDb()

	_, err = msDB.Exec("DELETE FROM appinstancemetrics")
	FailOnError("clean appinstancemetrics", err)

	_, err = msDB.Exec("DELETE from policy_json")
	FailOnError("clean policy_json", err)

	AddAppPolicy(msDB, "an-app-id", "1234")

	err = msDB.Close()
	FailOnError("close db", err)
	return []byte(ms)
}, func(pathsByte []byte) {
	msPath = string(pathsByte)

	msPort = 7000 + GinkgoParallelProcess()
	wsPort = 7100 + GinkgoParallelProcess()
	healthport = 8000 + GinkgoParallelProcess()
	cfg.Server.Port = msPort
	cfg.Health.Port = healthport
	cfg.Logging.Level = "info"

	dbUrl := GetDbUrl()
	cfg.DB.InstanceMetricsDB = db.DatabaseConfig{
		URL:                   dbUrl,
		MaxOpenConnections:    10,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 10 * time.Second,
	}
	cfg.DB.PolicyDB = db.DatabaseConfig{
		URL:                   dbUrl,
		MaxOpenConnections:    10,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 10 * time.Second,
	}

	cfg.Collector.CollectInterval = 10 * time.Second
	cfg.Collector.RefreshInterval = 30 * time.Second
	cfg.Collector.SaveInterval = 5 * time.Second
	cfg.Collector.MetricCacheSizePerApp = 100
	cfg.Collector.TLS.KeyFile = filepath.Join(testCertDir, "metricserver.key")
	cfg.Collector.TLS.CertFile = filepath.Join(testCertDir, "metricserver.crt")
	cfg.Collector.TLS.CACertFile = filepath.Join(testCertDir, "autoscaler-ca.crt")
	cfg.HttpClientTimeout = 10 * time.Second
	cfg.NodeAddrs = []string{"localhost"}
	cfg.NodeIndex = 0
	cfg.Collector.WSPort = wsPort
	cfg.Collector.WSKeepAliveTime = 1 * time.Minute

	cfg.Collector.PersistMetrics = true
	cfg.Collector.EnvelopeProcessorCount = 5
	cfg.Collector.EnvelopeChannelSize = 1000
	cfg.Collector.MetricChannelSize = 1000

	cfg.Health.HealthCheckUsername = "metricsserverhealthcheckuser"
	cfg.Health.HealthCheckPassword = "metricsserverhealthcheckpassword"

	configFile = writeConfig(&cfg)

	//nolint:staticcheck  // SA1019 TODO: https://github.com/cloudfoundry/app-autoscaler-release/issues/548
	tlsConfig, err := cfhttp.NewTLSConfig(
		filepath.Join(testCertDir, "metricserver.crt"),
		filepath.Join(testCertDir, "metricserver.key"),
		filepath.Join(testCertDir, "autoscaler-ca.crt"))
	Expect(err).NotTo(HaveOccurred())
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	healthHttpClient = &http.Client{}
})

func OpenDb() *sqlx.DB {
	dbUrl := GetDbUrl()
	database, err := db.GetConnection(dbUrl)
	FailOnError("could not open connection to DB", err)

	msDB, err := sqlx.Open(database.DriverName, database.DSN)
	FailOnError("could not open DB", err)
	return msDB
}

func AddAppPolicy(db *sqlx.DB, appId, appGuid string) {
	policy := `
		{
 			"instance_min_count": 1,
  			"instance_max_count": 5
		}`
	query := db.Rebind("INSERT INTO policy_json(app_id, policy_json, guid) values(?, ?, ?)")
	_, err := db.Exec(query, appId, policy, appGuid)
	FailOnError("insert policy failed", err)
}

var _ = SynchronizedAfterSuite(func() {
	os.Remove(configFile.Name())
}, func() {
	gexec.CleanupBuildArtifacts()
})

func writeConfig(c *config.Config) *os.File {
	cfg, err := os.CreateTemp("", "ms")
	Expect(err).NotTo(HaveOccurred())
	defer cfg.Close()

	bytes, err := yaml.Marshal(c)
	Expect(err).NotTo(HaveOccurred())

	_, err = cfg.Write(bytes)
	Expect(err).NotTo(HaveOccurred())

	return cfg
}

type MetricsServerRunner struct {
	configPath string
	startCheck string
	Session    *gexec.Session
}

func NewMetricsServerRunner() *MetricsServerRunner {
	return &MetricsServerRunner{
		configPath: configFile.Name(),
		startCheck: "metricsserver.started",
	}
}

func (ms *MetricsServerRunner) Start() {
	// #nosec G204
	msSession, err := gexec.Start(exec.Command(
		msPath,
		"-c",
		ms.configPath,
	),
		gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[32m[mc]\x1b[0m ", GinkgoWriter),
		gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[32m[mc]\x1b[0m ", GinkgoWriter),
	)
	Expect(err).NotTo(HaveOccurred())

	if ms.startCheck != "" {
		Eventually(msSession.Buffer, 2).Should(gbytes.Say(ms.startCheck))
	}

	ms.Session = msSession
}

func (ms *MetricsServerRunner) Interrupt() {
	if ms.Session != nil {
		ms.Session.Interrupt().Wait(5 * time.Second)
	}
}

func (ms *MetricsServerRunner) KillWithFire() {
	if ms.Session != nil {
		ms.Session.Kill().Wait(5 * time.Second)
	}
}
