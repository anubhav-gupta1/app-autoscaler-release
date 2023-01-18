package run_performance_test

import (
	"acceptance/config"
	. "acceptance/helpers"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/cloudfoundry/cf-test-helpers/v2/workflowhelpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	cfg            *config.Config
	setup          *workflowhelpers.ReproducibleTestSuiteSetup
	orgName        string
	spaceName      string
	errorsScaleOut sync.Map
	errorsScaleIn  sync.Map
)

func TestSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	cfg = config.LoadConfig()
	cfg.Prefix = "autoscaler-performance"
	setup = workflowhelpers.NewTestSuiteSetup(cfg)
	RunSpecs(t, "Performance Test Suite")
}

var _ = BeforeSuite(func() {
	// use smoke test to avoid creating a new user
	setup = workflowhelpers.NewSmokeTestSuiteSetup(cfg)

	if cfg.UseExistingOrganization && !cfg.UseExistingSpace {
		orgGuid := GetOrgGuid(cfg, cfg.ExistingOrganization)
		spaces := GetTestSpaces(orgGuid, cfg)
		Expect(len(spaces)).To(Equal(1), "Found more than one space in existing org %s", cfg.ExistingOrganization)
		cfg.ExistingSpace = spaces[0]
	} else {
		workflowhelpers.AsUser(setup.AdminUserContext(), cfg.DefaultTimeoutDuration(), func() {
			orgName, spaceName = FindExistingOrgAndSpace(cfg)
		})

		Expect(orgName).ToNot(Equal(""), "orgName has not been determined")
		Expect(spaceName).ToNot(Equal(""), "spaceName has not been determined")

		cfg.ExistingOrganization = orgName
		cfg.ExistingSpace = spaceName
	}

	cfg.UseExistingOrganization = true
	cfg.UseExistingSpace = true

	setup = workflowhelpers.NewTestSuiteSetup(cfg)

	setup.Setup()

	if cfg.IsServiceOfferingEnabled() {
		CheckServiceExists(cfg, setup.TestSpace.SpaceName(), cfg.ServiceName)
	}
})

var _ = AfterSuite(func() {

	errorsScaleOut.Range(func(appName, err interface{}) bool {
		fmt.Printf("scale-out app error: %s: %s \n", appName, err.(error).Error())
		return true
	})

	errorsScaleOut.Range(func(appName, err interface{}) bool {
		fmt.Printf("scale-in app error: %s: %s \n", appName, err.(error).Error())
		return true
	})

	if os.Getenv("SKIP_TEARDOWN") == "true" {
		fmt.Println("Skipping Teardown...")
	} else {
		fmt.Println("TODO: Cleanup test...")
		setup.Teardown()
	}
})
