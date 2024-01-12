gen-go:
	get-graphql-schema https://anchor.zeet.co/graphql > schema_0.graphql
	get-graphql-schema https://anchor.zeet.co/v1/graphql > schema_1.graphql
	go generate ./...

gen:
	go run github.com/Khan/genqlient genqlient_0.yaml
	go run github.com/Khan/genqlient genqlient_1.yaml
