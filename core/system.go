package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

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

type System struct {
	NotificationSystem *Messager
	Options            *Options
}

type StoredData struct {
	Data  []Message
	Subs  []string
	Mutex sync.Mutex
}

type Message struct {
	SubDomain string
	Owner     string
	Url       string
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

// New function Creates a new system instance with the specified notification system and options.
func New(NotificationSystem *Messager, Option Options) *System {
	return &System{
		NotificationSystem: NotificationSystem,
		Options:            &Option,
	}
}

// Run function Runs the system in a loop, periodically fetching data and sending notifications based on calculations.
func (s *System) Run() {
	ticker := time.NewTicker(time.Duration(s.Options.Delay) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, v := range s.calculateData(s.Pull()) {
				s.NotificationSystem.sendMessage(v)
				time.Sleep(2 * time.Second)
			}
		}
	}
}

// Pull function pulls the content of the "hackerone_data.json" file from the specified URL and returns it as a byte array.
func (s *System) Pull() []byte {
	resp, err := http.Get("https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/hackerone_data.json")
	if err != nil {
		fmt.Println("\033[31m[!] Network Error\033[0m")
		syscall.Exit(0)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("\033[31m[!] Read File Error\033[0m")
		syscall.Exit(0)
	}

	return data
}

// CalculateData function processes a byte slice of data to calculate values using concurrent processing, goroutines, and a wait group. It prints out messages to indicate progress and results.
func (s *System) calculateData(data []byte) []Message {
	fmt.Println("\033[32m[+] System Started !\033[0m")

	var Data []JsonData
	var wg sync.WaitGroup

	CollectedMessage := StoredData{Data: []Message{}, Subs: []string{}}
	_ = json.Unmarshal(data, &Data)
	SavedData := s.openData(Data)

	for _, Pr := range Data {
		wg.Add(1)
		go func(program JsonData) {
			CollectedMessage.Mutex.Lock()
			for _, item := range program.Targets.InScope {
				if item.AssetType == "URL" {
					CollectedMessage.Subs = append(CollectedMessage.Subs, item.AssetIdentifier)
					if item.EligibleForBounty {
						if !s.Contains(SavedData, item.AssetIdentifier) {
							CollectedMessage.Data = append(CollectedMessage.Data, Message{
								SubDomain: item.AssetIdentifier,
								Owner:     program.Name,
								Url:       program.URL,
							})
						}
					}
				}
			}
			CollectedMessage.Mutex.Unlock()
			wg.Done()
		}(Pr)
	}

	wg.Wait()

	s.saveData(CollectedMessage.Subs)
	if len(CollectedMessage.Data) == 0 {
		fmt.Println("\u001B[35m[-] No Change \u001B[0m")
	} else {
		fmt.Printf("\u001B[35m[+] %d Change \u001B[0m\n", len(CollectedMessage.Data))
	}

	return CollectedMessage.Data
}

// SaveData function saves data to a JSON file for future retrieval.
func (s *System) saveData(Subs []string) {
	jsonData, err := json.Marshal(Subs)
	if err != nil {
		fmt.Println("\033[31m[!] Marshal Data Error\033[0m")
		syscall.Exit(0)
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	filePath := filepath.Join("data", "Scopes.json")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("\033[31m[!] Save Pulled File Error\033[0m")
		syscall.Exit(0)
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		fmt.Println("\033[31m[!] Save Pulled File Error\033[0m")
		syscall.Exit(0)
	}
}

// OpenData function opens or creates a JSON file to store and retrieve data.
func (s *System) openData(Data []JsonData) []string {
	if _, err := os.Stat("data/Scopes.json"); os.IsNotExist(err) {
		var Subs []string
		for _, program := range Data {
			for _, item := range program.Targets.InScope {
				if item.AssetType == "URL" {
					Subs = append(Subs, item.AssetIdentifier)
				}
			}
		}

		jsonData, err := json.Marshal(Subs)
		if err != nil {
			fmt.Println("\033[31m[!] Marshal Data Error\033[0m")
			syscall.Exit(0)
		}

		if err := os.MkdirAll("data", 0755); err != nil {
			fmt.Println("\033[31m[!] Create Directory File Error\033[0m")
			syscall.Exit(0)
		}

		filePath := filepath.Join("data", "Scopes.json")
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("\033[31m[!] Save Pulled File Error\033[0m")
			syscall.Exit(0)

		}
		defer file.Close()

		if _, err := file.Write(jsonData); err != nil {
			fmt.Println("\033[31m[!] Write Pulled File Error\033[0m")
			syscall.Exit(0)
		}

		return Subs
	} else {
		file, err := os.OpenFile("data/Scopes.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("\033[31m[!] Open Pulled File Error\033[0m")
			syscall.Exit(0)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("\033[31m[!] Error reading file\033[0m")
			syscall.Exit(0)
		}

		var SavedSubs []string
		err = json.Unmarshal(data, &SavedSubs)
		if err != nil {
			fmt.Println("\033[31m[!] Error unmarshalling JSON data\033[0m")
			fmt.Println("\033[34m[+] Working on the problem\033[0m")

			err = file.Truncate(0)
			if err != nil {
				fmt.Println("\u001B[31m[!] Error deleting file content\u001B[0m")
				syscall.Exit(0)
			}

			var Subs []string
			for _, program := range Data {
				for _, item := range program.Targets.InScope {
					if item.AssetType == "URL" {
						Subs = append(Subs, item.AssetIdentifier)
					}
				}
			}

			jsonData, err := json.Marshal(Subs)
			if err != nil {
				fmt.Println("\033[31m[!] Marshal Data Error\033[0m")
				syscall.Exit(0)
			}

			_, err = file.Write(jsonData)
			if err != nil {
				fmt.Println("\033[31m[!] Write Pulled File Error\033[0m")
				syscall.Exit(0)
			}

			fmt.Println("\033[34m[+] Problem has been solved\033[0m")
			return Subs
		} else {
			fmt.Println("\033[33m[+] " + "Count: " + strconv.Itoa(len(SavedSubs)) + "\033[0m")
			return SavedSubs
		}

	}
}

// Contains function checks if an item exists in a list.
func (s *System) Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}
