package SuricataMatch

import (
	"GoFender/Utils"
	"bufio"
	ac "github.com/BobuSumisu/aho-corasick"
	"github.com/google/gonids"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type RuleInfo struct {
	Regex          string
	BackRegex      string
	ContentPattern []byte
	Msg            string
}

var (
	RuleSet []RuleInfo
	ACTrie  *ac.Trie
)

func GetRuleFiles(filePath string) (Files []string, err error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Files, err
	}
	rulesStat, _ := os.Stat(filePath)
	switch mode := rulesStat.Mode(); {
	case mode.IsDir():
		err = filepath.Walk(filePath, func(filePath string, fileObj os.FileInfo, err error) error {
			rulesObj, err := os.Open(filePath)
			defer rulesObj.Close()
			if err == nil && path.Ext(rulesObj.Name()) == ".rules" {
				Files = append(Files, filePath)
			}
			return nil
		})
	case mode.IsRegular():
		rulesObj, err := os.Open(filePath)
		defer rulesObj.Close()
		if err == nil && path.Ext(rulesObj.Name()) == ".rules" {
			Files = append(Files, filePath)
		}
	}
	return Files, err
}

func CheckRules(rule RuleInfo) bool {
	if len(rule.Regex) < 5 || len(rule.BackRegex) < 5 || rule.ContentPattern == nil || len(rule.Msg) < 1 {
		return false
	}
	return true
}

func RulesToMatch(rule *gonids.Rule) []RuleInfo {
	rein := RuleInfo{}
	reins := make([]RuleInfo, 0)
	if rule != nil {
		rein.Msg = rule.Description
		rein.Regex = strings.Replace(strings.Replace(rule.RE(), " ", "", -1), "\n", "", -1)
		if rule.Contents() != nil {
			for _, content := range rule.Contents() {
				if content.Pattern != nil {
					rein.ContentPattern = content.Pattern
					rein.BackRegex = strings.Replace(strings.Replace(content.ToRegexp(), " ", "", -1), "\n", "", -1)
				}
				reins = append(reins, rein)
			}
		}
	}
	return reins
}

func RulesParse(rulesPath string) []RuleInfo {
	ruleSlice := make([]RuleInfo, 0)
	rulesFiles, err := GetRuleFiles(rulesPath)
	if err != nil {
		log.Fatal("Can't get the snort rules files:", err)
	}
	for _, ruleFile := range rulesFiles {
		file, err := os.Open(ruleFile)
		if err != nil {
			log.Printf("Cannot open text file: %s, err: [%v]", ruleFile, err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			rule, err := gonids.ParseRule(line)
			if err == nil && len(rule.String()) > 0 {
				ruletmp := RulesToMatch(rule)
				for _, rule := range ruletmp {
					if CheckRules(rule) {
						ruleSlice = append(ruleSlice, rule)
					}
					continue
				}
			}
		}
	}
	return ruleSlice
}

func PacketMatch(packet *Utils.CommonPacket, rules []RuleInfo) (float64, bool, []string) {
	rulesChan := make(chan RuleInfo, len(rules))
	for _, rule := range rules {
		rulesChan <- rule
	}
	MsgChan := make(chan string, 1)
	for i := 0; i < len(rules); i++ {
		go func(ch chan RuleInfo) {
			if rulesChan == nil {
				return
			}
			rule := <-ch
			re, err := regexp.Compile(rule.Regex)
			if err != nil {
				re, err = regexp.Compile(rule.BackRegex)
				if err != nil {
					log.Println("Can't compile the regex:", err)
					return
				}
			}
			if tmp := re.Find(packet.ComPacketData); tmp != nil {
				MsgChan <- rule.Msg
			}
		}(rulesChan)
	}
	re := ACTrie.Match(packet.ComPacketData)
	msgslice := make([]string, 0)
	if len(re) > 2 && len(MsgChan) > 0 {
		for i := 0; i < len(MsgChan); i++ {
			msg := <-MsgChan
			if len(msg) > 1 {
				msgslice = append(msgslice, msg)
			}
		}
		return calcTotalScore(removeZeroContent(re)), true, msgslice
	} else {
		return 0, false, nil
	}
}
