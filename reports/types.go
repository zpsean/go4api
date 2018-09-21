/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
	// "go4api/lib/testcase"
)

type TcStats struct {
	StatusStats map[string]map[string]int
	StatusStatsPercentage map[string]map[string]float64
}

type MutationStats struct {
	HttpUrl string
	HttpMethod string
	MutationPart string
	MutationRule string
	HttpStatus int
	TestStatus string 
	Count int
}

type MutationDetails struct {
	HttpUrl string
	HttpMethod string
	MutationPart string
	MutationRule string
	HttpStatus int
	TestStatus string 
	MutationMessage string
}