package main


import "github.com/xuri/excelize"
func main(){

	xlsx:=excelize.NewFile()
	tabhead := map[string]string{"A1": "Apple", "B1": "Orange", "C1": "Pear"}
	values := map[string]int{"A2": 2, "B2": 3, "C2": 3, "A3": 5, "B3": 2, "C3": 4, "A4": 6, "B4": 7, "C4": 8}

	for k, v := range tabhead {
		xlsx.SetCellValue("Sheet1", k, v)
	}

	for k, v := range values {
		xlsx.SetCellValue("Sheet1", k, v)
	}

	xlsx.SaveAs("D:/1.xlsx")

}