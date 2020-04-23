package cmd

import (
	"fmt"
	"net/http/cookiejar"

	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsRollbackCommit(config gctsRollbackCommitOptions, telemetryData *telemetry.CustomData) {
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
	err := rollbackCommit(&config, telemetryData, &c, httpClient)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func rollbackCommit(config *gctsRollbackCommitOptions, telemetryData *telemetry.CustomData, command execRunner, httpClient piperhttp.Sender) error {

	cookieJar, cookieErr := cookiejar.New(nil)
	if cookieErr != nil {
		return fmt.Errorf("rollback commit failed: %w", cookieErr)
	}
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	type rollbackResultBody struct {
		RID          string `json:"rid"`
		CheckoutTime int    `json:"checkoutTime"`
		FromCommit   string `json:"fromCommit"`
		ToCommit     string `json:"toCommit"`
		Caller       string `json:"caller"`
		Request      string `json:"request"`
		Type         string `json:"type"`
	}

	type rollbackResponseBody struct {
		Result    []rollbackResultBody `json:"result"`
		Log       []logs               `json:"log"`
		Exception exception            `json:"exception"`
		ErrorLogs []logs               `json:"errorLog"`
	}

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/getHistory?sap-client=" + config.Client

	resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil || httpErr != nil {
		return fmt.Errorf("rollback commit failed: %w", httpErr)
	}

	var response rollbackResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)

	if parsingErr != nil {
		log.Entry().Warning(parsingErr)
	}

	var deployParams []string
	if config.Commit != "" {
		deployParams = []string{"gctsDeployCommit", "--username", config.Username, "--password", config.Password, "--host", config.Host, "--client", config.Client, "--repository", config.Repository, "--commit", config.Commit}
	} else if response.Result[0].FromCommit != "" {
		deployParams = []string{"gctsDeployCommit", "--username", config.Username, "--password", config.Password, "--host", config.Host, "--client", config.Client, "--repository", config.Repository, "--commit", response.Result[0].FromCommit}
	} else {
		return fmt.Errorf("no commit to rollback to identified")
	}

	deployErr := command.RunExecutable("./piper", deployParams...)

	if deployErr != nil {
		return fmt.Errorf("rollback commit failed: %w", deployErr)
	}

	log.Entry().
		WithField("repository", config.Repository).
		Infof("rollback was successfull")
	return nil
}
