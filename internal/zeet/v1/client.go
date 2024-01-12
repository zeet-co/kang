package v1

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"k8s.io/client-go/transport"
)

type Client struct {
	gql  graphql.Client
	http *http.Client
	path string
}

func New(token string) *Client {
	path := "https://anchor.zeet.co/v1/graphql"

	httpClient := newHTTPClient(token)
	gqlClient := newGraphQLClient(path, httpClient)

	return &Client{path: path, gql: gqlClient, http: httpClient}
}

func newHTTPClient(token string) *http.Client {
	tp := http.DefaultTransport

	return &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, tp),
	}
}

func newGraphQLClient(path string, httpClient *http.Client) graphql.Client {
	return graphql.NewClient(path, httpClient)
}
