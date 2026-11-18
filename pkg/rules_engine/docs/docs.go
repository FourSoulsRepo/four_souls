// Package docs holds our rules text in rules/ (R-04) for tests: the rule
// IDs and the open questions a test may name as its source (ADR 034).
// Only tests import it; the engine and the website do not.
package docs

import (
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

// Rules is the rules text: rules/*.md.
//
//go:embed rules/*.md
var Rules embed.FS

// RuleDef finds a rule definition: "**R-STACK-03**".
var RuleDef = regexp.MustCompile(`\*\*(R-[A-Z]+-\d+)\*\*`)

var (
	// questionHead is a question's heading: "### Q-07 Choosing a dead player".
	questionHead = regexp.MustCompile(`(?m)^### (Q-\d+)\b`)
	// oldRulingRow is a row of the R-07 table: "| 4 | You can damage …".
	oldRulingRow = regexp.MustCompile(`^\| (\d+) \|`)
)

// RuleIDs returns every rule ID the rules text defines.
func RuleIDs() (map[string]bool, error) {
	files, err := fs.Glob(Rules, "rules/*.md")
	if err != nil {
		return nil, err
	}
	ids := map[string]bool{}
	for _, f := range files {
		data, err := fs.ReadFile(Rules, f)
		if err != nil {
			return nil, err
		}
		for _, m := range RuleDef.FindAllStringSubmatch(string(data), -1) {
			ids[m[1]] = true
		}
	}
	return ids, nil
}

// Questions returns the IDs of open-questions.md: each question
// ("Q-07") and each old ruling re-checked ("R-07 #4").
func Questions() (map[string]bool, error) {
	data, err := fs.ReadFile(Rules, "rules/open-questions.md")
	if err != nil {
		return nil, err
	}
	ids := map[string]bool{}
	for _, m := range questionHead.FindAllStringSubmatch(string(data), -1) {
		ids[m[1]] = true
	}
	inOld := false
	for line := range strings.Lines(string(data)) {
		if strings.HasPrefix(line, "## ") {
			inOld = strings.Contains(line, "(R-07)")
			continue
		}
		if m := oldRulingRow.FindStringSubmatch(line); inOld && m != nil {
			ids[fmt.Sprintf("R-07 #%s", m[1])] = true
		}
	}
	return ids, nil
}
