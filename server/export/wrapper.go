package export

import (
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
	if left {
		border.Left = "thin"
		border.LeftColor = color
	}
	if right {
		border.Right = "thin"
		border.RightColor = color
	}
	if top {
		border.Top = "thin"
		border.TopColor = color
	}
	if bottom {
		border.Bottom = "thin"
		border.BottomColor = color
	}
	style := cell.GetStyle()
	style.Border = *border
	style.ApplyBorder = true
	cell.SetStyle(style)
}
