package export

import (
	"time"

	"github.com/tealeg/xlsx"
)

const red = "FFC00000"
const white = "FFFFFFFF"
const black = "FF000000"

func createCell(row *xlsx.Row, content string) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetString(content)
	return cell
}

func createBoolCell(row *xlsx.Row, boolean bool) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetBool(boolean)
	return cell
}

func createDateCell(row *xlsx.Row, date time.Time) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetDate(date)
	return cell
}

func createFormattedValueCell(row *xlsx.Row, numberFormat string) *xlsx.Cell {
	cell := row.AddCell()
	cell.SetString(numberFormat)
	cell.NumFmt = "0%"
	//nolint:errcheck
	cell.FormattedValue()
	return cell
}

func createMergedCell(row *xlsx.Row, content string, size int) *xlsx.Cell {
	if size < 1 {
		size = 1
	}
	cell := row.AddCell()
	for i := 1; i < size; i++ {
		row.AddCell()
	}
	cell.Merge(size-1, 0)
	cell.SetString(content)
	return cell
}

func colorRow(row *xlsx.Row, bgColor string, color string) {
	for _, cell := range row.Cells {
		colorCell(cell, bgColor, color)
	}
}

func colorCell(cell *xlsx.Cell, bgColor string, color string) {
	style := cell.GetStyle()
	fill := xlsx.NewFill("solid", bgColor, color)
	style.Fill = *fill
	style.ApplyFill = true
	cell.SetStyle(style)
}

func rotateCell(cell *xlsx.Cell, degrees int) {
	style := cell.GetStyle()
	style.Alignment.TextRotation = degrees
	style.ApplyAlignment = true
	cell.SetStyle(style)
}

func modifySheetAlignment(sheet *xlsx.Sheet, horizontal string, vertical string) {
	for _, row := range sheet.Rows {
		modifyRowAlignment(row, horizontal, vertical)
	}
}

func modifyRowAlignment(row *xlsx.Row, horizontal string, vertical string) {
	for _, cell := range row.Cells {
		modifyCellAlignment(cell, horizontal, vertical)
	}
}

func modifyCellAlignment(cell *xlsx.Cell, horizontal string, vertical string) {
	style := cell.GetStyle()
	style.Alignment.Horizontal = horizontal
	style.Alignment.Vertical = vertical
	style.ApplyAlignment = true
	cell.SetStyle(style)
}

func modifySheetBorder(sheet *xlsx.Sheet, color string) {
	for _, row := range sheet.Rows {
		modifyRowBorder(row, color)
	}
}

func modifyRowBorder(row *xlsx.Row, color string) {
	for _, cell := range row.Cells {
		modifyCellBorder(cell, true, true, true, true, color)
	}
}

func modifyCellBorder(cell *xlsx.Cell, left bool, right bool, top bool, bottom bool, color string) {
	border := xlsx.DefaultBorder()
	const borderSize = "thin"
	if left {
		border.Left = borderSize
		border.LeftColor = color
	}
	if right {
		border.Right = borderSize
		border.RightColor = color
	}
	if top {
		border.Top = borderSize
		border.TopColor = color
	}
	if bottom {
		border.Bottom = borderSize
		border.BottomColor = color
	}
	style := cell.GetStyle()
	style.Border = *border
	style.ApplyBorder = true
	cell.SetStyle(style)
}

func setWidthCols(sheet *xlsx.Sheet, width float64) {
	//nolint:errcheck
	sheet.SetColWidth(0, len(sheet.Cols), width)
}
