package types

type CancerStudy struct {
	ID                    string `json:"studyId,omitempty"`
	CancerStudyIdentifier string `json:"cancerStudyIdentified,omitempty"`
	CancerTypeID          string `json:"cancerTypeId,omitempty"`
	Name                  string `json:"name,omitempty"`
	Desc                  string `json:"description,omitempty"`
	PublicStudy           bool   `json:"publicStudy,omitempty"`
	Pmid                  string `json:"pmid,omitempty"`
	Citation              string `json:"citation,omitempty"`
	Groups                string `json:"groups,omitempty"`
	Status                int    `json:"status,omitempty"`
	ImportDate            string `json:"importDate,omitempty"`
	ReferenceGenome       string `json:"referenceGenome,omitempty"`
}
