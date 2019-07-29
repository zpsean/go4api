/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package keyword

import (
    // "fmt"
    // "time"
    // "os"
    // "sort"
    "strings"
)

// for keyword to ts / tc execution, there are two options:
// option 1: convert to ts / tc (i.e. temp files) format then execution
// option 2: mapping the keywords with existing ts/tc (i.e. treat the ts/tc as library, but can accept params)
//
// here use option 2, steps are:
// (1). scan the keyword files, get all the testcases and their variables
// (2). lookup ts/tc, generate the cases to be executed
// (3). reporting
// i.e. go4api -run -K -kw xx/xx/*.keyword -tc xxxx/ -tsuite xxxxx -jsFuncs xxx -r xxx -tr xxx

func (gKw *GKeyWord) PopulateSettingsOriginalContent (startLine int, endLine int, lines []string) {
	var originalContentTemp []string

	if startLine != endLine {
		for i := startLine + 1; i <= endLine; i ++ {
			originalContentTemp = append(originalContentTemp, lines[i])
		}
	}

	settings := &Settings {
        StartLine:       startLine,
        EndLine:         endLine,
        OriginalContent: originalContentTemp,
    }

    gKw.Settings = settings
}

func (gKw *GKeyWord) PopulateTestCasesOriginalContent (startLine int, endLine int, lines []string) {
	var originalContentTemp []string

	if startLine != endLine {
		for i := startLine + 1; i <= endLine; i ++ {
			originalContentTemp = append(originalContentTemp, lines[i])
		}
	}

	testCases := &TestCases {
        StartLine:       startLine,
        EndLine:         endLine,
        OriginalContent: originalContentTemp,
    }

    gKw.TestCases = testCases
	
}

func (gKw *GKeyWord) PopulateVariablesOriginalContent (startLine int, endLine int, lines []string) {
	var originalContentTemp []string

	if startLine != endLine {
		for i := startLine + 1; i <= endLine; i ++ {
			originalContentTemp = append(originalContentTemp, lines[i])
		}
	}

	variables := &Variables {
        StartLine:       startLine,
        EndLine:         endLine,
        OriginalContent: originalContentTemp,
    }

    gKw.Variables = variables
	
}

func (gKw *GKeyWord) ParseSettingsOriginalContent () {
	for _, line := range gKw.Settings.OriginalContent {
		switch {
		case strings.HasPrefix(line, "ID"):
			str := strings.TrimLeft(line, "ID:")
			str = strings.TrimSpace(str)
			gKw.Settings.ID = str
		case strings.HasPrefix(line, "TestSuites"):
			str := strings.TrimLeft(line, "TestSuites:")
			str = strings.TrimSpace(str)
			gKw.Settings.TestSuitePaths = strings.Split(str, ",")
		case strings.HasPrefix(line, "BasicTestCases"):
			str := strings.TrimLeft(line, "BasicTestCases:")
			str = strings.TrimSpace(str)
			gKw.Settings.BasicTestCasePaths = strings.Split(str, ",")
		case strings.HasPrefix(line, "JsFuncs"):
			str := strings.TrimLeft(line, "JsFuncs:")
			str = strings.TrimSpace(str)
			gKw.Settings.JsFuncPaths = strings.Split(str, ",")
		}
	}
}

func (gKw *GKeyWord) ParseTestCasesOriginalContent () {
	var parsedTestCases []*KWTestCase

	for i, line := range gKw.TestCases.OriginalContent {
		str := strings.TrimSpace(line)

		if len(str) > 0 {
			var parsedTestCase []string

			fields := strings.Fields(str)
			kwTcName := strings.Join(fields, "-")
			
			// note, here for name only, arg handle to be added later 
			parsedTestCase = append(parsedTestCase, kwTcName)

			kwTestCase := &KWTestCase {
				OriginalLine:     i + gKw.TestCases.StartLine + 1,
				OriginalTestCase: line,
				KWTestCaseName:	  kwTcName,
				ParsedTestCase:   parsedTestCase,
			}

			parsedTestCases = append(parsedTestCases, kwTestCase)
		}
	}
	
	gKw.TestCases.ParsedTestCases = parsedTestCases
}

func (gKw *GKeyWord) ParseVariablesOriginalContent () {
	
}
