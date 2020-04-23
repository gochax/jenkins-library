package cmd

import (
	"testing"

	"github.com/SAP/jenkins-library/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestGctsRollbackCommitSuccess(t *testing.T) {

	execRunner := mock.ExecMockRunner{}
	config := gctsRollbackCommitOptions{
		Host:       "http://testHost.com:50000",
		Client:     "000",
		Repository: "testRepo",
		Username:   "testUser",
		Password:   "testPassword",
	}

	t.Run("rollback one commit by commit history", func(t *testing.T) {

		httpClient := httpMock{StatusCode: 200, ResponseBody: `{
			"result": [
				{
					"rid": "my-repository",
					"checkoutTime": 20180606130524,
					"fromCommit": "f1cdb6a032c1d8187c0990b51e94e8d8bb9898b2",
					"toCommit": "7763c593906d3cb92e264733b620df9c27f2e1b8",
					"caller": "JOHNDOE",
					"request": "SIDK1234567",
					"type": "PULL"
				}
			]
		}`}

		err := rollbackCommit(&config, nil, &execRunner, &httpClient)

		if assert.NoError(t, err) {

			t.Run("check url", func(t *testing.T) {
				assert.Equal(t, "http://testHost.com:50000/sap/bc/cts_abapvcs/repository/testRepo/getHistory?sap-client=000", httpClient.URL)
			})

			t.Run("check method", func(t *testing.T) {
				assert.Equal(t, "GET", httpClient.Method)
			})

			t.Run("check user", func(t *testing.T) {
				assert.Equal(t, "testUser", httpClient.Options.Username)
			})

			t.Run("check password", func(t *testing.T) {
				assert.Equal(t, "testPassword", httpClient.Options.Password)
			})

			t.Run("check CLI call", func(t *testing.T) {
				assert.Equal(t, "./piper", execRunner.Calls[0].Exec)
			})

			t.Run("check CLI call parameters", func(t *testing.T) {
				assert.Equal(t, []string{"gctsDeployCommit", "--username", "testUser", "--password", "testPassword", "--host", "http://testHost.com:50000", "--client", "000", "--repository", "testRepo", "--commit", "f1cdb6a032c1d8187c0990b51e94e8d8bb9898b2"}, execRunner.Calls[0].Params)
			})

		}
	})
}

func TestGctsRollbackToGivenCommitSuccess(t *testing.T) {

	execRunner := mock.ExecMockRunner{}
	config := gctsRollbackCommitOptions{
		Host:       "http://testHost.com:50000",
		Client:     "000",
		Repository: "testRepo",
		Username:   "testUser",
		Password:   "testPassword",
		Commit:     "8aeebd1a125d5d27499bd30699f4db2e79f51ee7",
	}

	t.Run("rollback to given commit", func(t *testing.T) {

		httpClient := httpMock{StatusCode: 200, ResponseBody: `{
				"result": [
					{
						"rid": "my-repository",
						"checkoutTime": 20180606130524,
						"fromCommit": "f1cdb6a032c1d8187c0990b51e94e8d8bb9898b2",
						"toCommit": "7763c593906d3cb92e264733b620df9c27f2e1b8",
						"caller": "JOHNDOE",
						"request": "SIDK1234567",
						"type": "PULL"
					}
				]
			}`}

		err := rollbackCommit(&config, nil, &execRunner, &httpClient)

		if assert.NoError(t, err) {

			t.Run("check CLI call", func(t *testing.T) {
				assert.Equal(t, "./piper", execRunner.Calls[0].Exec)
			})

			t.Run("check CLI call parameters", func(t *testing.T) {
				assert.Equal(t, []string{"gctsDeployCommit", "--username", "testUser", "--password", "testPassword", "--host", "http://testHost.com:50000", "--client", "000", "--repository", "testRepo", "--commit", "8aeebd1a125d5d27499bd30699f4db2e79f51ee7"}, execRunner.Calls[0].Params)
			})

		}
	})

}

func TestGctsRollbackCommitFailure(t *testing.T) {

	execRunner := mock.ExecMockRunner{}
	config := gctsRollbackCommitOptions{
		Host:       "http://testHost.com:50000",
		Client:     "000",
		Repository: "testRepo",
		Username:   "testUser",
		Password:   "testPassword",
	}

	t.Run("no commit history & no specified commit", func(t *testing.T) {

		httpClient := httpMock{StatusCode: 200, ResponseBody: `{
			"result": [
				{
					"rid": "my-repository",
					"checkoutTime": 20180606130524,
					"toCommit": "7763c593906d3cb92e264733b620df9c27f2e1b8",
					"caller": "JOHNDOE",
					"request": "SIDK1234567",
					"type": "PULL"
				}
			]
		}`}

		err := rollbackCommit(&config, nil, &execRunner, &httpClient)

		assert.EqualError(t, err, "no commit to rollback to identified")
	})

	t.Run("http error when getting commit history", func(t *testing.T) {

		httpClient := httpMock{StatusCode: 500, ResponseBody: `{
			"log": [
				{
					"time": 20180606130524,
					"user": "JENKINS",
					"section": "REPOSITORY_FACTORY",
					"action": "CREATE_REPOSITORY",
					"severity": "INFO",
					"message": "Start action CREATE_REPOSITORY review",
					"code": "GCTS.API.410"
				}
			],
			"errorLog": [
				{
					"time": 20180606130524,
					"user": "JENKINS",
					"section": "REPOSITORY_FACTORY",
					"action": "CREATE_REPOSITORY",
					"severity": "INFO",
					"message": "Start action CREATE_REPOSITORY review",
					"code": "GCTS.API.410"
				}
			],
			"exception": {
				"message": "repository_not_found",
				"description": "Repository not found",
				"code": 404
			}
		}`}

		err := rollbackCommit(&config, nil, &execRunner, &httpClient)

		assert.EqualError(t, err, "rollback commit failed: a http error occurred")
	})

}
