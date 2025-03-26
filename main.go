package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/openact/api"
	"github.com/openact/kit/win"
	"github.com/openact/lib/utils"
	"github.com/openact/productParser/conf"
	"log"
	"os"
	"sort"
	"time"
)

//go:embed vault
var fs embed.FS

func main() {
	utils.SetLogOutput()
	start := time.Now()
	//Licensing
	Dpath := win.ReleaseVault(&fs)
	valid := win.LoadYamlLicense(Dpath, "level1.yaml", "2025-06-30")

	if valid == false {
		return
	}
	for _, run := range conf.Runs {
		lib := api.Library{Name: run.Name, Variables: make(map[string]map[string]api.Definition)}
		folderPath := conf.OutputPath + "/" + run.Folder
		utils.InitializePath(folderPath)

		for i, repo := range run.CodeRepo {
			fmt.Println(i, repo)
			prodName, _ := utils.FilePathToName(repo)
			utils.InitializePath(folderPath + "/" + prodName)
			fmt.Println(run.Name)
			api.ParseProduct(repo, &lib)
		}

		var varNames []string
		for varName, _ := range lib.Variables {
			varNames = append(varNames, varName)
		}
		sort.Strings(varNames)

		//print out the library
		f, err := os.Create(folderPath + "/" + "LibPivot" + ".csv")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		writer := csv.NewWriter(f)

		header := []string{"Name", "Description"}
		for _, repo := range run.CodeRepo {
			prodName, _ := utils.FilePathToName(repo)
			header = append(header, prodName, prodName)
		}

		err = writer.Write(header)
		if err != nil {
			panic(err)
		}

		for _, varName := range varNames {
			row := []string{varName}
			varMap := lib.Variables[varName]
			defnMap := make(map[string]string)
			for i, repo := range run.CodeRepo {
				prodName, _ := utils.FilePathToName(repo)
				prodPath := folderPath + "/" + prodName

				definition, ok := varMap[prodName]
				cat := definition.Cat
				defnType := api.LookupIdent(definition.Type)
				defn := ""

				if !ok {
					defn = ""
				} else {
					defn = definition.Defn
					mdFile, err := os.Create(prodPath + "/" + varName + ".md")
					if err != nil {
						panic(err)
					}

					mdFile.WriteString(varName + "\n")
					mdFile.WriteString("\n")
					mdFile.WriteString(cat + " " + defnType + " Definition" + "\n")
					mdFile.WriteString("```" + "\n")
					mdFile.WriteString(definition.Defn + "\n")
					mdFile.WriteString("```" + "\n")
					mdFile.WriteString("\n")
					mdFile.Close()
				}

				if sameAsProd, exists := defnMap[defn]; exists {
					defn = fmt.Sprintf("Same as %s", sameAsProd)
				} else {
					if defn != "" {
						defnMap[defn] = prodName
					}
				}

				if i == 0 {
					row = append(row, varMap[prodName].Desc)
				}
				if defnType == "Extended Formula" {
					defn = "Extended Formula"
				} else if defnType == "t-Dependent Extended Formula" {
					defn = "t-Dependent Extended Formula"
				}
				row = append(row, cat, defn)
			}
			err = writer.Write(row)
		}
		writer.Flush()
	}
	//logging
	{
		elapsed := time.Since(start)
		fmt.Printf("Product Formula Parser ended and it took %s", elapsed)
		log.Printf("Product Formula Parser ended and it took %s", elapsed)
	}
}
