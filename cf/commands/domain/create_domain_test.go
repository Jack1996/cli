package domain_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/commands/domain"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create domain command", func() {

	var (
		requirementsFactory *testreq.FakeReqFactory
		ui                  *testterm.FakeUI
		domainRepo          *testapi.FakeDomainRepository
		configRepo          configuration.ReadWriter
	)

	BeforeEach(func() {
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
		domainRepo = &testapi.FakeDomainRepository{}
		configRepo = testconfig.NewRepositoryWithAccessToken(configuration.TokenInfo{Username: "my-user"})
	})

	runCommand := func(args ...string) {
		ui = new(testterm.FakeUI)
		cmd := domain.NewCreateDomain(ui, configRepo, domainRepo)
		testcmd.RunCommand(cmd, args, requirementsFactory)
		return
	}

	It("fails with usage", func() {
		runCommand("")
		Expect(ui.FailedWithUsage).To(BeTrue())

		runCommand("org1")
		Expect(ui.FailedWithUsage).To(BeTrue())

		runCommand("org1", "example.com")
		Expect(ui.FailedWithUsage).To(BeFalse())
	})

	Context("checks login", func() {
		It("passes when logged in", func() {
			runCommand("my-org", "example.com")
			Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
			Expect(requirementsFactory.OrganizationName).To(Equal("my-org"))
		})

		It("fails when not logged in", func() {
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: false}

			runCommand("my-org", "example.com")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	It("creates a domain", func() {
		org := models.Organization{}
		org.Name = "myOrg"
		org.Guid = "myOrg-guid"
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, Organization: org}
		runCommand("myOrg", "example.com")

		Expect(domainRepo.CreateDomainName).To(Equal("example.com"))
		Expect(domainRepo.CreateDomainOwningOrgGuid).To(Equal("myOrg-guid"))
		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating domain", "example.com", "myOrg", "my-user"},
			[]string{"OK"},
		))
	})
})
