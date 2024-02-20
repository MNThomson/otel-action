package main

import (
	"fmt"
	"os"
	"strings"
)

type configType struct {
	otelEndpoint string
	otelHeaders  map[string]string
	serviceName  string
	runID        string
	repository   string
	owner        string
	repo         string
	githubToken  string
}

func getConfig() (configType, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_ENDPOINT")
	if len(endpoint) == 0 {
		return configType{}, fmt.Errorf("missing env: OTEL_EXPORTER_ENDPOINT")
	}

	headers := os.Getenv("OTEL_EXPORTER_HEADERS")
	if len(headers) == 0 {
		return configType{}, fmt.Errorf("missing env: OTEL_EXPORTER_HEADERS")
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if len(serviceName) == 0 {
		return configType{}, fmt.Errorf("missing env: OTEL_SERVICE_NAME")
	}

	repository := os.Getenv("GITHUB_REPOSITORY")
	if len(repository) == 0 {
		return configType{}, fmt.Errorf("missing env: GITHUB_REPOSITORY")
	}

	runID := os.Getenv("GITHUB_RUN_ID")
	if len(runID) == 0 {
		return configType{}, fmt.Errorf("missing env: GITHUB_RUN_ID")
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	repoDetails := strings.Split(repository, "/")
	if len(repoDetails) != 2 {
		return configType{}, fmt.Errorf("invalid env: GITHUB_REPOSITORY")
	}

	headersMap := map[string]string{}

	headersSplit := strings.Split(headers, ",")
	for i := 0; i < len(headersSplit); i++ {
		splitHeaders := strings.Split(headersSplit[i], ":")
		if len(splitHeaders) != 2 {
			return configType{}, fmt.Errorf("invalid env: OTEL_EXPORTER_HEADERS")
		}
		headersMap[splitHeaders[0]] = splitHeaders[1]

	}

	return configType{
		otelEndpoint: endpoint,
		otelHeaders:  headersMap,
		serviceName:  serviceName,
		runID:        runID,
		repository:   repository,
		owner:        repoDetails[0],
		repo:         repoDetails[1],
		githubToken:  githubToken,
	}, nil
}
