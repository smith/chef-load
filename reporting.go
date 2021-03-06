package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-chef/chef"
	uuid "github.com/satori/go.uuid"
)

const rubyDateTime = "2006-01-02 15:04:05 -0700"

func reportingRunStart(nodeClient chef.Client, nodeName string, runUUID uuid.UUID, startTime time.Time) int {
	startRunBody := map[string]string{
		"action":     "start",
		"run_id":     runUUID.String(),
		"start_time": startTime.Format(rubyDateTime),
	}
	data, err := chef.JSONReader(startRunBody)
	if err != nil {
		fmt.Println(err)
	}

	req, err := nodeClient.NewRequest("POST", "reports/nodes/"+nodeName+"/runs", data)
	req.Header.Set("X-Ops-Reporting-Protocol-Version", "0.1.0")
	res, err := nodeClient.Do(req, nil)
	if err != nil && res.StatusCode != 404 {
		// can't print res here if it is nil
		// fmt.Println(res.StatusCode)
		fmt.Println(err)
		return res.StatusCode
	}
	defer res.Body.Close()
	return res.StatusCode
}

func reportingRunStop(nodeClient chef.Client, nodeName string, runUUID uuid.UUID, startTime time.Time, endTime time.Time, rl runList) int {
	endRunBody := map[string]interface{}{
		"action":          "end",
		"data":            map[string]interface{}{},
		"end_time":        endTime.Format(rubyDateTime),
		"resources":       []interface{}{},
		"run_list":        `["` + strings.Join(rl.toStringSlice(), `","`) + `"]`,
		"start_time":      startTime.Format(rubyDateTime),
		"status":          "success",
		"total_res_count": "0",
	}
	data, err := chef.JSONReader(endRunBody)
	if err != nil {
		fmt.Println(err)
	}

	req, err := nodeClient.NewRequest("POST", "reports/nodes/"+nodeName+"/runs/"+runUUID.String(), data)
	req.Header.Set("X-Ops-Reporting-Protocol-Version", "0.1.0")
	res, err := nodeClient.Do(req, nil)
	if err != nil {
		// can't print res here if it is nil
		// fmt.Println(res.StatusCode)
		fmt.Println(err)
		return res.StatusCode
	}
	defer res.Body.Close()
	return res.StatusCode
}
