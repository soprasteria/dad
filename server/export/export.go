package export

import (
	"bytes"

	"sort"

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

	createMergedCell(servicePkgRow, "Matrix Maturity", 4)

	createMergedCell(serviceNameRow, "", 4)

	createCell(serviceMaturityRow, "Domain")
	createCell(serviceMaturityRow, "Project")
	createCell(serviceMaturityRow, "Entity")
	createCell(serviceMaturityRow, "Service Center")

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
		projectRow := sheet.AddRow()

		entity, err := e.Database.Organizations.FindByIDBson(project.Entity)
		if err != nil {
			entity = types.Organization{Name: "N/A"}
		}

		serviceCenter, err := e.Database.Organizations.FindByIDBson(project.ServiceCenter)
		if err != nil {
			serviceCenter = types.Organization{Name: "N/A"}
		}

		createCell(projectRow, project.Domain)
		createCell(projectRow, project.Name)
		createCell(projectRow, entity.Name)
		createCell(projectRow, serviceCenter.Name)

		for _, pkg := range servicesMapSortedKeys {
			services := servicesMap[pkg]
			for _, service := range services {
				applicable := false
				for _, line := range project.Matrix {
					if line.Service == service.ID {
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
