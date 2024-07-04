package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type saveOutput struct {
	savedOutput []byte
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

func extractTarGz(gzipStream io.Reader, path string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}
	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}
		p := filepath.Join(path, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, 0755); err != nil {
				fmt.Println(err)
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(p)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()
		default:
			return fmt.Errorf(
				"ExtractTarGz: uknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}
	}

	return nil
}

func getFileName(name string) string {
	for {
		if !strings.Contains(name, ".") {
			return name
		}
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}
}

func defaultString(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

func main() {
	e := echo.New()
	g := e.Group(defaultString(os.Getenv("CONTEXT_PATH"), ""))

	// ping
	g.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// upload
	g.POST("/upload", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		studyBaseDir := os.Getenv("STUDY_BASE_DIR")
		err = extractTarGz(src, studyBaseDir)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var so saveOutput
		studyDir := studyBaseDir + getFileName(file.Filename) + "/"
		url := os.Getenv("CBIOPORTAL_URL")
		genePanelPresent := false

		// import gene panels if present
		fmt.Println("--------------importing gene panels, if present")
		err = filepath.Walk(studyDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			var fileData string
			if !info.IsDir() {
				fileData1, err := os.ReadFile(path)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				fileData = string(fileData1)
			}
			isGenePanelFile := !strings.Contains(path, "._") && strings.Contains(path, "_gene_panel_") && fileData != "" && strings.Contains(fileData, "stable_id") && strings.Contains(fileData, "description") && strings.Contains(fileData, "gene_list")
			fmt.Printf("---%s: %v\n", path, isGenePanelFile)
			if isGenePanelFile {
				genePanelPresent = true
				fmt.Printf("File present - %s\n", path)
				args := []string{
					"--data",
					path,
				}
				cmd := exec.Command("./importGenePanel.pl", args...)
				cmd.Dir = "core/scripts"
				cmd.Stdin = os.Stdin
				cmd.Stdout = &so
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				if err := cmd.Wait(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				fmt.Println("---imported")
			}
			return nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		fmt.Println("----------------done importing gene panels")

		fmt.Println("----------------clearing cache")
		if genePanelPresent {
			fmt.Printf("%s/api/cache\n", url)
			req, err := http.NewRequest(http.MethodDelete, url+"/api/cache", nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			req.Header.Set("X-API-KEY", os.Getenv("X-API-KEY"))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			resp_str, err := io.ReadAll(resp.Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			fmt.Printf("%s\n%d\n", resp_str, resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("error importing gene panels for the study")
			}
		}
		fmt.Println("----------------done clearing cache")

		fmt.Println("----------------importing study")
		args := []string{
			"-u",
			os.Getenv("CBIOPORTAL_URL"),
			"-s",
			studyDir,
			"-o",
		}
		cmd := exec.Command("./metaImport.py", args...)
		cmd.Dir = "core/scripts/importer"
		cmd.Stdin = os.Stdin
		cmd.Stdout = &so
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := cmd.Wait(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		fmt.Println("----------------done importing study")

		return c.String(http.StatusOK, file.Filename)
	})

	e.Logger.Fatal(e.Start(":9000"))
}
