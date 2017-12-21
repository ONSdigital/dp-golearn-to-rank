package dp_golearn_to_rank

import (
	"net/http"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"errors"
)

const (
	LearnToRankApi = "/_ltr"
)

type FeatureService struct {
	client *Client
}

func NewFeatureService(client *Client) (*FeatureService) {
	f := FeatureService{
		client: client,
	}
	return &f
}

func (featureService *FeatureService) DefaultFeatureStoreExists() (bool, error) {
	return featureService.FeatureStoreExists("")
}

func (featureService *FeatureService) FeatureStoreExists(featureStore string) (bool, error) {
	// Checks if the ltr feature store exists
	method := http.MethodGet
	path := LearnToRankApi
	if featureStore != "" {
		path = path + "/" + featureStore
	}

	// Perform the request - note that a 404 indicates the feature store does not exist
	// and is a perfectly valid response code
	resp, err := featureService.client.PerformRequest(method, path, nil, nil, 404)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		err := errors.New(fmt.Sprintf("error checking if feature store %s exists", path))
		return false, err
	}

	exists := resp.StatusCode == http.StatusOK;

	return exists, err
}

func (featureService *FeatureService) CreateFeatureStore(featureStore string) (*elastic.Response, error) {
	// Creates a feature store
	method := http.MethodPut
	path := LearnToRankApi + "/" + featureStore

	resp, err := featureService.client.PerformRequest(method, path, nil, nil)

	return resp, err
}
