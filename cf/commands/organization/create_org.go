package organization

import (
	"github.com/cloudfoundry/cli/cf"
	"github.com/cloudfoundry/cli/cf/api/organizations"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	. "github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type CreateOrg struct {
	ui      terminal.UI
	config  configuration.Reader
	orgRepo organizations.OrganizationRepository
}

func NewCreateOrg(ui terminal.UI, config configuration.Reader, orgRepo organizations.OrganizationRepository) (cmd CreateOrg) {
	cmd.ui = ui
	cmd.config = config
	cmd.orgRepo = orgRepo
	return
}

func (cmd CreateOrg) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-org",
		ShortName:   "co",
		Description: T("Create an org"),
		Usage:       T("CF_NAME create-org ORG"),
	}
}

func (cmd CreateOrg) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd CreateOrg) Run(c *cli.Context) {
	name := c.Args()[0]

	cmd.ui.Say(T("Creating org {{.OrgName}} as {{.Username}}...",
		map[string]interface{}{
			"OrgName":  terminal.EntityNameColor(name),
			"Username": terminal.EntityNameColor(cmd.config.Username())}))
	err := cmd.orgRepo.Create(name)
	if err != nil {
		if apiErr, ok := err.(errors.HttpError); ok && apiErr.ErrorCode() == errors.ORG_EXISTS {
			cmd.ui.Ok()
			cmd.ui.Warn(T("Org {{.OrgName}} already exists",
				map[string]interface{}{"OrgName": name}))
			return
		} else {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.ui.Ok()
	cmd.ui.Say(T("\nTIP: Use '{{.Command}}' to target new org",
		map[string]interface{}{"Command": terminal.CommandColor(cf.Name() + " target -o " + name)}))
}
