package cmd

import (
	"net/http/cookiejar"

	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsDeployCommit(config gctsDeployCommitOptions, telemetryData *telemetry.CustomData) error {

	client := piperhttp.Client{}
	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	client.SetOptions(clientOptions)

	type deployResponseBody struct {
		Trkorr     string    `json:"trkorr"`
		FromCommit string    `json:"fromCommit"`
		ToCommit   string    `json:"toCommit"`
		Log        []logs    `json:"log"`
		Exception  exception `json:"exception"`
		ErrorLogs  []logs    `json:"errorLog"`
	}

	url := "http://" + config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.RepositoryName +
		"/pullByCommit?sap-client=" + config.Client
	if config.Commit != "" {
		url = url + "&request=" + config.Commit
	}

	resp, err := client.SendRequest("GET", url, nil, nil, nil)
	if resp == nil {
		log.Entry().Fatal(err)
	}
	var response deployResponseBody
	if resp != nil {
		parsingErr := parseHTTPResponseBodyJSON(resp, &response)
		if parsingErr != nil {
			log.Entry().Warning(parsingErr)
		}
		if err != nil {
			log.Entry().WithError(err).
				WithField("repositoryName", config.RepositoryName).
				Fatalf("Failed to pull the latest commit: %v", response.Exception)
		}
		resp.Body.Close()
	}
	log.Entry().
		WithField("repositoryName", config.RepositoryName).
		Infof("Successfully pulled latest commit %v (previous commit was %v)", response.ToCommit, response.FromCommit)

	return nil
}
