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
	workflowID   string
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

	workflowID := os.Getenv("GITHUB_WORKFLOW")
	if len(workflowID) == 0 {
		return configType{}, fmt.Errorf("missing env: GITHUB_WORKFLOW")
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if len(githubToken) == 0 {
		return configType{}, fmt.Errorf("missing env: GITHUB_TOKEN")
	}

	repoDetails := strings.Split(repository, "/")
	if len(repoDetails) != 2 {
		return configType{}, fmt.Errorf("invalid env: GITHUB_REPOSITORY")
	}

	headersSplit := strings.Split(headers, ":")
	if len(headersSplit) != 2 {
		return configType{}, fmt.Errorf("invalid env: OTEL_EXPORTER_HEADERS")
	}

	headersMap := map[string]string{
		headersSplit[0]: headersSplit[1],
	}

	return configType{
		otelEndpoint: endpoint,
		otelHeaders:  headersMap,
		serviceName:  serviceName,
		workflowID:   workflowID,
		runID:        runID,
		repository:   repository,
		owner:        repoDetails[0],
		repo:         repoDetails[1],
		githubToken:  githubToken,
	}, nil
}
