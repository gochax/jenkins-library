package cmd

import (
	// "io/ioutil"
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"strings"

	// gabs "github.com/Jeffail/gabs/v2"
	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/pkg/errors"
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
		return errors.Wrap(cookieErr, "rollback commit failed")
	}
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	successCommit, _ := getLastSuccessfullCommit(config, telemetryData, httpClient, clientOptions)

	fmt.Println(successCommit)
	// url := config.Host +
	// 	"/sap/bc/cts_abapvcs/repository/" + config.Repository +
	// 	"/getHistory?sap-client=" + config.Client

	// resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	// defer func() {
	// 	if resp != nil && resp.Body != nil {
	// 		resp.Body.Close()
	// 	}
	// }()

	// if resp == nil || httpErr != nil {
	// 	return errors.Wrap(httpErr, "rollback commit failed")
	// }

	// bodyText, readErr := ioutil.ReadAll(resp.Body)

	// if readErr != nil {
	// 	return errors.Wrapf(readErr, "deploying commit failed")
	// }

	// response, parsingErr := gabs.ParseJSON([]byte(bodyText))

	// if parsingErr != nil {
	// 	return errors.Wrap(parsingErr, "deploying commit failed")
	// }

	// var deployParams []string
	// if config.Commit != "" {
	// 	deployParams = []string{"gctsDeployCommit", "--username", config.Username, "--password", config.Password, "--host", config.Host, "--client", config.Client, "--repository", config.Repository, "--commit", config.Commit}
	// } else if fromCommit, ok := response.Path("result.0.fromCommit").Data().(string); ok {
	// 	deployParams = []string{"gctsDeployCommit", "--username", config.Username, "--password", config.Password, "--host", config.Host, "--client", config.Client, "--repository", config.Repository, "--commit", fromCommit}
	// } else {
	// 	return errors.Errorf("no commit to rollback to identified")
	// }

	// deployErr := command.RunExecutable("./piper", deployParams...)

	// if deployErr != nil {
	// 	return errors.Wrap(deployErr, "rollback commit failed")
	// }

	log.Entry().
		WithField("repository", config.Repository).
		Infof("rollback was successfull")
	return nil
}

func getLastSuccessfullCommit(config *gctsRollbackCommitOptions, telemetryData *telemetry.CustomData, httpClient piperhttp.Sender, clientOptions piperhttp.ClientOptions) (string, error) {

	// commitList, _ := getCommits(config, telemetryData, httpClient)

	remoteURL, err := getRepoInfo(config, telemetryData, httpClient)
	if err != nil {
		return "", err
	}
	fmt.Println(remoteURL)

	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return "", err
	}
	resources := strings.Split(parsedURL.Path, "/")

	fmt.Println(resources)
	fmt.Println(parsedURL.Host)

	clientOptions.Token = "3a09064f3029f5a304d25532ef8f95d1dfa6da44"
	fmt.Println(clientOptions)
	httpClient.SetOptions(clientOptions)

	url := parsedURL.Host + "api/v3/repos" + parsedURL.Path + "/commits?sap-client=" + config.Client

	for _, commit := range commitList {
		resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

		defer func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}()

		// TODO anpassen
		if resp == nil || httpErr != nil {
			return nil, errors.Errorf("fail: %v", httpErr)
		}
	}

	return "", nil
}

func getCommits(config *gctsRollbackCommitOptions, telemetryData *telemetry.CustomData, httpClient piperhttp.Sender) ([]string, error) {

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/getCommit?sap-client=" + config.Client

	type commitsResponseBody struct {
		Commits []struct {
			ID          string `json:"id"`
			Author      string `json:"author"`
			AuthorMail  string `json:"authorMail"`
			Message     string `json:"message"`
			Description string `json:"description"`
			Date        string `json:"date"`
		} `json:"commits"`
	}

	resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// TODO anpassen
	if resp == nil || httpErr != nil {
		return nil, errors.Errorf("fail: %v", httpErr)
	}

	var response commitsResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)
	if parsingErr != nil {
		return []string{}, errors.Errorf("%v", parsingErr)
	}

	commitList := []string{}
	for _, commit := range response.Commits {
		commitList = append(commitList, commit.ID)
	}

	return commitList, nil
}

func getRemoteRepoURL(config *gctsRollbackCommitOptions, telemetryData *telemetry.CustomData, httpClient piperhttp.Sender) (string, error) {

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"?sap-client=" + config.Client

	type getRepoResponseBody struct {
		Result struct {
			Rid           string `json:"rid"`
			Name          string `json:"name"`
			Role          string `json:"role"`
			Type          string `json:"type"`
			Vsid          string `json:"vsid"`
			Status        string `json:"status"`
			Branch        string `json:"branch"`
			URL           string `json:"url"`
			Version       string `json:"version"`
			Objects       int    `json:"objects"`
			CurrentCommit string `json:"currentCommit"`
			Connection    string `json:"connection"`
			Config        []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"config"`
		} `json:"result"`
	}

	resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// TODO anpassen
	if resp == nil || httpErr != nil {
		return "", errors.Errorf("fail: %v", httpErr)
	}

	var response getRepoResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)
	if parsingErr != nil {
		return "", errors.Errorf("%v", parsingErr)
	}

	if response.Result.URL == "" {
		return "", errors.Errorf("no remote repository URL was configured")
	}

	return response.Result.URL, nil
}
