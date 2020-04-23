package cmd

import (
	"fmt"
	"net/http/cookiejar"

	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsDeployCommit(config gctsDeployCommitOptions, telemetryData *telemetry.CustomData) {
	// for command execution use Command
	c := command.Command{}
	// reroute command output to logging framework
	c.Stdout(log.Entry().Writer())
	c.Stderr(log.Entry().Writer())

	// for http calls import  piperhttp "github.com/SAP/jenkins-library/pkg/http"
	// and use a  &piperhttp.Client{} in a custom system
	// Example: step checkmarxExecuteScan.go
	httpClient := &piperhttp.Client{}

	// error situations should stop execution through log.Entry().Fatal() call which leads to an os.Exit(1) in the end
	err := deployCommit(&config, telemetryData, &c, httpClient)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func deployCommit(config *gctsDeployCommitOptions, telemetryData *telemetry.CustomData, command execRunner, httpClient piperhttp.Sender) error {

	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	type deployResponseBody struct {
		Trkorr     string    `json:"trkorr"`
		FromCommit string    `json:"fromCommit"`
		ToCommit   string    `json:"toCommit"`
		Log        []logs    `json:"log"`
		Exception  exception `json:"exception"`
		ErrorLogs  []logs    `json:"errorLog"`
	}

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/pullByCommit?sap-client=" + config.Client
	if config.Commit != "" {
		url = url + "&request=" + config.Commit
	}

	resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil || httpErr != nil {
		return fmt.Errorf("deploy commit failed: %w", httpErr)
	}

	var response deployResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)

	if parsingErr != nil {
		log.Entry().Warning(parsingErr)
	}

	log.Entry().
		WithField("repository", config.Repository).
		Infof("successfully deployed commit %v (previous commit was %v)", response.ToCommit, response.FromCommit)
	return nil
}
