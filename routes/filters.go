package routes

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/invisibl-cloud/cbioportal-apiserver/pkg/db"
	"github.com/invisibl-cloud/cbioportal-apiserver/pkg/types"

	"github.com/labstack/echo/v4"
)

func GetFilters(c echo.Context) (map[string][]string, error) {
	clnt := db.GetDBClient()
	filterType := c.QueryParam("filterType")
	var q string
	var id string
	switch filterType {
	case "byTreatment":
		id = "treatment"
		q = "select distinct `value` from clinical_event_data where `key` = 'AGENT';"
	case "bySourceSite":
		id = "sourceSite"
		q = "select distinct ATTR_VALUE from clinical_sample where `ATTR_ID` = 'TISSUE_SOURCE_SITE';"
	default:
		return nil, fmt.Errorf("unknown filter type `%s`", filterType)
	}
	results, err := clnt.Query(q)
	if err != nil {
		return nil, err
	}
	filters := []string{}
	for results.Next() {
		var t string
		err = results.Scan(&t)
		if err != nil {
			return nil, err
		}
		if t != "" {
			filters = append(filters, t)
		}
	}
	return map[string][]string{
		id: filters,
	}, nil
}

func GetStudiesWithFilters(c echo.Context) ([]types.CancerStudy, error) {
	clnt := db.GetDBClient()
	filterType := c.QueryParam("filterType")

	var q string
	var filterValues []string
	switch filterType {
	case "byTreatment":
		treatment := c.QueryParam("treatment")
		if treatment == "" {
			return nil, fmt.Errorf("missing query param `treatment` for filtering studies by treatment")
		}
		filterValues = strings.Split(treatment, ",")
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID
		FROM cancer_study
		JOIN patient ON patient.CANCER_STUDY_ID = cancer_study.CANCER_STUDY_ID
		JOIN clinical_event ON clinical_event.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_event_data ON clinical_event_data.CLINICAL_EVENT_ID = clinical_event.CLINICAL_EVENT_ID AND clinical_event_data.key = 'AGENT' AND clinical_event_data.value IN (%s);
		`
	case "bySourceSite":
		sourceSite := c.QueryParam("sourceSite")
		if sourceSite == "" {
			return nil, fmt.Errorf("missing query param `source site` for filtering studies by tissue source site")
		}
		filterValues = strings.Split(sourceSite, ",")
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID
		FROM cancer_study
		JOIN patient ON patient.CANCER_STUDY_ID = cancer_study.CANCER_STUDY_ID
		JOIN sample ON sample.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_sample ON clinical_sample.INTERNAL_ID = sample.INTERNAL_ID AND clinical_sample.ATTR_ID = 'TISSUE_SOURCE_SITE' AND clinical_sample.ATTR_VALUE IN (%s);
		`
	default:
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID 
		FROM cancer_study;
		`
	}

	q = fmt.Sprintf(q, "'"+strings.Join(filterValues, "', '")+"'")
	fmt.Println(q)

	results, err := clnt.Query(q)
	if err != nil {
		return nil, err
	}
	out := []types.CancerStudy{}
	for results.Next() {
		var t types.CancerStudy
		err = results.Scan(&t.ID, &t.Name, &t.Desc, &t.CancerTypeID)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}

func GetStudiesWithFiltersV2(c echo.Context) ([]types.CancerStudy, error) {
	filters := map[string][]string{}
	if err := c.Bind(&filters); err != nil {
		return nil, err
	}
	treatment := filters["treatment"]
	sourceSite := filters["sourceSite"]

	clnt := db.GetDBClient()
	var q string
	var results *sql.Rows
	var err error
	switch {
	case len(treatment) > 0 && len(sourceSite) == 0:
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID
		FROM cancer_study
		JOIN patient ON patient.CANCER_STUDY_ID = cancer_study.CANCER_STUDY_ID
		JOIN clinical_event ON clinical_event.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_event_data ON clinical_event_data.CLINICAL_EVENT_ID = clinical_event.CLINICAL_EVENT_ID AND clinical_event_data.key = 'AGENT' AND clinical_event_data.value IN (%s);
		`
		q = fmt.Sprintf(q, "'"+strings.Join(treatment, "', '")+"'")
		results, err = clnt.Query(q)
		if err != nil {
			return nil, err
		}
	case len(sourceSite) > 0 && len(treatment) == 0:
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID
		FROM cancer_study
		JOIN patient ON patient.CANCER_STUDY_ID = cancer_study.CANCER_STUDY_ID
		JOIN sample ON sample.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_sample ON clinical_sample.INTERNAL_ID = sample.INTERNAL_ID AND clinical_sample.ATTR_ID = 'TISSUE_SOURCE_SITE' AND clinical_sample.ATTR_VALUE IN (%s);
		`
		q = fmt.Sprintf(q, "'"+strings.Join(sourceSite, "', '")+"'")
		results, err = clnt.Query(q)
		if err != nil {
			return nil, err
		}
	case len(treatment) > 0 && len(sourceSite) > 0:
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID
		FROM cancer_study
		JOIN patient ON patient.CANCER_STUDY_ID = cancer_study.CANCER_STUDY_ID
		JOIN clinical_event ON clinical_event.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_event_data ON clinical_event_data.CLINICAL_EVENT_ID = clinical_event.CLINICAL_EVENT_ID AND clinical_event_data.key = 'AGENT' AND clinical_event_data.value IN (%s)
		JOIN sample ON sample.PATIENT_ID = patient.INTERNAL_ID
		JOIN clinical_sample ON clinical_sample.INTERNAL_ID = sample.INTERNAL_ID AND clinical_sample.ATTR_ID = 'TISSUE_SOURCE_SITE' AND clinical_sample.ATTR_VALUE IN (%s);
		`
		q = fmt.Sprintf(q, "'"+strings.Join(treatment, "', '")+"'", "'"+strings.Join(sourceSite, "', '")+"'")
		results, err = clnt.Query(q)
		if err != nil {
			return nil, err
		}
	default:
		q = `
		SELECT DISTINCT cancer_study.CANCER_STUDY_ID, cancer_study.NAME, cancer_study.DESCRIPTION, cancer_study.TYPE_OF_CANCER_ID 
		FROM cancer_study;
		`
		results, err = clnt.Query(q)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(q)

	out := []types.CancerStudy{}
	for results.Next() {
		var t types.CancerStudy
		err = results.Scan(&t.ID, &t.Name, &t.Desc, &t.CancerTypeID)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}
