package export

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/tealeg/xlsx"
)

// Export contains APIs entrypoints needed for accessing users
type Export struct {
	Database *mongo.DadMongo
}

//ExportAll exports all business data as a file
func (e *Export) ExportAll() (*bytes.Reader, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Plan de déploiement")
	if err != nil {
		return nil, err
	}

	servicePkgRow := sheet.AddRow()
	serviceNameRow := sheet.AddRow()
	serviceMaturityRow := sheet.AddRow()

	serviceNameRow.SetHeightCM(10)

	createMergedCell(servicePkgRow, "Matrix Maturity", 5)

	createMergedCell(serviceNameRow, "", 5)

	createCell(serviceMaturityRow, "Domain")
	createCell(serviceMaturityRow, "Project")
	createCell(serviceMaturityRow, "Business Unit")
	createCell(serviceMaturityRow, "Service Center")
	createCell(serviceMaturityRow, "Project Manager")

	services, err := e.Database.FunctionnalServices.FindAll()
	if err != nil {
		return nil, err
	}

	// Build a map of services indexed by their package name
	servicesMap := make(map[string][]types.FunctionnalService)
	for _, service := range services {
		servicesMap[service.Package] = append(servicesMap[service.Package], service)
	}

	// Keep a list of the sorted package names
	servicesMapSortedKeys := []string{}
	for key := range servicesMap {
		servicesMapSortedKeys = append(servicesMapSortedKeys, key)
	}
	sort.Strings(servicesMapSortedKeys)

	// Header generation: package and associated functionnal services
	for _, pkg := range servicesMapSortedKeys {
		services := servicesMap[pkg]

		createMergedCell(servicePkgRow, pkg, len(services)*2)
		for _, service := range services {
			nameCell := createMergedCell(serviceNameRow, service.Name, 2)
			rotateCell(nameCell, 90)
			createCell(serviceMaturityRow, "Progress")
			createCell(serviceMaturityRow, "Goal")
		}
	}

	// Generate a project row
	projects, err := e.Database.Projects.FindAll()
	for _, project := range projects {
		var comments []string
		projectRow := sheet.AddRow()

		businessUnit, err := e.Database.Entities.FindByIDBson(project.BusinessUnit)
		if err != nil {
			businessUnit = types.Entity{Name: "N/A"}
		}

		serviceCenter, err := e.Database.Entities.FindByIDBson(project.ServiceCenter)
		if err != nil {
			serviceCenter = types.Entity{Name: "N/A"}
		}

		projectManager, err := e.Database.Users.FindByIDBson(project.ProjectManager)
		if err != nil {
			projectManager = types.User{DisplayName: "N/A"}
		}

		createCell(projectRow, project.Domain)
		createCell(projectRow, project.Name)
		createCell(projectRow, businessUnit.Name)
		createCell(projectRow, serviceCenter.Name)
		createCell(projectRow, projectManager.DisplayName)

		// Iterate on each service in the correct order
		for _, pkg := range servicesMapSortedKeys {
			services := servicesMap[pkg]
			for _, service := range services {
				applicable := false
				// Iterate on the project matrix and print the data for the current service
				for _, line := range project.Matrix {
					if line.Service == service.ID {
						comments = append(comments, fmt.Sprintf("%s: %s: %s", pkg, service.Name, line.Comment))
						createCell(projectRow, types.Progress[line.Progress])
						createCell(projectRow, types.Progress[line.Goal])
						applicable = true
						break
					}
				}
				if !applicable {
					createCell(projectRow, "N/A")
					createCell(projectRow, "N/A")
				}
			}
		}
		createCell(projectRow, strings.Join(comments, "\n"))
	}

	colorRow(servicePkgRow, red, white)
	colorRow(serviceNameRow, red, white)
	colorRow(serviceMaturityRow, red, white)
	modifySheetAlignment(sheet, "center", "center")
	modifySheetBorder(sheet, black)

	// Write the file in-memory and returns is as a readable stream
	var b bytes.Buffer
	err = file.Write(&b)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b.Bytes()), nil
}
