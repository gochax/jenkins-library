package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGctsRunUnitTestsForAllRepoPackagesCommand(t *testing.T) {

	testCmd := GctsRunUnitTestsForAllRepoPackagesCommand()

	// only high level testing performed - details are tested in step generation procudure
	assert.Equal(t, "gctsRunUnitTestsForAllRepoPackages", testCmd.Use, "command name incorrect")

}
