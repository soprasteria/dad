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

// ProjectDataExport contains the data to put inside the export
type ProjectDataExport struct {
	Name                           string
	Description                    string
	BusinessUnit                   string
	ServiceCenter                  string
	Domain                         []string
	Client                         string
	ProjectManager                 string
	Deputies                       []string
	DocktorGroupName               string
	DocktorGroupURL                string
	Technologies                   []string
	Mode                           string
	VersionControlSystem           string
	DeliverablesInVersionControl   bool
	SourceCodeInVersionControl     bool
	SpecificationsInVersionControl bool
	Created                        time.Time
	Updated                        time.Time
	ServicesDataExport             map[int][]ServiceDataExport
}

// ServiceDataExport contains the data of each service
type ServiceDataExport struct {
	Progress  string
	Goal      string
	Priority  string
	DueDate   *time.Time
	Indicator string
	Comment   string
}

// Headers contains headers for each project
type Headers struct {
	MatrixHeaders        []string
	ServicesInfosHeaders []string
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

// generateExportHeaders generate headers for projects
func generateExportHeaders() Headers {

	// Headers contained inside the Matrix maturity column
	matrixMaturityColumnsHeaders := []string{
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

	// Headers for each service
	serviceInfosColumnsHeaders := []string{
		"Progress",
		"Goal",
		"Priority",
		"Due Date",
		"Indicator",
		"Comment"}

	headers := Headers{}
	headers.MatrixHeaders = matrixMaturityColumnsHeaders
	headers.ServicesInfosHeaders = serviceInfosColumnsHeaders

	return headers
}

// retrieveData retrieve all data from projects, projectToUsageIndicators and servicesMap. These data are meant to be used for the generation of the export
func (e *Export) retrieveData(servicesMapSortedKeys []string, servicesMap map[string][]types.FunctionalService, projects []types.Project, projectToUsageIndicators map[string][]types.UsageIndicator) map[string]ProjectDataExport {

	exportData := make(map[string]ProjectDataExport)

	allServiceIndicatorMap := getServiceIndicatorMap(projects, servicesMapSortedKeys, servicesMap, projectToUsageIndicators)

	for _, project := range projects {

		if project.Name == "" {
			log.WithField("project.Name", project.Name).Error(fmt.Sprintf("The project %v has no name", project.Name))
			continue
		} else {

			currentProjectDataExport := ProjectDataExport{}
			currentProjectDataExport.Name = project.Name
			currentProjectDataExport.Description = project.Description

			businessUnit, err := e.Database.Entities.FindByID(project.BusinessUnit)
			if err != nil {
				currentProjectDataExport.BusinessUnit = "N/A"
			} else {
				currentProjectDataExport.BusinessUnit = businessUnit.Name
			}

			serviceCenter, err := e.Database.Entities.FindByID(project.ServiceCenter)
			if err != nil {
				currentProjectDataExport.ServiceCenter = "N/A"
			} else {
				currentProjectDataExport.ServiceCenter = serviceCenter.Name
			}

			if len(project.Domain) == 0 {
				currentProjectDataExport.Domain = []string{"N/A"}
			} else {
				currentProjectDataExport.Domain = project.Domain
			}

			currentProjectDataExport.Client = project.Client

			projectManager, err := e.Database.Users.FindByID(project.ProjectManager)
			if err != nil {
				currentProjectDataExport.ProjectManager = "N/A"
			}
			currentProjectDataExport.ProjectManager = projectManager.DisplayName

			deputies := e.findDeputies(project)
			currentProjectDataExport.Deputies = deputies
			currentProjectDataExport.DocktorGroupName = project.DocktorGroupName
			currentProjectDataExport.DocktorGroupURL = project.DocktorGroupURL
			currentProjectDataExport.Technologies = project.Technologies
			currentProjectDataExport.Mode = project.Mode
			currentProjectDataExport.VersionControlSystem = project.VersionControlSystem
			currentProjectDataExport.DeliverablesInVersionControl = project.DeliverablesInVersionControl
			currentProjectDataExport.SourceCodeInVersionControl = project.SourceCodeInVersionControl
			currentProjectDataExport.SpecificationsInVersionControl = project.SpecificationsInVersionControl
			currentProjectDataExport.Created = project.Created
			currentProjectDataExport.Updated = project.Updated

			servicesDataExportMap := make(map[int][]ServiceDataExport)

			// Iterate on each service in the correct order
			for indexPkg, pkg := range servicesMapSortedKeys {
				services := servicesMap[pkg]

				servicesDataExport := make([]ServiceDataExport, 0)
				for _, service := range services {
					applicable := false
					// Iterate on the project matrix and print the data for the current service

					serviceDataExport := ServiceDataExport{}
					for _, line := range project.Matrix {
						if line.Service == service.ID {
							serviceDataExport.Progress = types.Progress[line.Progress]
							serviceDataExport.Goal = types.Progress[line.Goal]
							serviceDataExport.Priority = line.Priority
							// If the DueDate is nil, N/A will be written while the generateXlsx function
							serviceDataExport.DueDate = line.DueDate

							key := ServiceProjectEntry{ProjectName: project.Name, ServiceName: service.Name}
							serviceDataExport.Indicator = allServiceIndicatorMap[key]
							serviceDataExport.Comment = line.Comment

							applicable = true
							break
						}
					}
					if !applicable {
						serviceDataExport.Progress = "N/A"
						serviceDataExport.Goal = "N/A"
						serviceDataExport.Priority = "N/A"
						// If the DueDate is nil, N/A will be written while the generateXlsx function
						serviceDataExport.DueDate = nil

						key := ServiceProjectEntry{ProjectName: project.Name, ServiceName: service.Name}
						serviceDataExport.Indicator = allServiceIndicatorMap[key]
						serviceDataExport.Comment = ""
					}
					servicesDataExport = append(servicesDataExport, serviceDataExport)
				}
				servicesDataExportMap[indexPkg] = servicesDataExport
			}
			currentProjectDataExport.ServicesDataExport = servicesDataExportMap
			exportData[project.Name] = currentProjectDataExport
		}
	}
	return exportData
}

// generateXlsx generate the export with informations associated
func generateXlsx(servicesMapSortedKeys []string, servicesMap map[string][]types.FunctionalService, headersExport Headers, dataExport map[string]ProjectDataExport) (*bytes.Reader, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Plan de déploiement")
	if err != nil {
		return nil, err
	}

	servicePkgRow := sheet.AddRow()
	serviceNameRow := sheet.AddRow()
	serviceMaturityRow := sheet.AddRow()

	// Header generation: package and associated functional services
	createMergedCell(servicePkgRow, "Matrix Maturity", len(headersExport.MatrixHeaders))
	createMergedCell(serviceNameRow, "Export Date: "+time.Now().Format("02/01/2006"), len(headersExport.MatrixHeaders))
	for _, matrixHeader := range headersExport.MatrixHeaders {
		createCell(serviceMaturityRow, matrixHeader)
	}

	for _, pkg := range servicesMapSortedKeys {
		services := servicesMap[pkg]

		createMergedCell(servicePkgRow, pkg, len(services)*len(headersExport.ServicesInfosHeaders))
		for _, service := range services {
			nameCell := createMergedCell(serviceNameRow, service.Name, len(headersExport.ServicesInfosHeaders))
			rotateCell(nameCell, 90)
			for _, servicesInfo := range headersExport.ServicesInfosHeaders {
				createCell(serviceMaturityRow, servicesInfo)
			}
		}
	}

	// Keep a list of the projects names
	dataExportKeys := []string{}
	for key := range dataExport {
		dataExportKeys = append(dataExportKeys, key)
	}
	sort.Strings(dataExportKeys)

	// Generate all project rows
	for _, projectName := range dataExportKeys {
		projectRow := sheet.AddRow()

		createCell(projectRow, dataExport[projectName].Name)
		createCell(projectRow, dataExport[projectName].Description)
		createCell(projectRow, dataExport[projectName].BusinessUnit)
		createCell(projectRow, dataExport[projectName].ServiceCenter)
		createCell(projectRow, strings.Join(dataExport[projectName].Domain, "; "))
		createCell(projectRow, dataExport[projectName].Client)
		createCell(projectRow, dataExport[projectName].ProjectManager)
		createCell(projectRow, strings.Join(dataExport[projectName].Deputies, ", "))
		createCell(projectRow, dataExport[projectName].DocktorGroupName)
		createCell(projectRow, dataExport[projectName].DocktorGroupURL)
		createCell(projectRow, strings.Join(dataExport[projectName].Technologies, ", "))
		createCell(projectRow, dataExport[projectName].Mode)
		createCell(projectRow, dataExport[projectName].VersionControlSystem)
		createBoolCell(projectRow, dataExport[projectName].DeliverablesInVersionControl)
		createBoolCell(projectRow, dataExport[projectName].SourceCodeInVersionControl)
		createBoolCell(projectRow, dataExport[projectName].SpecificationsInVersionControl)
		createDateCell(projectRow, dataExport[projectName].Created)
		createDateCell(projectRow, dataExport[projectName].Updated)

		servicesDataExportSortedKeys := make([]int, 0)
		for k := range dataExport[projectName].ServicesDataExport {
			servicesDataExportSortedKeys = append(servicesDataExportSortedKeys, k)
		}
		sort.Ints(servicesDataExportSortedKeys)

		for _, sortedKey := range servicesDataExportSortedKeys {

			serviceDataExport := dataExport[projectName].ServicesDataExport[sortedKey]
			for _, service := range serviceDataExport {

				createFormattedValueCell(projectRow, service.Progress)
				createFormattedValueCell(projectRow, service.Goal)
				createCell(projectRow, service.Priority)
				if service.DueDate == nil {
					createCell(projectRow, "N/A")
				} else {
					createDateCell(projectRow, *service.DueDate)
				}
				createCell(projectRow, service.Indicator)
				createCell(projectRow, service.Comment)
			}
		}
	}

	// Presentation of the sheet
	serviceNameRow.SetHeightCM(10)
	colorRow(servicePkgRow, red, white)
	colorRow(serviceNameRow, red, white)
	colorRow(serviceMaturityRow, red, white)
	modifySheetAlignment(sheet, "center", "center")
	modifySheetBorder(sheet, black)

	// Width for all cells
	setWidthCols(sheet, 12.0)

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

	headersExport := generateExportHeaders()

	dataExport := e.retrieveData(servicesMapSortedKeys, servicesMap, projects, projectToUsageIndicators)

	return generateXlsx(servicesMapSortedKeys, servicesMap, headersExport, dataExport)
}
