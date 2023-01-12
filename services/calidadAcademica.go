package services

import (
	"UPCSToolkit/widgets"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func IndexOf(list []string, value string) int {
	for idx, val := range list {
		if val == value {
			return idx
		}
	}
	return -1
}

func getUniqueColumn(data [][]string, idx int) *[]string {
	var resultMap map[string]int = make(map[string]int) //create a map to avoid repeated values
	var result []string
	for i, row := range data {
		if i == 0 {
			continue
		}
		resultMap[row[idx]] = 1
	}

	for key := range resultMap {
		result = append(result, key)
	}

	return &result
}

func CreateReportCalidadAcademica(data [][]string) {
	var (
		School_Idx       int = IndexOf(data[0], "escuela")
		Major_Idx        int = IndexOf(data[0], "carrera")
		Project_Idx      int = IndexOf(data[0], "proyecto")
		Project_Type_Idx int = IndexOf(data[0], "tipo_proyecto_up")
		Finish_Idx       int = IndexOf(data[0], "terminacion_institucion")
		Punishment_Idx   int = IndexOf(data[0], "penalizacion")
	)

	//Checks that all the columns exists
	if School_Idx == -1 ||
		Major_Idx == -1 ||
		Project_Idx == -1 ||
		Project_Type_Idx == -1 ||
		Finish_Idx == -1 ||
		Punishment_Idx == -1 {
		widgets.NewPrompt("Error", "Error: the .csv selected does not meet the requirements")
		return // if not then is not a valid
	}
	var (
		school       func(int) string = func(idx int) string { return data[idx][School_Idx] }
		major        func(int) string = func(idx int) string { return data[idx][Major_Idx] }
		project      func(int) string = func(idx int) string { return data[idx][Project_Idx] }
		project_type func(int) string = func(idx int) string { return data[idx][Project_Type_Idx] }
		finish       func(int) string = func(idx int) string { return data[idx][Finish_Idx] }
		punishment   func(int) string = func(idx int) string { return data[idx][Punishment_Idx] }
	)

	var _ [6]string = [6]string{school(0), major(0), project(0), project_type(0), finish(0), punishment(0)}

	var majorList *[]string = getUniqueColumn(data, Major_Idx)
	var majorStats map[string][]int = make(map[string][]int)     //[major] = {Social- Social, Academico -Social}
	var projectStats map[string]string = make(map[string]string) //[major] = {Social- Social || Academico -Social}

	for _, major := range *majorList {
		majorStats[major] = []int{0, 0}
	}
	majorStats["total"] = []int{0, 0}
	var typeCount []int = []int{0, 0} //{Social - Social, Academico - Social}
	var invalidCount int = 0          //rows deleted

	for i := range data {
		if punishment(i) != "" || finish(i) != "Si" {
			if i != 0 {
				invalidCount++
			}
			continue
		}
		var (
			currentType    string = project_type(i)
			currentMajor   string = major(i)
			currentProject string = project(i)
		)
		if currentType == "Social - Social" {
			majorStats[currentMajor][0] += 1
			majorStats["total"][0] += 1
			if projectStats[currentProject] == "" {
				typeCount[0] += 1
			}

		} else {
			majorStats[currentMajor][1] += 1
			majorStats["total"][1] += 1
			if projectStats[currentProject] == "" {
				typeCount[1] += 1
			}

		}
		projectStats[currentProject] = currentType

	}
	showReport(majorStats, typeCount, invalidCount)
}

func showReport(majorStats map[string][]int, typeCount []int, invalidCount int) {
	reportWin := fyne.CurrentApp().NewWindow("Reporte Calidad Academica")
	//invalidLabel := widget.NewLabel(fmt.Sprintf("Alumnos Invalidos: %d", invalidCount))
	var majorKeys []string = make([]string, 0)
	for key := range majorStats {
		majorKeys = append(majorKeys, key)
	}
	table := widget.NewTable(
		//define the dimensions of the table
		func() (int, int) { return len(majorKeys), 3 },
		//specifies the widget
		func() fyne.CanvasObject { return widget.NewLabel("...") },
		//update the table
		func(tci widget.TableCellID, obj fyne.CanvasObject) {
			if tci.Col == 0 {
				text := majorKeys[tci.Row]
				if len(text) > 60 {
					text = text[:60]
				}
				obj.(*widget.Label).SetText(text)
				obj.(*widget.Label).Resize(fyne.NewSize(500, 10))
			} else {
				obj.(*widget.Label).SetText(fmt.Sprintf("%d", majorStats[majorKeys[tci.Row]][tci.Col-1]))
				obj.(*widget.Label).Resize(fyne.NewSize(100, 10))
			}

			if majorKeys[tci.Row] == "total" {
				obj.(*widget.Label).TextStyle.Bold = true
			}
		},
	)
	table.SetColumnWidth(0, 500)
	table.SetColumnWidth(1, 50)
	table.SetColumnWidth(2, 50)
	content := container.NewCenter(table)
	content.Resize(fyne.NewSize(600, 500))

	//container := container.NewVBox(invalidLabel, table)
	reportWin.Resize(fyne.NewSize(800, 500))
	reportWin.SetContent(table)
	reportWin.Show()
}
