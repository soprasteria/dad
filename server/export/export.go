package export

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/tealeg/xlsx"
)

// Export contains APIs entrypoints needed for accessing users
type Export struct {
	Database *mongo.DadMongo
}

func (e *Export) generateXlsx(projects []types.Project) (*bytes.Reader, error) {
	services, err := e.Database.FunctionalServices.FindAll()
	if err != nil {
		return nil, err
	}

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
		"Business",
		"Service Center",
		"Domain",
		"Client",
		"Project Manager",
		"Technologies",
		"Deployment Mode",
		"Version Control System",
		"Deliverables in VCS",
		"Source Code in VCS",
		"Specifications in VCS",
		"Creation Date",
		"Last Update",
		"Comments",
	}

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

	// Number of columns by service
	const nbColsService = 4

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
		}
	}

	// Generate a project row
	for _, project := range projects {
		var comments []string
		projectRow := sheet.AddRow()

		var businessUnit, serviceCenter types.Entity
		var projectManager types.User
		businessUnit, err = e.Database.Entities.FindByID(project.BusinessUnit)
		if err != nil {
			businessUnit = types.Entity{Name: "N/A"}
		}

		serviceCenter, err = e.Database.Entities.FindByID(project.ServiceCenter)
		if err != nil {
			serviceCenter = types.Entity{Name: "N/A"}
		}

		projectManager, err = e.Database.Users.FindByID(project.ProjectManager)
		if err != nil {
			projectManager = types.User{DisplayName: "N/A"}
		}

		if project.Domain == "" {
			project.Domain = "N/A"
		}

		createCell(projectRow, project.Name)
		createCell(projectRow, businessUnit.Name)
		createCell(projectRow, serviceCenter.Name)
		createCell(projectRow, project.Domain)
		createCell(projectRow, project.Client)
		createCell(projectRow, projectManager.DisplayName)
		createCell(projectRow, strings.Join(project.Technologies, ", "))
		createCell(projectRow, project.Mode)
		createCell(projectRow, project.VersionControlSystem)
		createBoolCell(projectRow, project.DeliverablesInVersionControl)
		createBoolCell(projectRow, project.SourceCodeInVersionControl)
		createBoolCell(projectRow, project.SpecificationsInVersionControl)
		createDateCell(projectRow, project.Created)
		createDateCell(projectRow, project.Updated)

		// Aggregate comments
		for _, pkg := range servicesMapSortedKeys {
			services := servicesMap[pkg]
			for _, service := range services {
				for _, line := range project.Matrix {
					if line.Service == service.ID {
						if line.Comment != "" {
							comments = append(comments, fmt.Sprintf("%s: %s: %s", pkg, service.Name, line.Comment))
						}
						break
					}
				}
			}
		}
		commentsString := strings.Join(comments, "\n")
		createCell(projectRow, commentsString)
		projectRow.SetHeightCM(0.5*float64(strings.Count(commentsString, "\n")) + 0.5)

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
						applicable = true
						break
					}
				}
				if !applicable {
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
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
func (e *Export) Export(projects []types.Project) (*bytes.Reader, error) {
	return e.generateXlsx(projects)
}
