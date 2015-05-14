// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"strings"

	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/provision"
	"github.com/tsuru/tsuru/repository"
)

// gitDeployCmds returns the list of commands that are used when the
// provisioner deploys a unit using the Git repository method.
func gitDeployCmds(app provision.App, version string) ([]string, error) {
	repo, err := repository.Manager().GetRepository(app.GetName())
	if err != nil {
		return nil, err
	}
	return deployCmds(app, "git", repo.ReadOnlyURL, version)
}

// archiveDeployCmds returns the list of commands that are used when the
// provisioner deploys a unit using the archive method.
func archiveDeployCmds(app provision.App, archiveURL string) ([]string, error) {
	return deployCmds(app, "archive", archiveURL)
}

func deployCmds(app provision.App, params ...string) ([]string, error) {
	deployCmd, err := config.GetString("docker:deploy-cmd")
	if err != nil {
		return nil, err
	}
	cmds := append([]string{deployCmd}, params...)
	host := app.Envs()["TSURU_HOST"].Value
	token := app.Envs()["TSURU_APP_TOKEN"].Value
	unitAgentCmds := []string{"tsuru_unit_agent", host, token, app.GetName(), `"` + strings.Join(cmds, " ") + `"`, "deploy"}
	finalCmd := strings.Join(unitAgentCmds, " ")
	return []string{"/bin/bash", "-lc", finalCmd}, nil
}
