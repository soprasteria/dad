package export

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/tealeg/xlsx"
)

// Export contains APIs entrypoints needed for accessing users
type Export struct {
	Database *mongo.DadMongo
}

// ServiceProjectEntry contains a specific service name for a specific project name
type ServiceProjectEntry struct {
	ProjectName string
	ServiceName string
}

// Status represents the different status possible for a service (like Jenkins)
type Status int

const (
	// Empty means that a the service does not have any project configuration. e.g. jenkins doesn't have a job
	Empty Status = iota
	// Undetermined means that a there is an incompatibilty in indicators results. e.g jenkins has jobs but no CPU activity is available
	Undetermined
	// Inactive means that a the service is configured but not used recently. e.g. jenkins has at least one job but its CPU usage is below the defined threshold
	Inactive
	// Active means that a the service is configured and used recently. e.g. jenkins has at least one job and its CPU usage is above the defined threshold
	Active
)

// statusStr represents the order of the Status, meaning the first status is the worse, and the last one is the best.
var statusStr = [...]string{
	"Empty",
	"Undetermined",
	"Inactive",
	"Active",
}

// statusMap is defining the matching between a string status and the real enum status.
// It's initialized in init function
var statusMap = make(map[string]Status)

func init() {
	for i, s := range statusStr {
		statusMap[s] = Status(i)
	}
}

// String function will return the string representation of a service Status (e.g. Jenkins)
func (status Status) String() string {
	return statusStr[status]
}

// GetStatus will return the enum representation of a service Status (e.g. Jenkins)
// returns an error if string status is unrecognized
func GetStatus(status string) (Status, error) {
	if v, ok := statusMap[status]; ok {
		return v, nil
	}
	return Undetermined, fmt.Errorf("Status %q does not exists", status)
}

func (e *Export) findDeputies(project types.Project) []string {
	var deputies []string
	for _, deputyID := range project.Deputies {
		deputy, err := e.Database.Users.FindByID(deputyID)
		if err != nil {
			deputy = types.User{DisplayName: "Invalid User"}
		}
		deputies = append(deputies, deputy.DisplayName)
	}
	return deputies
}

// getServiceToIndicatorUsage map which associates an indicator to a service if the indicator's service matches the service
func getServiceToIndicatorUsage(service types.FunctionalService, usageIndicators []types.UsageIndicator) map[string]types.UsageIndicator {
	var serviceToUsageIndicator = make(map[string]types.UsageIndicator)
	for _, service := range service.Services {
		for _, usageIndicator := range usageIndicators {
			if usageIndicator.Service == service {
				serviceToUsageIndicator[service] = usageIndicator
			}
		}
	}
	return serviceToUsageIndicator
}

// bestIndicatorStatus returns the best indicator status from an array of UsageIndicator which contains indicator status
func bestIndicatorStatus(service types.FunctionalService, usageIndicators []types.UsageIndicator) *Status {
	var currentStatus *Status
	if len(service.Services) > 0 && len(usageIndicators) > 0 {
		usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
		for _, currentService := range service.Services {
			if currentService == usageIndicator[currentService].Service {
				newStatus, err := GetStatus(usageIndicator[currentService].Status)
				if err != nil {
					log.WithError(err).Warn(fmt.Sprintf("The indicator status '%s' doesn't match the status list [Empty, Undetermined, Inactive, Active]", usageIndicator[currentService].Status))
				} else {
					if currentStatus == nil || *currentStatus < newStatus {
						currentStatus = &newStatus
					}
				}
			}
		}
	}
	return currentStatus
}

// getServiceIndicatorMap map which contains all indicator status for each services or by default N/A
func getServiceIndicatorMap(projects []types.Project, servicesMapSortedKeys []string, servicesMap map[string][]types.FunctionalService, projectToUsageIndicators map[string][]types.UsageIndicator) map[ServiceProjectEntry]string {

	serviceIndicatorMap := make(map[ServiceProjectEntry]string)

	for _, project := range projects {
		for _, pkg := range servicesMapSortedKeys {
			services := servicesMap[pkg]
			for _, service := range services {
				usageIndicators := projectToUsageIndicators[project.Name]
				newServiceProjectEntry := ServiceProjectEntry{
					ProjectName: project.Name,
					ServiceName: service.Name}
				status := bestIndicatorStatus(service, usageIndicators)
				if status != nil {
					serviceIndicatorMap[newServiceProjectEntry] = (*status).String()
				} else {
					serviceIndicatorMap[newServiceProjectEntry] = "N/A"
				}
			}
		}
	}
	return serviceIndicatorMap
}

func (e *Export) generateXlsx(projects []types.Project, services []types.FunctionalService, projectToUsageIndicators map[string][]types.UsageIndicator) (*bytes.Reader, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Plan de déploiement")
	if err != nil {
		return nil, err
	}

	servicePkgRow := sheet.AddRow()
	serviceNameRow := sheet.AddRow()
	serviceMaturityRow := sheet.AddRow()

	serviceNameRow.SetHeightCM(10)

	// Name of columns contained inside the Matrix maturity column
	matrixMaturityColumns := []string{
		"Project",
		"Description",
		"Business",
		"Service Center",
		"Consolidation Criteria",
		"Client",
		"Project Manager",
		"Deputies",
		"Docktor Group Name",
		"Docktor Group URL",
		"Technologies",
		"Deployment Mode",
		"Version Control System",
		"Deliverables in VCS",
		"Source Code in VCS",
		"Specifications in VCS",
		"Creation Date",
		"Last Update"}

	createMergedCell(servicePkgRow, "Matrix Maturity", len(matrixMaturityColumns))

	createMergedCell(serviceNameRow, "Export Date: "+time.Now().Format("02/01/2006"), len(matrixMaturityColumns))

	for _, column := range matrixMaturityColumns {
		createCell(serviceMaturityRow, column)
	}

	// Build a map of services indexed by their package name
	servicesMap := make(map[string][]types.FunctionalService)
	for _, service := range services {
		servicesMap[service.Package] = append(servicesMap[service.Package], service)
	}

	// Keep a list of the sorted package names
	servicesMapSortedKeys := []string{}
	for key := range servicesMap {
		servicesMapSortedKeys = append(servicesMapSortedKeys, key)
	}
	sort.Strings(servicesMapSortedKeys)

	allServiceIndicatorMap := getServiceIndicatorMap(projects, servicesMapSortedKeys, servicesMap, projectToUsageIndicators)

	// Number of columns by service
	const nbColsService = 6

	// Header generation: package and associated functional services
	for _, pkg := range servicesMapSortedKeys {
		services := servicesMap[pkg]

		createMergedCell(servicePkgRow, pkg, len(services)*nbColsService)
		for _, service := range services {
			nameCell := createMergedCell(serviceNameRow, service.Name, nbColsService)
			rotateCell(nameCell, 90)
			createCell(serviceMaturityRow, "Progress")
			createCell(serviceMaturityRow, "Goal")
			createCell(serviceMaturityRow, "Priority")
			createCell(serviceMaturityRow, "Due Date")
			createCell(serviceMaturityRow, "Indicator")
			createCell(serviceMaturityRow, "Comment")
		}
	}

	// Generate a project row
	for _, project := range projects {
		projectRow := sheet.AddRow()

		var businessUnit, serviceCenter types.Entity
		businessUnit, err = e.Database.Entities.FindByID(project.BusinessUnit)
		if err != nil {
			businessUnit = types.Entity{Name: "N/A"}
		}

		serviceCenter, err = e.Database.Entities.FindByID(project.ServiceCenter)
		if err != nil {
			serviceCenter = types.Entity{Name: "N/A"}
		}

		var projectManager types.User
		projectManager, err = e.Database.Users.FindByID(project.ProjectManager)
		if err != nil {
			projectManager = types.User{DisplayName: "N/A"}
		}

		deputies := e.findDeputies(project)

		if len(project.Domain) == 0 {
			project.Domain = []string{"N/A"}
		}

		createCell(projectRow, project.Name)
		createCell(projectRow, project.Description)
		createCell(projectRow, businessUnit.Name)
		createCell(projectRow, serviceCenter.Name)
		createCell(projectRow, strings.Join(project.Domain, "; "))
		createCell(projectRow, project.Client)
		createCell(projectRow, projectManager.DisplayName)
		createCell(projectRow, strings.Join(deputies, ", "))
		createCell(projectRow, project.DocktorGroupName)
		createCell(projectRow, project.DocktorGroupURL)
		createCell(projectRow, strings.Join(project.Technologies, ", "))
		createCell(projectRow, project.Mode)
		createCell(projectRow, project.VersionControlSystem)
		createBoolCell(projectRow, project.DeliverablesInVersionControl)
		createBoolCell(projectRow, project.SourceCodeInVersionControl)
		createBoolCell(projectRow, project.SpecificationsInVersionControl)
		createDateCell(projectRow, project.Created)
		createDateCell(projectRow, project.Updated)

		// Iterate on each service in the correct order
		for _, pkg := range servicesMapSortedKeys {
			services := servicesMap[pkg]
			for _, service := range services {
				applicable := false
				// Iterate on the project matrix and print the data for the current service
				for _, line := range project.Matrix {
					if line.Service == service.ID {
						createFormattedValueCell(projectRow, types.Progress[line.Progress])
						createFormattedValueCell(projectRow, types.Progress[line.Goal])
						createCell(projectRow, line.Priority)
						if line.DueDate != nil {
							createDateCell(projectRow, *line.DueDate)
						} else {
							createCell(projectRow, "N/A")
						}
						key := ServiceProjectEntry{ProjectName: project.Name, ServiceName: service.Name}
						createCell(projectRow, allServiceIndicatorMap[key])
						createCell(projectRow, line.Comment)
						applicable = true
						break
					}
				}

				if !applicable {
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
					key := ServiceProjectEntry{ProjectName: project.Name, ServiceName: service.Name}
					createCell(projectRow, allServiceIndicatorMap[key])
					createCell(projectRow, "")
				}
			}
		}
	}

	colorRow(servicePkgRow, red, white)
	colorRow(serviceNameRow, red, white)
	colorRow(serviceMaturityRow, red, white)
	modifySheetAlignment(sheet, "center", "center")
	modifySheetBorder(sheet, black)

	// Width for all cells
	const widthDate = 12.0
	setWidthCols(sheet, widthDate)

	// Write the file in-memory and returns is as a readable stream
	var b bytes.Buffer
	err = file.Write(&b)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b.Bytes()), nil
}

//Export exports some business data as a file
func (e *Export) Export(projects []types.Project, projectToUsageIndicators map[string][]types.UsageIndicator) (*bytes.Reader, error) {
	services, err := e.Database.FunctionalServices.FindAll()
	if err != nil {
		return nil, err
	}
	return e.generateXlsx(projects, services, projectToUsageIndicators)
}
