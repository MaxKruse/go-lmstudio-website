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

To re-generate the openapi spec, run `swag init --dir ./cmd/server/,./internal/api/v1/controllers,./internal/models/dtos/,./internal/models/dtos/request_dtos/` in the root directory. 

**This is important as the files provided in this repo are not being updated and only contain contact information.**

To create a new migration, first install migrate:

> go install github.com/golang-migrate/migrate/v4/cmd/migrate@4.18.1

Then run:

> migrate create -ext sql -dir ./migrations/ -seq <name>

This will create a new migration file in the `./migrations/` directory.

## Suggestions

For LLMs to use, in my personal testing with an RTX 4070 12GB, I conclude the following models to work well: 

- granite-3.1-8b-instruct
- internlm2_5-20b-chat

Other models might work, but flaky tool_chain usage has been observed.