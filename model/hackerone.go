package model

import "sync"

type JsonData struct {
	AllowsBountySplitting             bool   `json:"allows_bounty_splitting"`
	AverageTimeToBountyAwarded        *int   `json:"average_time_to_bounty_awarded"`
	AverageTimeToFirstProgramResponse *int   `json:"average_time_to_first_program_response"`
	AverageTimeToReportResolved       *int   `json:"average_time_to_report_resolved"`
	Handle                            string `json:"handle"`
	ID                                int    `json:"id"`
	ManagedProgram                    bool   `json:"managed_program"`
	Name                              string `json:"name"`
	OffersBounties                    bool   `json:"offers_bounties"`
	OffersSwag                        bool   `json:"offers_swag"`
	ResponseEfficiencyPercentage      int    `json:"response_efficiency_percentage"`
	SubmissionState                   string `json:"submission_state"`
	URL                               string `json:"url"`
	Website                           string `json:"website"`
	Targets                           struct {
		InScope []Scope `json:"in_scope"`
	} `json:"targets"`
}

type StoredData struct {
	Data  []Message
	Subs  []string
	Mutex sync.Mutex
}

type Message struct {
	SubDomain   string
	Owner       string
	Url         string
	MaxSeverity string
}

type Scope struct {
	AssetIdentifier            string `json:"asset_identifier"`
	AssetType                  string `json:"asset_type"`
	AvailabilityRequirement    string `json:"availability_requirement"`
	ConfidentialityRequirement string `json:"confidentiality_requirement"`
	EligibleForBounty          bool   `json:"eligible_for_bounty"`
	EligibleForSubmission      bool   `json:"eligible_for_submission"`
	Instruction                string `json:"instruction"`
	IntegrityRequirement       string `json:"integrity_requirement"`
	MaxSeverity                string `json:"max_severity"`
}
