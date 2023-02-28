package mocks

// ElasticClientStub -
type ElasticClientStub struct {
	GetMultipleCalled func(index string, requests []string) ([]byte, error)
}

// GetMultiple -
func (ecs *ElasticClientStub) GetMultiple(index string, requests []string) ([]byte, error) {
	if ecs.GetMultipleCalled != nil {
		return ecs.GetMultipleCalled(index, requests)
	}

	return nil, nil
}
