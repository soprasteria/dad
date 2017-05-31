package export

import (
	"testing"

	"github.com/soprasteria/dad/server/types"
	"github.com/stretchr/testify/assert"
)

func TestBestIndicatorStatus_ZeroTechnicalServiceZeroIndicator_emptyResult(t *testing.T) {

	// Given
	services := []string{}
	servicesToMatch := []types.UsageIndicator{}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Empty", status)
}

func TestBestIndicatorStatus_ZeroTechnicalServiceOneIndicator_emptyResult(t *testing.T) {

	// Given
	services := []string{}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Empty", status)
}

func TestBestIndicatorStatus_OneTechnicalServiceZeroIndicator_emptyResult(t *testing.T) {

	// Given
	services := []string{"jenkins"}
	servicesToMatch := []types.UsageIndicator{}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Empty", status)
}

func TestBestIndicatorStatus_OneTechnicalServiceOneIndicator_EmptyResult(t *testing.T) {

	// Given
	services := []string{"jenkins"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Empty", status)
}

func TestBestIndicatorStatus_OneTechnicalServiceOneIndicator_UndeterminedResult(t *testing.T) {

	// Given
	services := []string{"jenkins"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Undetermined", status)
}

func TestBestIndicatorStatus_OneTechnicalServiceOneIndicator_InactiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Inactive", status)
}

func TestBestIndicatorStatus_OneTechnicalServiceOneIndicator_ActiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Active", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_EmptyVSEmpty_EmptyResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Empty"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Empty", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_EmptyVSUndetermined_UndeterminedResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Undetermined"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Undetermined", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_EmptyVSInactive_InactiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Inactive"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Inactive", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_EmptyVSActive_ActiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Active"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Active", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_UndeterminedVSUndetermined_UndeterminedResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Undetermined"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Undetermined", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_UndeterminedVSInactive_InactiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Inactive"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Inactive", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_UndeterminedVSActive_ActiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Active"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Active", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_InactiveVSInactive_InactiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Inactive"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Inactive", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_InactiveVSActive_ActiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Active"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Active", status)
}

func TestBestIndicatorStatus_MultipleTechnicalServiceMultipleIndicator_ActiveVSActive_ActiveResult(t *testing.T) {

	// Given
	services := []string{"jenkins", "gitlabci"}
	servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}, {Service: "gitlabci", Status: "Active"}}

	// When
	status := bestIndicatorStatus(services, servicesToMatch)

	// Then
	assert.Equal(t, "Active", status)
}
