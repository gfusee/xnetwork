package reader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// GetElasticTemplatesAndPolicies will return elastic templates and policies
// TODO implement policies when will start to use it again
func GetElasticTemplatesAndPolicies(path string, indexes []string) (map[string]*bytes.Buffer, map[string]*bytes.Buffer, error) {
	indexTemplates := make(map[string]*bytes.Buffer)
	indexPolicies := make(map[string]*bytes.Buffer)
	var err error

	for _, index := range indexes {
		indexTemplates[index], err = getTemplateByIndex(path, index)
		if err != nil {
			return nil, nil, err
		}
	}

	return indexTemplates, indexPolicies, nil
}

func getTemplateByIndex(path string, index string) (*bytes.Buffer, error) {
	indexTemplate := &bytes.Buffer{}

	fileName := fmt.Sprintf("%s.json", index)
	filePath := filepath.Join(path, fileName)
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("getTemplateByIndex: %w, path %s, error %s", err, filePath, err.Error())
	}

	indexTemplate.Grow(len(fileBytes))
	_, err = indexTemplate.Write(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("getTemplateByIndex: %w, path %s, error %s", err, filePath, err.Error())
	}

	return indexTemplate, nil
}
