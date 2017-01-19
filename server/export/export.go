package export

import (
	"bytes"

	"github.com/soprasteria/dad/server/mongo"
	"github.com/tealeg/xlsx"
)

const red = "FFC00000"

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
	serviceNameRow.SetHeightCM(10)

	titleCell := servicePkgRow.AddCell()
	servicePkgRow.AddCell()
	titleCell.Merge(1, 0)
	titleCell.SetString("Matric Maturity")
	style := titleCell.GetStyle()

	fill := xlsx.NewFill("solid", red, "FFFFFFFF")
	style.Fill = *fill
	style.ApplyFill = true

	titleCell.SetStyle(style)

	domainCell := serviceNameRow.AddCell()
	domainCell.SetString("Domain")
	style = domainCell.GetStyle()

	fill = xlsx.NewFill("solid", red, "FFFFFFFF")
	style.Fill = *fill
	style.ApplyFill = true

	domainCell.SetStyle(style)

	projectCell := serviceNameRow.AddCell()
	projectCell.SetString("Project")
	style = projectCell.GetStyle()

	fill = xlsx.NewFill("solid", red, "FFFFFFFF")
	style.Fill = *fill
	style.ApplyFill = true

	projectCell.SetStyle(style)

	services, err := e.Database.FunctionnalServices.FindAll()
	if err != nil {
		return nil, err
	}

	servicesMap := make(map[string][]string)
	for _, service := range services {
		servicesMap[service.Package] = append(servicesMap[service.Package], service.Name)
	}

	for pkg, names := range servicesMap {
		pkgCell := servicePkgRow.AddCell()
		for i, name := range names {
			if i != len(names)-1 {
				servicePkgRow.AddCell()
			}
			nameCell := serviceNameRow.AddCell()
			style := nameCell.GetStyle()
			style.Alignment.TextRotation = 90
			style.ApplyAlignment = true

			fill := xlsx.NewFill("solid", red, "FFFFFFFF")
			style.Fill = *fill
			style.ApplyFill = true

			nameCell.SetStyle(style)
			nameCell.SetString(name)
		}
		pkgCell.Merge(len(names)-1, 0)

		style := pkgCell.GetStyle()
		fill := xlsx.NewFill("solid", red, "FFFFFFFF")
		style.Fill = *fill
		style.ApplyFill = true

		pkgCell.SetStyle(style)
		pkgCell.SetString(pkg)
	}

	// Write the file in-memory and returns is as a readable stream
	var b bytes.Buffer
	err = file.Write(&b)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b.Bytes()), nil
}
