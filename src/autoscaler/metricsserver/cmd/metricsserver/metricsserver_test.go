package main_test

import (
	. "code.cloudfoundry.org/app-autoscaler/src/autoscaler/testhelpers"
	"code.cloudfoundry.org/cfhttp"
	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("MetricsServer", func() {
	var (
		runner *MetricsServerRunner
	)

	BeforeEach(func() {
		runner = NewMetricsServerRunner()
	})

	AfterEach(func() {
		runner.KillWithFire()
	})

	Describe("Sending ws start and stop envelopes", func() {
		var ws websocket.Conn
		BeforeEach(func() {
			runner.Start()

			u := url.URL{Scheme: "wss", Host: fmt.Sprintf("localhost:%d", wsPort), Path: "/"}
			log.Printf("connecting to %s", u.String())
			//nolint:staticcheck  // SA1019 TODO: https://github.com/cloudfoundry/app-autoscaler-release/issues/548
			tlsConfig, err := cfhttp.NewTLSConfig(
				filepath.Join(testCertDir, "metricserver.crt"),
				filepath.Join(testCertDir, "metricserver.key"),
				filepath.Join(testCertDir, "autoscaler-ca.crt"),
			)
			ws, _, err := (&websocket.Dialer{
				Proxy:            http.ProxyFromEnvironment,
				HandshakeTimeout: 45 * time.Second,
				TLSClientConfig:  tlsConfig,
			}).Dial(u.String(), nil)

			FailOnError("failed open ws", err)

			DeferCleanup(func() {
				err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				FailOnError("Closed failed to close", err)
			})

		})
		It("try it out", func() {
			envelope := &loggregator_v2.Envelope{}
			envelope.Message = &loggregator_v2.Envelope_Timer{Timer: &loggregator_v2.Timer{Start: time.Now().UnixNano(), Stop: time.Now().Add(1 * time.Second).UnixNano()}}
			envelope.SourceId = "some-app"
			envelope.InstanceId = "0"
			envelope.GetTimer()

			message, err := proto.Marshal(envelope)
			FailOnError("proto.Marshal", err)

			err = ws.WriteMessage(websocket.BinaryMessage, message)
			FailOnError("WriteMessage", err)
		})

	})

	Describe("MetricsServer configuration check", func() {

		Context("with a missing config file", func() {
			BeforeEach(func() {
				runner.startCheck = ""
				runner.configPath = "bogus"
				runner.Start()
			})

			It("fails with an error", func() {
				Eventually(runner.Session).Should(Exit(1))
				Expect(runner.Session.Buffer()).To(Say("failed to open config file"))
			})
		})

		Context("with an invalid config file", func() {
			BeforeEach(func() {
				runner.startCheck = ""
				badfile, err := os.CreateTemp("", "bad-ms-config")
				Expect(err).NotTo(HaveOccurred())
				runner.configPath = badfile.Name()
				err = os.WriteFile(runner.configPath, []byte("bogus"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
				runner.Start()
			})

			AfterEach(func() {
				os.Remove(runner.configPath)
			})

			It("fails with an error", func() {
				Eventually(runner.Session).Should(Exit(1))
				Expect(runner.Session.Buffer()).To(Say("failed to read config file"))
			})
		})

		Context("with missing configuration", func() {
			BeforeEach(func() {
				runner.startCheck = ""
				missingConfig := cfg
				missingConfig.Server.Port = 7000 + GinkgoParallelProcess()
				missingConfig.Logging.Level = "debug"
				missingConfig.Collector.EnvelopeChannelSize = 0
				runner.configPath = writeConfig(&missingConfig).Name()
				runner.Start()
			})

			AfterEach(func() {
				os.Remove(runner.configPath)
			})

			It("should fail validation", func() {
				Eventually(runner.Session).Should(Exit(1))
				Expect(runner.Session.Buffer()).To(Say("failed to validate configuration"))
			})
		})
	})

	Describe("when interrupt is sent", func() {
		BeforeEach(func() {
			runner.Start()
		})

		It("should stop", func() {
			runner.Session.Interrupt()
			Eventually(runner.Session, 5).Should(Exit(0))
		})

	})

	Describe("MetricsServer REST API", func() {
		Context("when a request for metrics history comes", func() {
			BeforeEach(func() {
				runner.Start()
			})

			It("returns with a 200", func() {
				rsp, err := httpClient.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/apps/an-app-id/metric_histories/a-metric-type", msPort))
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				rsp.Body.Close()
			})
		})

	})

	Describe("when Health server is ready to serve RESTful API", func() {
		BeforeEach(func() {

			basicAuthConfig := cfg
			basicAuthConfig.Health.HealthCheckUsername = ""
			basicAuthConfig.Health.HealthCheckPassword = ""
			runner.configPath = writeConfig(&basicAuthConfig).Name()
			runner.Start()

		})
		Context("when a request to query health comes", func() {
			It("returns with a 200", func() {
				rsp, err := healthHttpClient.Get(fmt.Sprintf("http://127.0.0.1:%d", healthport))
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
				raw, _ := io.ReadAll(rsp.Body)
				healthData := string(raw)
				Expect(healthData).To(ContainSubstring("autoscaler_metricsserver_concurrent_http_request"))
				Expect(healthData).To(ContainSubstring("autoscaler_metricsserver_policyDB"))
				Expect(healthData).To(ContainSubstring("autoscaler_metricsserver_instanceMetricsDB"))
				Expect(healthData).To(ContainSubstring("go_goroutines"))
				Expect(healthData).To(ContainSubstring("go_memstats_alloc_bytes"))
				rsp.Body.Close()

			})
		})
	})

	Describe("when Health server is ready to serve RESTful API with basic Auth", func() {
		BeforeEach(func() {
			runner.Start()
		})

		Context("when username and password are incorrect for basic authentication during health check", func() {
			It("should return 401", func() {

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/health", healthport), nil)
				Expect(err).NotTo(HaveOccurred())

				req.SetBasicAuth("wrongusername", "wrongpassword")

				rsp, err := healthHttpClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when username and password are correct for basic authentication during health check", func() {
			It("should return 200", func() {

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/health", healthport), nil)
				Expect(err).NotTo(HaveOccurred())

				req.SetBasicAuth(cfg.Health.HealthCheckUsername, cfg.Health.HealthCheckPassword)

				rsp, err := healthHttpClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("when Health server is ready to serve RESTful API with basic Auth", func() {
		BeforeEach(func() {
			runner.Start()
		})

		Context("when username and password are incorrect for basic authentication during health check", func() {
			It("should return 401", func() {

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/health", healthport), nil)
				Expect(err).NotTo(HaveOccurred())

				req.SetBasicAuth("wrongusername", "wrongpassword")

				rsp, err := healthHttpClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when username and password are correct for basic authentication during health check", func() {
			It("should return 200", func() {

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/health", healthport), nil)
				Expect(err).NotTo(HaveOccurred())

				req.SetBasicAuth(cfg.Health.HealthCheckUsername, cfg.Health.HealthCheckPassword)

				rsp, err := healthHttpClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	//TODO : Add test cases for testing WebServer endpoints
})
