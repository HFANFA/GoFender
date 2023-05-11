package YaraMatch

import (
	"github.com/hillu/go-yara/v4"
	"github.com/orcaman/concurrent-map"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Scanner struct {
	Rules *yara.Rules
}

type FileScanResult struct {
	FileName string
	Matches  []yara.MatchRule
}

type FileResult struct {
	Filename    string
	Namespace   string
	Rule        string
	Description string
}

func LoadRules(rulesData string) (*yara.Rules, error) {
	rules, err := yara.LoadRules(rulesData)
	return rules, err
}

func (s *Scanner) ScanFile(filename string) (error, *FileScanResult) {
	sl, _ := yara.NewScanner(s.Rules)
	var m yara.MatchRules
	err := sl.SetCallback(&m).ScanFile(filename)
	err = s.Rules.ScanFile(filename, 0, 10, nil)
	result := &FileScanResult{FileName: filename, Matches: nil}
	return err, result
}

func (s *Scanner) ScanFiles(filename string) {
	files, err := GetFiles(filename)
	if err == nil {
		var wg sync.WaitGroup
		wg.Add(len(files))
		for _, f := range files {
			//SaveFileResult(s.ScanFile(f))
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				SaveFileResult(s.ScanFile(filename))
			}(f)
			waitTimeout(&wg, 60)
		}
	}
}
func SaveFileResult(err error, result *FileScanResult) {
	FileResultMap := cmap.New()
	if err == nil && len(result.Matches) > 0 {
		FileResultMap.Set(result.FileName, result)
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

func GetFiles(filePath string) (Files []string, err error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Files, err
	}
	rulesStat, _ := os.Stat(filePath)
	switch mode := rulesStat.Mode(); {
	case mode.IsDir():
		err = filepath.Walk(filePath, func(filePath string, fileObj os.FileInfo, err error) error {
			rulesObj, err := os.Open(filePath)
			defer rulesObj.Close()
			if err == nil {
				Files = append(Files, filePath)
			}
			return nil
		})
	case mode.IsRegular():
		rulesObj, err := os.Open(filePath)
		defer rulesObj.Close()
		if err == nil {
			Files = append(Files, filePath)
		}
	}
	return Files, err
}
