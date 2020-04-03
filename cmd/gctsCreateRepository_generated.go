// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type gctsCreateRepositoryOptions struct {
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	RepositoryName string `json:"repositoryName,omitempty"`
	Host           string `json:"host,omitempty"`
	Client         string `json:"client,omitempty"`
	GithubURL      string `json:"githubURL,omitempty"`
	Role           string `json:"role,omitempty"`
	VSID           string `json:"vSID,omitempty"`
	Type           string `json:"type,omitempty"`
}

// GctsCreateRepositoryCommand Creates a Git repository
func GctsCreateRepositoryCommand() *cobra.Command {
	metadata := gctsCreateRepositoryMetadata()
	var stepConfig gctsCreateRepositoryOptions
	var startTime time.Time

	var createGctsCreateRepositoryCmd = &cobra.Command{
		Use:   "gctsCreateRepository",
		Short: "Creates a Git repository",
		Long:  `Creates a local Git repository if it does not already exist.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			startTime = time.Now()
			log.SetStepName("gctsCreateRepository")
			log.SetVerbose(GeneralConfig.Verbose)
			return PrepareConfig(cmd, &metadata, "gctsCreateRepository", &stepConfig, config.OpenPiperFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetry.Send(&telemetryData)
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, "gctsCreateRepository")
			gctsCreateRepository(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
		},
	}

	addGctsCreateRepositoryFlags(createGctsCreateRepositoryCmd, &stepConfig)
	return createGctsCreateRepositoryCmd
}

func addGctsCreateRepositoryFlags(cmd *cobra.Command, stepConfig *gctsCreateRepositoryOptions) {
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.RepositoryName, "repositoryName", os.Getenv("PIPER_repositoryName"), "Specifies the name (ID) of the repsitory to be cloned")
	cmd.Flags().StringVar(&stepConfig.Host, "host", os.Getenv("PIPER_host"), "Specifies the host address of the ABAP system including the port")
	cmd.Flags().StringVar(&stepConfig.Client, "client", os.Getenv("PIPER_client"), "Specifies the client of the ABAP system to be adressed")
	cmd.Flags().StringVar(&stepConfig.GithubURL, "githubURL", os.Getenv("PIPER_githubURL"), "URL of the corresponding remote repository")
	cmd.Flags().StringVar(&stepConfig.Role, "role", os.Getenv("PIPER_role"), "Role of the local repository. Choose between 'TARGET' and 'SOURCE'. Local repositories with a TARGET role will NOT be able to be the source of code changes.")
	cmd.Flags().StringVar(&stepConfig.VSID, "vSID", os.Getenv("PIPER_vSID"), "Virtual SID of the local repository. The vSID corresponds to the transport route that delivers content to the remote Git repository.")
	cmd.Flags().StringVar(&stepConfig.Type, "type", "Git", "Type of the used source code management tool. So far, only Git is supported.")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("repositoryName")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("client")
}

// retrieve step metadata
func gctsCreateRepositoryMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:    "gctsCreateRepository",
			Aliases: []config.Alias{},
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "username",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "password",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "repositoryName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "host",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "client",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "githubURL",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "role",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "vSID",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "type",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
				},
			},
		},
	}
	return theMetaData
}
