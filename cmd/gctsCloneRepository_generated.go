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

type gctsCloneRepositoryOptions struct {
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	RepositoryName string `json:"repositoryName,omitempty"`
	Host           string `json:"host,omitempty"`
	Client         string `json:"client,omitempty"`
}

// GctsCloneRepositoryCommand Clones a Git repository
func GctsCloneRepositoryCommand() *cobra.Command {
	metadata := gctsCloneRepositoryMetadata()
	var stepConfig gctsCloneRepositoryOptions
	var startTime time.Time

	var createGctsCloneRepositoryCmd = &cobra.Command{
		Use:   "gctsCloneRepository",
		Short: "Clones a Git repository",
		Long:  `Clones a Git repository from a remote repository to a local repository. To be able to execute this step, the corresponding local repository has to exist on the local ABAP system.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			startTime = time.Now()
			log.SetStepName("gctsCloneRepository")
			log.SetVerbose(GeneralConfig.Verbose)
			err := PrepareConfig(cmd, &metadata, "gctsCloneRepository", &stepConfig, config.OpenPiperFile)
			if err != nil {
				return err
			}
			return nil
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
			telemetry.Initialize(GeneralConfig.NoTelemetry, "gctsCloneRepository")
			gctsCloneRepository(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
		},
	}

	addGctsCloneRepositoryFlags(createGctsCloneRepositoryCmd, &stepConfig)
	return createGctsCloneRepositoryCmd
}

func addGctsCloneRepositoryFlags(cmd *cobra.Command, stepConfig *gctsCloneRepositoryOptions) {
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.RepositoryName, "repositoryName", os.Getenv("PIPER_repositoryName"), "Specifies the name (ID) of the repsitory to be cloned")
	cmd.Flags().StringVar(&stepConfig.Host, "host", os.Getenv("PIPER_host"), "Specifies the host address of the ABAP system including the port")
	cmd.Flags().StringVar(&stepConfig.Client, "client", os.Getenv("PIPER_client"), "Specifies the client of the ABAP system to be adressed")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("repositoryName")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("client")
}

// retrieve step metadata
func gctsCloneRepositoryMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:    "gctsCloneRepository",
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
						Aliases:     []config.Alias{{Name: "username"}},
					},
					{
						Name:        "password",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "password"}},
					},
					{
						Name:        "repositoryName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "repositoryName"}},
					},
					{
						Name:        "host",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "host"}},
					},
					{
						Name:        "client",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "client"}},
					},
				},
			},
		},
	}
	return theMetaData
}
