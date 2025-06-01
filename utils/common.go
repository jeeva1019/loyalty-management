package utils

import (
	"encoding/json"
	"fmt"
	"loyality_points/helpers"
	"net/http"
	"strconv"
)

type RspStruct struct {
	Status string `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Result any    `json:"result,omitempty"`
}

func ResponseWriter(debug *helpers.HelperStruct, resp http.ResponseWriter, result any) {
	debug.Info("ResponseWriter(+)")

	respStr, err := json.Marshal(result)
	if err != nil {
		debug.Error("error occured at marshal", err)
		return
	}

	fmt.Fprintln(resp, string(respStr))
	debug.Info("ResponseWriter(-)")
}

func ResponseConstructor(status, msg string, val any) string {
	var finalRes RspStruct

	finalRes.Status = status
	finalRes.Msg = msg
	finalRes.Result = val

	bodyStr, err := json.Marshal(finalRes)
	if err != nil {
		return ""
	}

	return string(bodyStr)
}

func GetPaginationValue(start, end string) (int, int, int) {
	page := 1
	pageSize := 20

	if start != "" {
		if p, err := strconv.Atoi(start); err == nil && p > 0 {
			page = p
		}
	}

	if end != "" {
		if ps, err := strconv.Atoi(end); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	return page, pageSize, offset
}
