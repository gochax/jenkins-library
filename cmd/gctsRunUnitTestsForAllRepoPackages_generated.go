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

type gctsRunUnitTestsForAllRepoPackagesOptions struct {
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Repository string `json:"repository,omitempty"`
	Host       string `json:"host,omitempty"`
	Client     string `json:"client,omitempty"`
}

// GctsRunUnitTestsForAllRepoPackagesCommand Runs all existing unit tests for the repository packages
func GctsRunUnitTestsForAllRepoPackagesCommand() *cobra.Command {
	metadata := gctsRunUnitTestsForAllRepoPackagesMetadata()
	var stepConfig gctsRunUnitTestsForAllRepoPackagesOptions
	var startTime time.Time

	var createGctsRunUnitTestsForAllRepoPackagesCmd = &cobra.Command{
		Use:   "gctsRunUnitTestsForAllRepoPackages",
		Short: "Runs all existing unit tests for the repository packages",
		Long:  `This step will execute every unit test associated with a package belonging to the specified local repository on an ABAP system.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			startTime = time.Now()
			log.SetStepName("gctsRunUnitTestsForAllRepoPackages")
			log.SetVerbose(GeneralConfig.Verbose)
			err := PrepareConfig(cmd, &metadata, "gctsRunUnitTestsForAllRepoPackages", &stepConfig, config.OpenPiperFile)
			if err != nil {
				return err
			}
			log.RegisterSecret(stepConfig.Username)
			log.RegisterSecret(stepConfig.Password)
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
			telemetry.Initialize(GeneralConfig.NoTelemetry, "gctsRunUnitTestsForAllRepoPackages")
			gctsRunUnitTestsForAllRepoPackages(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
		},
	}

	addGctsRunUnitTestsForAllRepoPackagesFlags(createGctsRunUnitTestsForAllRepoPackagesCmd, &stepConfig)
	return createGctsRunUnitTestsForAllRepoPackagesCmd
}

func addGctsRunUnitTestsForAllRepoPackagesFlags(cmd *cobra.Command, stepConfig *gctsRunUnitTestsForAllRepoPackagesOptions) {
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password to authenticate to the ABAP system")
	cmd.Flags().StringVar(&stepConfig.Repository, "repository", os.Getenv("PIPER_repository"), "Specifies the name (ID) of the local repsitory on the ABAP system")
	cmd.Flags().StringVar(&stepConfig.Host, "host", os.Getenv("PIPER_host"), "Specifies the protocol and host adress, including the port. Please provide in the format '<protocol>://<host>:<port>'")
	cmd.Flags().StringVar(&stepConfig.Client, "client", os.Getenv("PIPER_client"), "Specifies the client of the ABAP system to be adressed")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("client")
}

// retrieve step metadata
func gctsRunUnitTestsForAllRepoPackagesMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:    "gctsRunUnitTestsForAllRepoPackages",
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
						Name:        "repository",
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
				},
			},
		},
	}
	return theMetaData
}
