package export

import (
	"bytes"

	"github.com/soprasteria/dad/server/mongo"
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

	createMergedCell(servicePkgRow, "Matric Maturity", 4)
	createMergedCell(serviceNameRow, "", 4)
	createCell(serviceMaturityRow, "Domain")
	createCell(serviceMaturityRow, "Project")
	createCell(serviceMaturityRow, "Entity")
	createCell(serviceMaturityRow, "Service Center")

	services, err := e.Database.FunctionnalServices.FindAll()
	if err != nil {
		return nil, err
	}

	servicesMap := make(map[string][]string)
	for _, service := range services {
		servicesMap[service.Package] = append(servicesMap[service.Package], service.Name)
	}

	for pkg, names := range servicesMap {
		createMergedCell(servicePkgRow, pkg, len(names)*2)
		for _, name := range names {

			nameCell := createMergedCell(serviceNameRow, name, 2)
			rotateCell(nameCell, 90)
			createCell(serviceMaturityRow, "Progress")
			createCell(serviceMaturityRow, "Goal")
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
