module github.com/mondracode/ambrosia-atlas-api

go 1.19

require (
	github.com/99designs/gqlgen v0.17.20
	github.com/vektah/gqlparser/v2 v2.5.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
)

replace (
	github.com/mondracode/ambrosia-atlas-api/graph => ./graph
	github.com/mondracode/ambrosia-atlas-api/graph/generated => ./graph/generated
	github.com/mondracode/ambrosia-atlas-api/graph/model => ./graph/model
)
