package routes

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/invisibl-cloud/cbioportal-apiserver/pkg/cmd"
	"github.com/invisibl-cloud/cbioportal-apiserver/pkg/utils"
	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	studyBaseDir := os.Getenv("STUDY_BASE_DIR")
	err = utils.ExtractTarGz(src, studyBaseDir)
	if err != nil {
		return err
	}

	studyDir := studyBaseDir + utils.GetFilenameWithoutExtn(file.Filename) + "/"
	url := os.Getenv("CBIOPORTAL_URL")
	genePanelPresent := false

	// import gene panels if present
	fmt.Println("--------------importing gene panels, if present")
	err = filepath.Walk(studyDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var fileData string
		if !info.IsDir() {
			fileData1, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fileData = string(fileData1)
		}
		isGenePanelFile := !strings.Contains(path, "._") && strings.Contains(path, "_gene_panel_") && fileData != "" && strings.Contains(fileData, "stable_id") && strings.Contains(fileData, "description") && strings.Contains(fileData, "gene_list")
		fmt.Printf("---%s: %v\n", path, isGenePanelFile)
		if isGenePanelFile {
			genePanelPresent = true
			fmt.Printf("File present - %s\n", path)
			err := cmd.New("./importGenePanel.pl", "core/scripts", "--data", path).Execute()
			if err != nil {
				return err
			}
			fmt.Println("---imported")
		}
		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("----------------done importing gene panels")

	// flush cache
	fmt.Println("----------------clearing cache")
	if genePanelPresent {
		fmt.Printf("%s/api/cache\n", url)
		res, err := req.NewClient().R().SetHeader("X-API-KEY", os.Getenv("X-API-KEY")).Delete(url + "/api/cache")
		if err != nil {
			return err
		}
		if !res.IsSuccessState() || res.IsErrorState() {
			return fmt.Errorf("error clearing cache: %s", res.String())
		}
	}
	fmt.Println("----------------done clearing cache")

	// validate study
	fmt.Println("----------------validating study")
	err = cmd.New("./validateData.py", "core/scripts/importer", "-s", studyDir, "-n").Execute()
	if err != nil {
		return err
	}
	fmt.Println("----------------done validating study")

	// import study
	fmt.Println("----------------importing study")
	err = cmd.New("./metaImport.py", "core/scripts/importer", "-u", os.Getenv("CBIOPORTAL_URL"), "-s", studyDir, "-o").Execute()
	if err != nil {
		return err
	}
	fmt.Println("----------------done importing study")

	return nil
}
