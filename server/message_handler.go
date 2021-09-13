package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	cx "cloud.google.com/go/dialogflow/cx/apiv3"
	"github.com/sainiajay/backend.ajaysaini.dev/services/bot"
	cxpb "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
)

const (
	location_id = "global"
	project_id  = "golden-tempest-325806"
	agent_id    = "74b701bc-a608-477d-a144-59df280ed813"
)

func detectIntent(text string) (*cxpb.DetectIntentResponse, error) {
	session_id := "abcd-1234"
	session_path := fmt.Sprintf("projects/%s/locations/%s/agents/%s/sessions/%s", project_id, location_id, agent_id, session_id)
	log.Printf("Session_path: %s", session_path)
	ctx := context.Background()
	c, err := cx.NewSessionsClient(ctx)
	if err != nil {
		log.Printf("Error occurred %v", err)
		return nil, err
	}
	defer c.Close()
	text_input := cxpb.QueryInput_Text{
		Text: &cxpb.TextInput{
			Text: text,
		},
	}
	query_input := cxpb.QueryInput{
		LanguageCode: "en",
		Input:        &text_input,
	}
	req := &cxpb.DetectIntentRequest{
		Session:    session_path,
		QueryInput: &query_input,
	}
	return c.DetectIntent(ctx, req)
}

type server struct {
	bot.UnimplementedBotServiceServer
}

func (s *server) HandleUserMessage(ctx context.Context, message *bot.Message) (*bot.Message, error) {
	response, err := detectIntent(message.GetBody())
	if err != nil {
		log.Printf("detectIntent Failed %v; %v", response, err)
	}
	var responseText []string
	for _, message := range response.QueryResult.ResponseMessages {
		responseText = append(responseText, message.GetText().GetText()...)
	}
	text := strings.Join(responseText, "\n")
	log.Printf("Received 1: %v", text)
	return &bot.Message{Body: text}, nil
}
