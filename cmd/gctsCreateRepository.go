package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsCreateRepository(config gctsCreateRepositoryOptions, telemetryData *telemetry.CustomData) error {

	client := piperhttp.Client{}
	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	client.SetOptions(clientOptions)

	// var response createResponseBody
	// err := getRepoInfo(config, &client, &response)
	// if err != nil {
	// 	log.Entry().Warning(err)
	// }

	type repoData struct {
		RID             string `json:"rid"`
		Name            string `json:"name"`
		Role            string `json:"role"`
		Type            string `json:"type"`
		VSID            string `json:"vsid"`
		GithubURLstring string `json:"url"`
	}

	type createRequestBody struct {
		Repository string   `json:"repository"`
		Data       repoData `json:"data"`
	}

	type repoConfig struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Category string `json:"category"`
	}

	type createResultBody struct {
		RID         string       `json:"rid"`
		Name        string       `json:"name"`
		Role        string       `json:"role"`
		Type        string       `json:"type"`
		VSID        string       `json:"vsid"`
		Status      string       `json:"status"`
		Branch      string       `json:"branch"`
		URL         string       `json:"url"`
		CreatedBy   string       `json:"createdBy"`
		CreatedDate string       `json:"createdDate"`
		Connection  string       `json:"connection"`
		Config      []repoConfig `json:"config"`
	}

	type createResponseBody struct {
		Repository createResultBody `json:"repository"`
		Exception  string           `json:"exception"`
	}

	reqBody := createRequestBody{
		Repository: config.RepositoryName,
		Data: repoData{
			RID:             config.RepositoryName,
			Name:            config.RepositoryName,
			Role:            config.Role,
			Type:            config.Type,
			VSID:            config.VSID,
			GithubURLstring: config.GithubURL,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Add("Accept", "application/json")

	url := "http://" + config.Host +
		"/sap/bc/cts_abapvcs/repository?sap-client=" + config.Client

	resp, err := client.SendRequest("POST", url, bytes.NewBuffer(jsonBody), header, nil)
	if resp == nil {
		log.Entry().Fatal(err)
	}
	var response createResponseBody
	if resp != nil {
		parsingErr := parseHTTPResponseBodyJSON(resp, &response)
		if parsingErr != nil {
			log.Entry().Warning(parsingErr)
		}
		if err != nil {
			if resp.StatusCode == 500 && response.Exception == "Repository already exists" {
				log.Entry().
					WithField("repositoryName", config.RepositoryName).
					Info("The repository already exists locally")
				return nil
			}
			log.Entry().WithError(err).
				WithField("repositoryName", config.RepositoryName).
				Fatalf("Creating the repository failed: %v", response.Exception)
		}
		resp.Body.Close()
	}

	log.Entry().
		WithField("repositoryName", config.RepositoryName).
		Info("Successfully created the local repository")

	return nil
}

// func checkIfRepoExists(options gctsCreateRepositoryOptions, client piperhttp.Sender, response interface{}) (bool, error) {

// }

// func getRepoInfo(config gctsCreateRepositoryOptions, client piperhttp.Sender, response interface{}) error {

// 	url := "http://" +
// 		config.Host + "/sap/bc/cts_abapvcs/repository/" +
// 		config.RepositoryName +
// 		"?sap-client=" +
// 		config.Client

// 	resp, err := client.SendRequest("GET", url, nil, nil, nil)
// 	if err != nil {
// 		return fmt.Errorf("Request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	err = parseHTTPResponseBodyJSON(resp, &response)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func parseHTTPResponseBodyJSON(resp *http.Response, response interface{}) error {
	if resp == nil {
		return fmt.Errorf("http response was nil")
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read HTTP response body: %w", err)
	}
	json.Unmarshal(bodyText, &response)
	// fmt.Println(resp.Status)
	// fmt.Println(string(bodyText))
	// fmt.Println(response)

	return nil
}
