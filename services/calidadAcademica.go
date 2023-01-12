package services

import (
	"UPCSToolkit/widgets"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

func getDataMatrixFromStats(stats map[string][]int, headers []string) (result [][]string) {
	result = append(result, headers)
	for key, value := range stats {
		if key == "total" {
			continue
		}
		row := make([]string, 0)
		row = append(row, key)
		row = append(row, fmt.Sprint(value[0]))
		row = append(row, fmt.Sprint(value[1]))
		result = append(result, row)
	}

	return
}

func showReport(majorStats map[string][]int, typeCount []int, invalidCount int) {
	//creates the report window
	reportWin := fyne.CurrentApp().NewWindow("Reporte Calidad Academica")

	titleLabel := canvas.NewText("Reporte de Calidad Academica", widgets.DARK_BLUE)
	titleLabel.TextSize = 20

	projectTitle := canvas.NewText(
		"Conteo de Proyectos",
		widgets.DARK_BLUE,
	)
	projectTitle.TextSize = 14
	projectLabel := canvas.NewText(
		fmt.Sprintf(
			"Social-Social: %d     Academico-Social: %d     Total: %d",
			typeCount[0], typeCount[1], typeCount[0]+typeCount[1],
		),
		widgets.DARK_BLUE,
	)
	projectLabel.TextSize = 12
	headerContainer := container.NewVBox(titleLabel, projectTitle, projectLabel)

	var majorTableData [][]string = getDataMatrixFromStats(majorStats, []string{"Carrera", "Social - Social", "Academico - Social"})

	//appending the totals to the major Table
	majorTableData = append(majorTableData, []string{"sub-total", fmt.Sprint(majorStats["total"][0]), fmt.Sprint(majorStats["total"][1])})
	majorTableData = append(majorTableData, []string{"sub-total", fmt.Sprint(majorStats["total"][0] + majorStats["total"][1]), ""})
	majorTable := widget.NewTable(
		//define the dimensions of the table
		func() (int, int) { return len(majorTableData), 3 },

		//specifies the widget
		func() fyne.CanvasObject { return canvas.NewText("...", color.Black) },

		//fills the table
		func(tci widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*canvas.Text)
			//just apply bold style to the first and last row
			label.TextStyle.Bold = (tci.Row == 0 || tci.Row >= len(majorTableData)-2)
			label.Text = majorTableData[tci.Row][tci.Col]

			if tci.Row == 0 {
				label.Color = widgets.GOLD
				label.Alignment = fyne.TextAlignCenter
				return
			}
			label.Color = widgets.DARK_GREY

			if tci.Col == 0 {
				label.Alignment = fyne.TextAlignLeading //align left
				if len(label.Text) > 60 {
					label.Text = label.Text[:60]
				}
				if tci.Row >= len(majorTableData)-2 {
					label.Alignment = fyne.TextAlignTrailing //align right
				}
				return
			}

			label.Alignment = fyne.TextAlignCenter

		},
	)
	majorTable.SetColumnWidth(0, 400)
	majorTable.SetColumnWidth(1, 150)
	majorTable.SetColumnWidth(2, 150)

	reportLayout := layout.NewBorderLayout(headerContainer, nil, nil, nil)
	reportContainer := container.New(reportLayout, headerContainer, majorTable)

	//container := container.NewVBox(invalidLabel, table)
	reportWin.Resize(fyne.NewSize(800, 600))
	reportWin.SetContent(reportContainer)
	reportWin.Show()
}
