package core

import (
	"encoding/json"
	"fmt"
	"github.com/NImaism/ScopeDetective/model"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type System struct {
	NotificationSystem *Messager
	Options            *Options
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
	ticker := time.NewTicker(time.Duration(s.Options.Delay+3) * time.Minute)
	defer ticker.Stop()

	s.NotificationSystem.sendLog("```yaml\n - 📡 Detective Initiates HackerOne Monitoring ! ```")
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
func (s *System) calculateData(data []byte) []model.Message {
	fmt.Println("\033[32m[+] System Started !\033[0m")
	s.NotificationSystem.sendLog("```yaml\n - 🔍 Detective Begins Document Inspection ! ```")

	var Data []model.JsonData
	var wg sync.WaitGroup

	CollectedMessage := model.StoredData{Data: []model.Message{}, Subs: []string{}}
	_ = json.Unmarshal(data, &Data)
	SavedData := s.openData(Data)

	for _, Pr := range Data {
		wg.Add(1)
		go func(program model.JsonData) {
			CollectedMessage.Mutex.Lock()
			for _, item := range program.Targets.InScope {
				if item.AssetType == "URL" && item.EligibleForSubmission {
					CollectedMessage.Subs = append(CollectedMessage.Subs, item.AssetIdentifier)
					if !Contains(SavedData, item.AssetIdentifier) && (s.Options.Vdp || item.EligibleForBounty) {
						CollectedMessage.Data = append(CollectedMessage.Data, model.Message{
							SubDomain:   item.AssetIdentifier,
							Owner:       program.Name,
							Url:         program.URL,
							MaxSeverity: item.MaxSeverity,
						})
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
		s.NotificationSystem.sendLog("```yaml\n - 📜 Detective Discovers No Pertinent Evidence !```")
		fmt.Println("\u001B[35m[-] No Change \u001B[0m")
	} else {

		s.NotificationSystem.sendLog("```yaml\n - 🔮 Detective Makes Significant Discovery !```")
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
func (s *System) openData(Data []model.JsonData) []string {
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
			if len(SavedSubs) == 0 {
				syscall.Exit(0)
			}
			fmt.Println("\033[33m[+] " + "Count: " + strconv.Itoa(len(SavedSubs)) + "\033[0m")
			return SavedSubs
		}

	}
}
