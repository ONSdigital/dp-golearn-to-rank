package client

import (
	"gopkg.in/olivere/elastic.v5"
	"context"
	"net/url"
	"net/http"
	"fmt"
	"github.com/pkg/errors"
)

const (
	DefaultUrl = "http://127.0.0.1:9200"
)

// Embedds an elasticsearch client
type Client struct {
	c *elastic.Client
	ctx context.Context
}

func NewClient() (*Client, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	c := &Client{
		c: client,
		ctx: ctx,
	}

	return c, nil
}

func (client *Client) BaseClient() (*elastic.Client) {
	// Returns the base elasticsearch client
	return client.c
}

func (client *Client) PerformRequest(method, path string, params url.Values, body interface{}, ignoreErrors ...int) (*elastic.Response, error) {
	return client.c.PerformRequest(client.ctx, method, path, params, body, ignoreErrors...)
}

func (client *Client) DefaultFeatureStoreExists() (bool, error) {
	return client.FeatureStoreExists("")
}

func (client *Client) FeatureStoreExists(featureStore string) (bool, error) {
	// Checks if the ltr feature store exists
	method := http.MethodGet
	path := "/_ltr"
	if featureStore != "" {
		path = fmt.Sprintf("%s/%s/", path, featureStore)
	}

	// Perform the request - note that a 404 indicates the feature store does not exist
	// and is a perfectly valid response code
	resp, err := client.PerformRequest(method, path, nil, nil, 404)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		err := errors.New(fmt.Sprintf("error checking if feature store %s exists", path))
		return false, err
	}

	exists := resp.StatusCode == http.StatusOK;

	return exists, err
}
