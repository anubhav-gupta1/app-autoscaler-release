package pre_upgrade_test

import (
	"acceptance/config"
	"acceptance/helpers"
	"fmt"
	"testing"
	"time"

	"github.com/KevinJCross/cf-test-helpers/v2/workflowhelpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	cfg   *config.Config
	setup *workflowhelpers.ReproducibleTestSuiteSetup
	nodeAppDropletPath string
)

func TestSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	cfg = config.LoadConfig(t)
	setup = workflowhelpers.NewTestSuiteSetup(cfg)

	RunSpecs(t, "Pre Upgrade Test Suite")
}

var _ = BeforeSuite(func() {

	fmt.Println("Clearing down existing test orgs/spaces...")
	setup = workflowhelpers.NewTestSuiteSetup(cfg)

	workflowhelpers.AsUser(setup.AdminUserContext(), cfg.DefaultTimeoutDuration(), func() {
		orgs := helpers.GetTestOrgs(cfg)

		for _, org := range orgs {
			orgName, _, spaceName, _ := helpers.GetOrgSpaceNamesAndGuids(cfg, org)
			if spaceName != "" {
				helpers.DeleteOrgWithTimeout(cfg, orgName, time.Duration(60)*time.Second)
			}
		}
	})

	fmt.Println("Clearing down existing test orgs/spaces... Complete")
	setup.Setup()

	workflowhelpers.AsUser(setup.AdminUserContext(), cfg.DefaultTimeoutDuration(), func() {
		if cfg.ShouldEnableServiceAccess() {
			helpers.EnableServiceAccess(cfg, setup.GetOrganizationName())
		}
	})

	if cfg.IsServiceOfferingEnabled() {
		helpers.CheckServiceExists(cfg, setup.TestSpace.SpaceName(), cfg.ServiceName)
	}

	fmt.Println("creating droplet")
	nodeAppDropletPath = helpers.CreateDroplet(*cfg)
})
