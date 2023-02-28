package elastic

import (
	"fmt"
	"net/http"
)

const numOfErrorsToExtractBulkResponse = 5

// bulkRequestResponse defines the structure of a bulk request response
type bulkRequestResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			Status int `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

func extractErrorFromBulkResponse(response *bulkRequestResponse) error {
	count := 0
	errorsString := ""
	for _, item := range response.Items {
		if item.Index.Status < http.StatusBadRequest {
			continue
		}

		count++
		errorsString += fmt.Sprintf("{ status code: %d, error type: %s, reason: %s }\n", item.Index.Status, item.Index.Error.Type, item.Index.Error.Reason)

		if count == numOfErrorsToExtractBulkResponse {
			break
		}
	}

	return fmt.Errorf("%s", errorsString)
}
