# Golang + LMStudio 

This project is a proof-of-concept (and learning exercise for myself) on how to write a webserver with regular functionality and then adding integration with LM Studio for real-time LLM interactions (including tool usage specific for the project itself).

An openapi spec will be provided using the tool [Swaggo](https://github.com/swaggo/swag).

## Requirements

This project requires the following:

- [go 1.23](https://go.dev/doc/install)
- [LM Studio](https://lmstudio.ai/docs/api/server)
- [Swag](https://github.com/swaggo/swag) and [Echo-Swagger](https://github.com/swaggo/echo-swagger)
- [migrate](https://github.com/golang-migrate/migrate) for migrations

## Tips

To re-generate the openapi spec, run `swag init --dir ./cmd/server/,./internal/api/v1/` in the root directory. 

**This is important as the files provided in this repo are not being updated and only contain contact information.**