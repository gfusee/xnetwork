package process

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type object = map[string]interface{}

func encodeQuery(query object) (bytes.Buffer, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return bytes.Buffer{}, fmt.Errorf("error encoding query: %s", err.Error())
	}

	return buff, nil
}

func getAll() *bytes.Buffer {
	obj := object{
		"query": object{
			"match_all": object{},
		},
	}

	encoded, _ := encodeQuery(obj)

	return &encoded
}

func getWithTimestamp(start, stop int64) *bytes.Buffer {
	obj := object{
		"query": object{
			"range": object{
				"timestamp": object{
					"gte": start,
					"lte": stop,
				},
			},
		},
		"_source": true,
		"sort": []interface{}{
			object{
				"timestamp": object{
					"order": "asc",
				},
			},
		},
	}

	encoded, _ := encodeQuery(obj)

	return &encoded
}

type generalElasticResponse struct {
	Hits struct {
		Hits []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func extractSourceFromEsResponse(response generalElasticResponse) map[string]json.RawMessage {
	hits := response.Hits.Hits
	recordsMap := make(map[string]json.RawMessage, len(hits))
	for i := 0; i < len(hits); i++ {
		recordsMap[hits[i].ID] = hits[i].Source
	}

	return recordsMap
}
