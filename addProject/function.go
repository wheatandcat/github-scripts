package addproject

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v29/github"
)

func GitHubEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	webhookEvent, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event := webhookEvent.(type) {
	case *github.IssuesEvent:
		if err := processIssuesEvent(ctx, event); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
