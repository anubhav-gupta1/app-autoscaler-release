package cf_test

import (
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/cf"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
	. "code.cloudfoundry.org/app-autoscaler/src/autoscaler/testhelpers"
	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	"net/http"
)

var _ = Describe("Cf client Service Instances", func() {

	var (
		conf            *cf.Config
		cfc             *cf.Client
		fakeCC          *MockServer
		fakeLoginServer *Server
		err             error
		logger          lager.Logger
	)

	var setCfcClient = func(maxRetries int) {
		conf = &cf.Config{}
		conf.API = fakeCC.URL()
		conf.MaxRetries = maxRetries
		conf.MaxRetryWaitMs = 1
		cfc = cf.NewCFClient(conf, logger, clock.NewClock())
		err = cfc.Login()
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		fakeCC = NewMockServer()
		fakeLoginServer = NewServer()
		fakeCC.Add().Info(fakeLoginServer.URL())
		fakeLoginServer.RouteToHandler("POST", cf.PathCFAuth, RespondWithJSONEncoded(http.StatusOK, cf.Tokens{
			AccessToken: "test-access-token",
			ExpiresIn:   12000,
		}))
		logger = lager.NewLogger("cf")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
		setCfcClient(0)
	})

	AfterEach(func() {
		if fakeCC != nil {
			fakeCC.Close()
		}
		if fakeLoginServer != nil {
			fakeLoginServer.Close()
		}
	})

	Describe("GetServiceInstancesInOrg", func() {

		When("list service instances succeeds", func() {
			BeforeEach(func() {
				fakeCC.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/v3/service_instances", "organization_guids=some_guid&service_plan_guids=broker_service_plan_guid"),
						RespondWith(http.StatusCreated, LoadFile("service_instances.json"), http.Header{"Content-Type": []string{"application/json"}}),
					),
				)
			})

			It("returns correct struct", func() {
				numServices, err := cfc.GetServiceInstancesInOrg("some_guid", "broker_service_plan_guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(numServices).To(Equal(2))
			})
		})

		When("list service instances returns a 500 code", func() {
			BeforeEach(func() {
				fakeCC.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/v3/service_instances", "organization_guids=some_guid&service_plan_guids=broker_service_plan_guid"),
						RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError),
					),
				)
			})

			It("should return correct error", func() {
				_, err := cfc.GetServiceInstancesInOrg("some_guid", "broker_service_plan_guid")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(MatchRegexp(`failed GetServiceInstancesInOrg org\(some_guid\), servicePlan\(broker_service_plan_guid\): .*cf.Response\[.*cf.ServiceInstance\].*GET.*'UnknownError'.*`)))
			})
		})

	})

})
