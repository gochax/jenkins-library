package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsRunUnitTestsForAllRepoPackages(config gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData) {
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
	err := runUnitTestsForAllRepoPackages(&config, telemetryData, &c, httpClient)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func runUnitTestsForAllRepoPackages(config *gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, command execRunner, httpClient piperhttp.Sender) error {

	cookieJar, cookieErr := cookiejar.New(nil)
	if cookieErr != nil {
		return fmt.Errorf("execution of unit tests failed: %w", cookieErr)
	}
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	repoObjects, getPackageErr := getPackageList(config, telemetryData, httpClient)

	if getPackageErr != nil {
		return fmt.Errorf("%w", getPackageErr)
	}

	discHeader, discError := discoverServer(config, telemetryData, httpClient)

	if discError != nil {
		return fmt.Errorf("%w", discError)
	}

	if discHeader.Get("X-Csrf-Token") == "" {
		return fmt.Errorf("could not retrieve x-csrf-token from server")
	}

	header := make(http.Header)
	header.Add("x-csrf-token", discHeader.Get("X-Csrf-Token"))
	header.Add("Accept", "application/xml")
	header.Add("Content-Type", "application/vnd.sap.adt.abapunit.testruns.result.v1+xml")

	for _, object := range repoObjects {
		executeTestsErr := executeTestsForPackage(config, telemetryData, httpClient, header, object)

		if executeTestsErr != nil {
			return fmt.Errorf("%w", executeTestsErr)
		}
	}

	log.Entry().
		WithField("repository", config.Repository).
		Info("all unit tests were successfull")
	return nil
}

func discoverServer(config *gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, client piperhttp.Sender) (*http.Header, error) {

	url := config.Host +
		"/sap/bc/adt/core/discovery?sap-client=" + config.Client

	header := make(http.Header)
	header.Add("Accept", "application/atomsvc+xml")
	header.Add("x-csrf-token", "fetch")
	header.Add("saml2", "disabled")

	disc, httpErr := client.SendRequest("GET", url, nil, header, nil)

	defer func() {
		if disc != nil && disc.Body != nil {
			disc.Body.Close()
		}
	}()

	if disc == nil || disc.Header == nil || httpErr != nil {
		if httpErr != nil {
			return nil, fmt.Errorf("discovery of the ABAP server failed: %v", httpErr)
		}
		return nil, fmt.Errorf("discovery of the ABAP server failed: http response or header are <nil>")
	}

	return &disc.Header, nil
}

func executeTestsForPackage(config *gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, client piperhttp.Sender, header http.Header, packageName string) error {

	var xmlBody = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<aunit:runConfiguration
			xmlns:aunit="http://www.sap.com/adt/aunit">
			<external>
					<coverage active="false"/>
			</external>
			<options>
					<uriType value="semantic"/>
					<testDeterminationStrategy sameProgram="true" assignedTests="false" appendAssignedTestsPreview="true"/>
					<testRiskLevels harmless="true" dangerous="true" critical="true"/>
					<testDurations short="true" medium="true" long="true"/>
			</options>
			<adtcore:objectSets
					xmlns:adtcore="http://www.sap.com/adt/core">
					<objectSet kind="inclusive">
							<adtcore:objectReferences>
									<adtcore:objectReference adtcore:uri="/sap/bc/adt/packages/` + packageName + `"/>
							</adtcore:objectReferences>
					</objectSet>
			</adtcore:objectSets>
	</aunit:runConfiguration>`)

	url := config.Host +
		"/sap/bc/adt/abapunit/testruns?sap-client=" + config.Client

	resp, httpErr := client.SendRequest("POST", url, bytes.NewBuffer(xmlBody), header, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil || httpErr != nil {
		return fmt.Errorf("execution of unit tests failed: %w", httpErr)
	}

	var response runResult
	parsingErr := parseHTTPResponseBodyXML(resp, &response)
	if parsingErr != nil {
		log.Entry().Warning(parsingErr)
	}

	aunitError := parseAUnitResponse(&response)
	if aunitError != nil {
		return fmt.Errorf("%w", aunitError)
	}

	return nil
}

func parseAUnitResponse(response *runResult) error {
	var node string
	aunitError := false

	for _, program := range response.Program {
		log.Entry().Infof("testing class %v", program.Name)
		for _, testClass := range program.TestClasses.TestClass {
			log.Entry().Infof("using test class %v", testClass.Name)
			for _, testMethod := range testClass.TestMethods.TestMethod {
				node = testMethod.Name
				if len(testMethod.Alerts.Alert) > 0 {
					log.Entry().Errorf("%v - error", node)
					aunitError = true
				} else {
					log.Entry().Infof("%v - ok", node)
				}
			}
		}
	}
	if aunitError {
		return fmt.Errorf("some unit tests failed")
	}
	return nil
}

func getPackageList(config *gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, client piperhttp.Sender) ([]string, error) {

	type object struct {
		Pgmid       string `json:"pgmid"`
		Object      string `json:"object"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}

	type objectsResponseBody struct {
		Objects   []object  `json:"objects"`
		Log       []logs    `json:"log"`
		Exception exception `json:"exception"`
		ErrorLogs []logs    `json:"errorLog"`
	}

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/getObjects?sap-client=" + config.Client

	resp, httpErr := client.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil || httpErr != nil {
		return []string{}, fmt.Errorf("failed to get repository objects: %v", httpErr)
	}

	var response objectsResponseBody
	parsingErr := parseHTTPResponseBodyJSON(resp, &response)
	if parsingErr != nil {
		return []string{}, fmt.Errorf("%v", parsingErr)
	}

	repoObjects := []string{}
	for _, object := range response.Objects {
		if object.Type == "DEVC" {
			repoObjects = append(repoObjects, object.Object)
		}
	}

	return repoObjects, nil
}

func parseHTTPResponseBodyXML(resp *http.Response, response interface{}) error {
	if resp == nil {
		return fmt.Errorf("cannot parse HTTP response with value <nil>")
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read HTTP response body: %w", err)
	}
	xml.Unmarshal(bodyText, &response)

	return nil
}

type runResult struct {
	XMLName xml.Name `xml:"runResult"`
	Text    string   `xml:",chardata"`
	Aunit   string   `xml:"aunit,attr"`
	Program []struct {
		Text        string `xml:",chardata"`
		URI         string `xml:"uri,attr"`
		Type        string `xml:"type,attr"`
		Name        string `xml:"name,attr"`
		URIType     string `xml:"uriType,attr"`
		Adtcore     string `xml:"adtcore,attr"`
		TestClasses struct {
			Text      string `xml:",chardata"`
			TestClass []struct {
				Text             string `xml:",chardata"`
				URI              string `xml:"uri,attr"`
				Type             string `xml:"type,attr"`
				Name             string `xml:"name,attr"`
				URIType          string `xml:"uriType,attr"`
				NavigationURI    string `xml:"navigationUri,attr"`
				DurationCategory string `xml:"durationCategory,attr"`
				RiskLevel        string `xml:"riskLevel,attr"`
				TestMethods      struct {
					Text       string `xml:",chardata"`
					TestMethod []struct {
						Text          string `xml:",chardata"`
						URI           string `xml:"uri,attr"`
						Type          string `xml:"type,attr"`
						Name          string `xml:"name,attr"`
						ExecutionTime string `xml:"executionTime,attr"`
						URIType       string `xml:"uriType,attr"`
						NavigationURI string `xml:"navigationUri,attr"`
						Unit          string `xml:"unit,attr"`
						Alerts        struct {
							Text  string `xml:",chardata"`
							Alert []struct {
								Text     string `xml:",chardata"`
								Kind     string `xml:"kind,attr"`
								Severity string `xml:"severity,attr"`
								Title    string `xml:"title"`
								Details  struct {
									Text   string `xml:",chardata"`
									Detail []struct {
										Text     string `xml:",chardata"`
										AttrText string `xml:"text,attr"`
										Details  struct {
											Text   string `xml:",chardata"`
											Detail []struct {
												Text     string `xml:",chardata"`
												AttrText string `xml:"text,attr"`
											} `xml:"detail"`
										} `xml:"details"`
									} `xml:"detail"`
								} `xml:"details"`
								Stack struct {
									Text       string `xml:",chardata"`
									StackEntry struct {
										Text        string `xml:",chardata"`
										URI         string `xml:"uri,attr"`
										Type        string `xml:"type,attr"`
										Name        string `xml:"name,attr"`
										Description string `xml:"description,attr"`
									} `xml:"stackEntry"`
								} `xml:"stack"`
							} `xml:"alert"`
						} `xml:"alerts"`
					} `xml:"testMethod"`
				} `xml:"testMethods"`
			} `xml:"testClass"`
		} `xml:"testClasses"`
	} `xml:"program"`
}
