package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

func gctsRunUnitTestsForAllRepoPackages(config gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData) error {

	client := piperhttp.Client{}
	cookieJar, _ := cookiejar.New(nil)
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	client.SetOptions(clientOptions)

	repoObjects, err := getPackageList(config, telemetryData, &client)
	if err != nil {
		log.Entry().Error(err)
	}

	url := "http://" + config.Host +
		"/sap/bc/adt/core/discovery?sap-client=" + config.Client

	header := make(http.Header)
	header.Add("Accept", "application/atomsvc+xml")
	header.Add("x-csrf-token", "fetch")
	header.Add("saml2", "disabled")

	disc, err := client.SendRequest("GET", url, nil, header, nil)
	if disc == nil {
		log.Entry().Fatal(err)
	}
	if disc != nil {
		if err != nil {
			log.Entry().WithError(err).
				Error("Discovery of ABAP server failed")
		}
	}
	header.Set("x-csrf-token", disc.Header.Get("X-Csrf-Token"))
	header.Set("Accept", "application/xml")
	header.Set("Content-Type", "application/vnd.sap.adt.abapunit.testruns.result.v1+xml")
	disc.Body.Close()

	for _, object := range repoObjects {
		err := executeTestsForPackage(config, telemetryData, &client, header, object)
		if err != nil {
			log.Entry().Fatalf("%v", err)
		}
	}
	log.Entry().
		WithField("repositoryName", config.RepositoryName).
		Info("All tests were successfull")

	return nil
}

func executeTestsForPackage(config gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, client piperhttp.Sender, header http.Header, packageName string) error {

	var xmlBody = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<aunit:runConfiguration xmlns:aunit="http://www.sap.com/adt/aunit">
		<external>
			<coverage active="false"/>
		</external>
		<options>
			<uriType value="semantic"/>
			<testDeterminationStrategy sameProgram="true" assignedTests="false" appendAssignedTestsPreview="true"/>
			<testRiskLevels harmless="true" dangerous="true" critical="true"/>
			<testDurations short="true" medium="true" long="true"/>
		</options>
		<adtcore:objectSets xmlns:adtcore="http://www.sap.com/adt/core">
			<objectSet kind="inclusive">
				<adtcore:objectReferences>
					<adtcore:objectReference adtcore:uri="/sap/bc/adt/packages/` +
		packageName + `"/>
				</adtcore:objectReferences>
			</objectSet>
		</adtcore:objectSets>
	</aunit:runConfiguration>`)

	url := "http://" + config.Host +
		"/sap/bc/adt/abapunit/testruns?sap-client=" + config.Client

	resp, HTTPErr := client.SendRequest("POST", url, bytes.NewBuffer(xmlBody), header, nil)
	if resp == nil {
		return fmt.Errorf("Request failed: %v", HTTPErr)
	}
	var response runResult
	if resp != nil {
		parsingErr := parseHTTPResponseBodyXML(resp, &response)
		if parsingErr != nil {
			return fmt.Errorf("%v", parsingErr)
		}
		if HTTPErr != nil {
			return fmt.Errorf("%v", HTTPErr)
		}
		resp.Body.Close()
	}

	aunitError := parseAUnitResponse(&response)
	if aunitError != nil {
		return fmt.Errorf("%v", aunitError)
	}

	return nil
}

func parseAUnitResponse(response *runResult) error {
	var node string
	aunitError := false

	for _, program := range response.Program {
		log.Entry().Infof("Testing class %v", program.Name)
		for _, testClass := range program.TestClasses.TestClass {
			log.Entry().Infof("With test class %v", testClass.Name)
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
		return fmt.Errorf("Some unit tests failed")
	}
	return nil
}

func getPackageList(config gctsRunUnitTestsForAllRepoPackagesOptions, telemetryData *telemetry.CustomData, client piperhttp.Sender) ([]string, error) {

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

	url := "http://" + config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.RepositoryName +
		"/getObjects?sap-client=" + config.Client

	resp, HTTPErr := client.SendRequest("GET", url, nil, nil, nil)
	if resp == nil {
		return nil, fmt.Errorf("Request failed: %v", HTTPErr)
	}
	var response objectsResponseBody
	if resp != nil {
		parsingErr := parseHTTPResponseBodyJSON(resp, &response)
		if parsingErr != nil {
			return nil, fmt.Errorf("%v", parsingErr)
		}
		if HTTPErr != nil {
			return nil, fmt.Errorf("%v", HTTPErr)
		}
		resp.Body.Close()
	}

	var repoObjects []string
	for _, object := range response.Objects {
		if object.Type == "DEVC" {
			repoObjects = append(repoObjects, object.Object)
		}
	}
	return repoObjects, nil
}

func parseHTTPResponseBodyXML(resp *http.Response, response interface{}) error {
	if resp == nil {
		return fmt.Errorf("http response was nil")
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read HTTP response body: %w", err)
	}
	xml.Unmarshal(bodyText, &response)

	return nil
}

// func parseHTTP(resp *http.Response, err error, errMessage string){

// }

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
