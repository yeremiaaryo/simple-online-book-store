#!/bin/bash

export POSTGRESQL_URL='postgres://admin:admin@localhost:5432/gotu?sslmode=disable'
run:
	@ go run cmd/main/main.go

migrate-create:
	@ migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@ migrate -database ${POSTGRESQL_URL} -path scripts/migrations up

migrate-down:
	@ migrate -database ${POSTGRESQL_URL} -path scripts/migrations down

mock:
	@ go generate -x ./...