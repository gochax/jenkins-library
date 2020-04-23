package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsCloneRepository(config gctsCloneRepositoryOptions, telemetryData *telemetry.CustomData) {
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
	err := cloneRepository(&config, telemetryData, &c, httpClient)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func cloneRepository(config *gctsCloneRepositoryOptions, telemetryData *telemetry.CustomData, command execRunner, httpClient piperhttp.Sender) error {

	cookieJar, cookieErr := cookiejar.New(nil)
	if cookieErr != nil {
		return fmt.Errorf("cloning the repository failed: %w", cookieErr)
	}
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	type cloneResultBody struct {
		RID          string `json:"rid"`
		CheckoutTime int    `json:"checkoutTime"`
		FromCommit   string `json:"fromCommit"`
		ToCommit     string `json:"toCommit"`
		Caller       string `json:"caller"`
		Request      string `json:"request"`
		Type         string `json:"type"`
	}

	type cloneResponseBody struct {
		Result    cloneResultBody `json:"result"`
		Log       []logs          `json:"log"`
		Exception exception       `json:"exception"`
		ErrorLogs []logs          `json:"errorLog"`
	}

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Add("Accept", "application/json")

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/clone?sap-client=" + config.Client

	resp, httpErr := httpClient.SendRequest("POST", url, nil, header, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil {
		return fmt.Errorf("cloning the repository failed: %w", httpErr)
	}

	var response cloneResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)

	if parsingErr != nil {
		return fmt.Errorf("creating repository on the ABAP system %v failed: %w", config.Host, parsingErr)
	}

	if httpErr != nil {
		if resp.StatusCode == 500 && response.ErrorLogs[1].Code == "GCTS.CLIENT.1420" {
			log.Entry().
				WithField("repository", config.Repository).
				Info("the repository has already been cloned")
			return nil
		}
		return fmt.Errorf("cloning the repository failed: %w", httpErr)
	}

	log.Entry().
		WithField("repository", config.Repository).
		Info("successfully cloned the Git repository to the local repository")
	return nil
}

type exception struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Code        int    `json:"code"`
}

type logs struct {
	Time     int    `json:"time"`
	User     string `json:"user"`
	Section  string `json:"section"`
	Action   string `json:"action"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Code     string `json:"code"`
}

func parseHTTPResponseBodyJSON(resp *http.Response, response interface{}) error {
	if resp == nil {
		return fmt.Errorf("cannot parse HTTP response with value <nil>")
	}

	bodyText, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("cannot read HTTP response body: %w", readErr)
	}

	marshalErr := json.Unmarshal(bodyText, &response)
	if marshalErr != nil {
		return fmt.Errorf("cannot parse HTTP response as JSON: %w", marshalErr)
	}

	return nil
}
