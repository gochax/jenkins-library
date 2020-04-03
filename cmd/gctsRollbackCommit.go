package cmd

import (
	// "fmt"
	"net/http/cookiejar"

	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsRollbackCommit(config gctsRollbackCommitOptions, telemetryData *telemetry.CustomData) error {

	client := piperhttp.Client{}
	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	client.SetOptions(clientOptions)

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

	url := "http://" + config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.RepositoryName +
		"/getHistory?sap-client=" + config.Client

	resp, err := client.SendRequest("GET", url, nil, nil, nil)
	if resp == nil {
		log.Entry().Fatalf("Request failed: %v", err)
	}
	var response rollbackResponseBody
	if resp != nil {
		parsingErr := parseHTTPResponseBodyJSON(resp, &response)
		if parsingErr != nil {
			log.Entry().Warning(parsingErr)
		}
		if err != nil {
			log.Entry().WithError(err).WithField("StatusCode", resp.Status).Fatalf("Could not get repository commit history %v", response.Exception)
		}
		c := command.Command{}
		deployParams := []string{"gctsDeployCommit", "--username", config.Username, "--password", config.Password, "--host", config.Host, "--client", config.Client, "--repositoryName", config.RepositoryName, "--commit", response.Result[0].FromCommit}
		deployErr := c.RunExecutable("./piper", deployParams...)
		// TODO decide whats an ERROR and whats FATAL
		if deployErr != nil {
			log.Entry().WithError(deployErr).Fatalf("Failed to deploy commit %v", response.Result[0].FromCommit)
		}
		resp.Body.Close()
	}
	log.Entry().
		WithField("repositoryName", config.RepositoryName).
		Infof("Rollback was successfull")

	return nil
}
