package cmd

import (
	"net/http"
	"net/http/cookiejar"

	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsCloneRepository(config gctsCloneRepositoryOptions, telemetryData *telemetry.CustomData) error {

	client := piperhttp.Client{}
	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	client.SetOptions(clientOptions)

	type cloneResultBody struct {
		RID          string `json:"rid"`
		CheckoutTime string `json:"checkoutTime"`
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

	url := "http://" + config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.RepositoryName +
		"/clone?sap-client=" + config.Client

	resp, err := client.SendRequest("POST", url, nil, header, nil)

	if resp == nil {
		log.Entry().Fatal(err)
	}
	var response cloneResponseBody
	if resp != nil {
		parsingErr := parseHTTPResponseBodyJSON(resp, &response)
		if parsingErr != nil {
			log.Entry().Warning(parsingErr)
		}
		if err != nil {
			if resp.StatusCode == 500 && response.ErrorLogs[1].Code == "GCTS.CLIENT.1420" {
				log.Entry().
					WithField("repositoryName", config.RepositoryName).
					Info("The repository has already been cloned")
				return nil
			}
			log.Entry().WithError(err).
				WithField("repositoryName", config.RepositoryName).
				Fatalf("Cloning the repository failed %v", response.Exception)
		}
		resp.Body.Close()
	}
	log.Entry().
		WithField("repositoryName", config.RepositoryName).
		Info("Successfully cloned the Git repository to the local repository")

	// defer resp.Body.Close()

	// var response cloneResponseBody
	// bodyText, readErr := ioutil.ReadAll(resp.Body)
	// // TODO learn how to respond to errors correctly
	// if readErr != nil {
	// 	log.Entry().
	// 		Warning("Could not read HTTP response")
	// }
	// json.Unmarshal(bodyText, &response)

	// if cloneErr != nil {
	// 	if resp.StatusCode == 500 {
	// 		if response.ErrorLogs[1].Severity == "ERROR" &&
	// 			response.ErrorLogs[1].Code == "GCTS.CLIENT.1420" {
	// 			log.Entry().
	// 				WithField("repositoryName", config.RepositoryName).
	// 				Info("The repository has already been cloned")
	// 			return nil
	// 		}
	// 		log.Entry().WithError(cloneErr).
	// 			WithField("repositoryName", config.RepositoryName).
	// 			Fatal("Cloning the repository failed")
	// 	}
	// 	log.Entry().WithError(cloneErr).
	// 		Warning("We got an error back, but no idea what kind")
	// }
	// log.Entry().
	// 	WithField("repositoryName", config.RepositoryName).
	// 	Info("Successfully cloned the requested Git repository")

	return nil
}

type exception struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type logs struct {
	Time     string `json:"time"`
	User     string `json:"user"`
	Section  string `json:"section"`
	Action   string `json:"action"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Code     string `json:"code"`
}
