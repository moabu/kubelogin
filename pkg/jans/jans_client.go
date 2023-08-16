package jans

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"net/http"
	"net/http/httputil"
)

var (
	ErrorBadRequest = fmt.Errorf("bad request")
	ErrorNotFound   = fmt.Errorf("not found")
)

// requestParams is used as a conveneince struct to pass parameters to the
// request method.
type requestParams struct {
	method      string
	path        string
	accept      string
	contentType string
	token       string
	payload     []byte
	resp        any
	returnRaw   bool
}

// Client is the client via which we can interact with all
// necessary Jans APIs.
type Client struct {
	host          string
	clientId      string
	clientSecret  string
	skipTLSVerify bool
}

// NewClient creates a new client, which will connect to a server
// at the provided host, using the given credentials.
func NewClient(host, clientId, clientSecret string) (*Client, error) {
	return &Client{
		host:          host,
		clientId:      clientId,
		clientSecret:  clientSecret,
		skipTLSVerify: false,
	}, nil
}

// NewInsecureClient creates a new client, which will connect to a server
// at the provided host, using the given credentials. Unlike NewClient,
// this client will skip TLS verification. This should only be used for
// development and testing purposes.
func NewInsecureClient(host, clientId, clientSecret string) (*Client, error) {
	return &Client{
		host:          host,
		clientId:      clientId,
		clientSecret:  clientSecret,
		skipTLSVerify: true,
	}, nil
}

// request performs an HTTP request of the requested method to the given path.
// The token is used as authorization header. If the request entity is not nil,
// it is marshaled into JSON and used as request body. If the response value
// is not nil, the response data is unmarshaled into it. The response value
// has to be of a pointer type.
func (c *Client) request(ctx context.Context, params requestParams) error {

	if c.host == "" {
		return fmt.Errorf("host is not set")
	}

	if params.path == "" {
		return fmt.Errorf("no request path provided")
	}

	url := fmt.Sprintf("%s%s", c.host, params.path)

	req, err := http.NewRequestWithContext(ctx, params.method, url, bytes.NewReader(params.payload))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add("Accept", params.accept)
	req.Header.Add("Content-Type", params.contentType)

	if params.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", params.token))
	} else {
		req.SetBasicAuth(c.clientId, c.clientSecret)
	}

	tr := &http.Transport{}
	if c.skipTLSVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{Transport: tr}

	// reqDump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	return fmt.Errorf("could not dump request: %w", err)
	// }
	// fmt.Printf("Request:\n%s\n", string(reqDump))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Errorf("could not dump response: %w", err)
	}
	fmt.Printf("Reponse:\n%s\n", string(respDump))

	if resp.StatusCode == 400 {
		// try to read error message
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrorBadRequest
		}

		return fmt.Errorf("%w: %v", ErrorBadRequest, string(data))
	}

	if resp.StatusCode == 404 {
		return ErrorNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// try to read error message
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrorBadRequest
		}

		return fmt.Errorf("did not get correct response code (%v): %v", resp.Status, string(data))
	}

	if params.resp == nil {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	if len(data) == 0 || params.resp == nil {
		return nil
	}

	// if json.Valid(data) {
	// 	return fmt.Errorf("response is not valid json")
	// }

	if params.returnRaw {
		if reflect.ValueOf(params.resp).Kind() != reflect.Ptr {
			return fmt.Errorf("response destination is not a pointer")
		}
		if reflect.ValueOf(params.resp).Elem().Kind() != reflect.String {
			return fmt.Errorf("response destination is not a string pointer")
		}
		*params.resp.(*string) = string(data)
		return nil
	}

	if err = json.Unmarshal(data, params.resp); err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}

	return nil
}
