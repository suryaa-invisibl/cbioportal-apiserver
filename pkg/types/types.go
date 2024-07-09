package types

type CancerStudy struct {
	ID                          string `json:"studyId,omitempty"`
	CancerStudyIdentifier       string `json:"cancerStudyIdentified,omitempty"`
	CancerTypeID                string `json:"cancerTypeId,omitempty"`
	Name                        string `json:"name,omitempty"`
	Desc                        string `json:"description,omitempty"`
	PublicStudy                 bool   `json:"publicStudy,omitempty"`
	Pmid                        string `json:"pmid,omitempty"`
	Citation                    string `json:"citation,omitempty"`
	Groups                      string `json:"groups,omitempty"`
	Status                      int    `json:"status,omitempty"`
	ImportDate                  string `json:"importDate,omitempty"`
	ReferenceGenome             string `json:"referenceGenome,omitempty"`
	AllSampleCount              int    `json:"allSampleCount,omitempty"`
	SequencedSampleCount        int    `json:"sequencedSampleCount,omitempty"`
	CnaSampleCount              int    `json:"cnaSampleCount,omitempty"`
	MrnaRnaSeqSampleCount       int    `json:"mrnaRnaSeqSampleCount,omitempty"`
	MrnaRnaSeqV2SampleCount     int    `json:"mrnaRnaSeqV2SampleCount,omitempty"`
	MrnaMicroarraySampleCount   int    `json:"mrnaMicroarraySampleCount,omitempty"`
	MiRnaSampleCount            int    `json:"miRnaSampleCount,omitempty"`
	MethylationHm27SampleCount  int    `json:"methylationHm27SampleCount,omitempty"`
	RppaSampleCount             int    `json:"rppaSampleCount,omitempty"`
	MassSpectrometrySampleCount int    `json:"massSpectrometrySampleCount,omitempty"`
	CompleteSampleCount         int    `json:"completeSampleCount,omitempty"`
	ReadPermission              bool   `json:"readPermission,omitempty"`
	TreatmentCount              int    `json:"treatmentCount,omitempty"`
	StructuralVariantCount      int    `json:"structuralVariantCount,omitempty"`
	CancerType                  struct {
		Name           string `json:"name,omitempty"`
		DedicatedColor string `json:"dedicatedColor,omitempty"`
		ShortName      string `json:"shortName,omitempty"`
		Parent         string `json:"parent,omitempty"`
		CancerTypeID   string `json:"cancerTypeId,omitempty"`
	} `json:"cancerType,omitempty"`
}
