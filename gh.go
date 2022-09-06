package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/go-github/v47/github"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
)

func createTraces(ctx context.Context, conf configType) error {
	var token *http.Client
	if len(conf.githubToken) != 0 {
		token = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: conf.githubToken},
		))
	}

	client := github.NewClient(token)

	runID, err := strconv.ParseInt(conf.runID, 10, 64)
	if err != nil {
		return err
	}

	workflowData, _, err := client.Actions.GetWorkflowRunByID(ctx, conf.owner, conf.repo, runID)
	if err != nil {
		return err
	}

	ctx, workflowSpan := tracer.Start(ctx, *workflowData.Name, trace.WithTimestamp(workflowData.CreatedAt.Time))
	defer workflowSpan.End(trace.WithTimestamp(workflowData.UpdatedAt.Time))

	jobs, _, err := client.Actions.ListWorkflowJobs(ctx, conf.owner, conf.repo, runID, &github.ListWorkflowJobsOptions{})
	if err != nil {
		return err
	}

	for _, job := range jobs.Jobs {
		ctx, jobSpan := tracer.Start(ctx, *job.Name, trace.WithTimestamp(job.GetStartedAt().Time))

		for _, step := range job.Steps {
			_, stepSpan := tracer.Start(ctx, *step.Name, trace.WithTimestamp(step.GetStartedAt().Time))

			if step.CompletedAt != nil {
				stepSpan.End(trace.WithTimestamp(step.CompletedAt.Time))
			} else {
				stepSpan.End()
			}
		}

		if job.CompletedAt != nil {
			jobSpan.End(trace.WithTimestamp(job.CompletedAt.Time))
		} else {
			jobSpan.End()
		}
	}

	return nil
}
