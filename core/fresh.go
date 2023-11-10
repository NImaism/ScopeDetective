package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/NImaism/ScopeDetective/model"
	httpx "github.com/projectdiscovery/httpx/runner"
	subfinder "github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Fresh struct {
	NotificationSystem *Messager
	Options            *Options
}

// NewFresh function Creates a new Fresh instance with the specified notification system and options.
func NewFresh(NotificationSystem *Messager, Option Options) *Fresh {
	return &Fresh{
		NotificationSystem: NotificationSystem,
		Options:            &Option,
	}
}

func (F *Fresh) Run() {
	ticker := time.NewTicker(time.Duration(F.Options.Delay) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			F.Start()
		}
	}
}

func (F *Fresh) Start() {
	if len(F.Options.WildCards) == 0 {
		syscall.Exit(0)
	}

	var wg sync.WaitGroup
	var allSubs []string
	var subMutex sync.Mutex

	subFinder := F.GenerateSubRunner()

	for _, v := range F.Options.WildCards {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			subs := F.GetSubs(domain, subFinder)

			subMutex.Lock()
			defer subMutex.Unlock()

			for _, d := range subs {
				if d != "" {
					allSubs = append(allSubs, d)
				}
			}
		}(v)
	}
	wg.Wait()

	checkedSubs := F.CheckSub(allSubs)
	savedSubs := F.OpenData(checkedSubs)

	F.CompareData(savedSubs, checkedSubs)

	F.SaveData(checkedSubs)
}

func (F *Fresh) CompareData(Saved []model.Sub, New []model.Sub) {
	savedMap := make(map[string]model.Sub)
	newMap := make(map[string]model.Sub)

	for _, s := range Saved {
		savedMap[s.URL] = s
	}

	for _, s := range New {
		newMap[s.URL] = s
	}

	for url, data := range newMap {
		saved, ok := savedMap[url]
		if !ok {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 New Service Is Up \n- Title: %s \n- Status: %t \n- Technology: %s\n- Code: %d  ```", newMap[url].Title, newMap[url].Status, newMap[url].Technology, newMap[url].Code), newMap[url].URL)
			continue
		}

		if data.Status != saved.Status {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 Change Status Detected \n- oldStatus: %t \n- newStatus: %t ```", savedMap[url].Status, newMap[url].Status), newMap[url].URL)
		}

		if HaveDifferent(saved.Code, data.Code) {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 Change Code Detected \n- oldCode: %d \n- newCode: %d ```", savedMap[url].Code, newMap[url].Code), newMap[url].URL)
		}

		if data.Words != saved.Words {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 Change Word Count Detected \n- oldCount: %d \n- newCount: %d ```", savedMap[url].Words, newMap[url].Words), newMap[url].URL)
		}

		if HaveDifferent(saved.Technology, data.Technology) {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 Change Technology Detected \n- oldTechs: %s \n- newTechs: %s ```", savedMap[url].Technology, newMap[url].Technology), newMap[url].URL)
		}

		if data.Title != saved.Title {
			F.NotificationSystem.sendSubMessage(fmt.Sprintf("```yaml\n - 💸 Change Title Detected \n- oldTitle: %s \n- newTitle: %s ```", savedMap[url].Title, newMap[url].Title), newMap[url].URL)
		}

	}

}

func (F *Fresh) CheckSub(subs []string) []model.Sub {
	var result []model.Sub
	options := F.GenerateHttpxRunner(subs, &result)

	httpxRunner, err := httpx.New(options)
	if err != nil {
		fmt.Printf("\033[31m[!] Feild To Run HTTPX \033[0m\n")
	}

	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()

	return result
}

func (F *Fresh) GetSubs(domain string, subFinder *subfinder.Runner) []string {
	output := &bytes.Buffer{}

	if err := subFinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		fmt.Printf("\033[31m[!] Feild To Enumerate %s\033[0m\n", domain)
	}

	return strings.Split(output.String(), "\n")
}

// SaveData function saves data to a JSON file for future retrieval.
func (F *Fresh) SaveData(Subs []model.Sub) {
	jsonData, err := json.Marshal(Subs)
	if err != nil {
		fmt.Println("\033[31m[!] Marshal Data Error\033[0m")
		syscall.Exit(0)
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	filePath := filepath.Join("data", "Subs.json")
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

func (F *Fresh) GenerateSubRunner() *subfinder.Runner {
	subFinderRunner, err := subfinder.NewRunner(&subfinder.Options{
		Silent:             true,
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
	})

	if err != nil {
		fmt.Println("\033[31m[!] Can't Create SubFinder Instance\033[0m")
	}

	return subFinderRunner
}

func (F *Fresh) GenerateHttpxRunner(sub []string, output *[]model.Sub) *httpx.Options {
	return &httpx.Options{
		Threads:             10,
		Silent:              true,
		Methods:             "GET",
		ExtractTitle:        true,
		TechDetect:          true,
		FollowRedirects:     true,
		FollowHostRedirects: true,
		InputTargetHost:     sub,
		OnResult: func(r httpx.Result) {
			if r.Err != nil {
				*output = append(*output, model.Sub{
					Title:      r.Input,
					URL:        r.URL,
					Technology: nil,
					Words:      0,
					Code:       nil,
					Status:     false,
				})
				return
			}

			*output = append(*output, model.Sub{
				Title:      r.Title,
				URL:        r.URL,
				Technology: r.Technologies,
				Words:      r.Words,
				Code:       intListToStringList(r.ChainStatusCodes),
				Status:     true,
			})
		},
	}
}

// OpenData function opens or creates a JSON file to store and retrieve data.
func (F *Fresh) OpenData(Data []model.Sub) []model.Sub {
	if _, err := os.Stat("data/Subs.json"); os.IsNotExist(err) {
		jsonData, err := json.Marshal(Data)
		if err != nil {
			fmt.Println("\033[31m[!] Marshal Data Error\033[0m")
			syscall.Exit(0)
		}

		if err := os.MkdirAll("data", 0755); err != nil {
			fmt.Println("\033[31m[!] Create Directory File Error\033[0m")
			syscall.Exit(0)
		}

		filePath := filepath.Join("data", "Subs.json")
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

		return Data
	} else {
		file, err := os.OpenFile("data/Subs.json", os.O_RDWR, 0644)
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

		var SavedData []model.Sub
		err = json.Unmarshal(data, &SavedData)
		if err != nil {
			fmt.Println("\033[31m[!] Error unmarshalling JSON data\033[0m")
			fmt.Println("\033[34m[+] Working on the problem\033[0m")

			err = file.Truncate(0)
			if err != nil {
				fmt.Println("\u001B[31m[!] Error deleting file content\u001B[0m")
				syscall.Exit(0)
			}

			jsonData, err := json.Marshal(Data)
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

			return Data
		} else {
			if len(SavedData) == 0 {
				syscall.Exit(0)
			}
			fmt.Println("\033[33m[+] " + "Subs Count: " + strconv.Itoa(len(SavedData)) + "\033[0m")

			return SavedData
		}

	}
}
