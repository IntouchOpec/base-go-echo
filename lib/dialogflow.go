package lib

import (
	"context"
	"fmt"
	"log"
	"strconv"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/api/option"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type DialogFlowProcessor struct {
	ProjectID        string
	AuthJSONFilePath string
	Lang             string
	TimeZone         string
	SessionClient    *dialogflow.SessionsClient
	ctx              context.Context
}

type NLPResponse struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

var dp DialogFlowProcessor

func (dp *DialogFlowProcessor) init(a ...string) error {
	dp.ProjectID = a[0]
	dp.AuthJSONFilePath = a[1]
	dp.Lang = a[2]
	dp.TimeZone = a[3]

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.AuthJSONFilePath))

	if err != nil {
		log.Fatal("Error in auth with DialogFlow")
		return err
	}

	dp.SessionClient = sessionClient
	return nil
}

func (dp *DialogFlowProcessor) processNLP(rewMessage, username string) (r NLPResponse) {
	sessionID := username

	request := dialogflowpb.DetectIntentRequest{
		Session: fmt.Sprintf("projects/agent/sessions/%s", dp.ProjectID, sessionID),
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         rewMessage,
					LanguageCode: dp.Lang,
				},
			},
		},
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone: dp.TimeZone,
		},
	}

	response, err := dp.SessionClient.DetectIntent(dp.ctx, &request)
	if err != nil {
		log.Fatal("Error in communication with Dialogflow %s", err.Error())
		return
	}
	queryResult := response.GetQueryResult()
	if queryResult.Intent != nil {
		r.Intent = queryResult.Intent.DisplayName
		r.Confidence = float32(queryResult.IntentDetectionConfidence)
	}
	r.Entities = make(map[string]string)
	params := queryResult.Parameters.GetFields()
	if len(params) > 0 {
		for index, param := range params {
			fmt.Printf("Param %s: %s (%s)", index, param.GetStringValue(), param.String())
			extractedValue := extractDialogflowEntities(param)
			r.Entities[index] = extractedValue
		}
	}
	return
	// request := dialogflowp
	// dialogflowdp.DeleteIn

}

func extractDialogflowEntities(p *structpb.Value) (extractedEntity string) {
	kind := p.GetKind()
	switch kind.(type) {
	case *structpb.Value_StringValue:
		return p.GetStringValue()
	case *structpb.Value_NumberValue:
		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
	case *structpb.Value_BoolValue:
		return strconv.FormatBool(p.GetBoolValue())
	case *structpb.Value_StructValue:
		s := p.GetStructValue()
		fields := s.GetFields()
		extractedEntity := ""
		for key, value := range fields {
			if key == "amount" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(value.GetNumberValue(), 'f', 6, 64))
			}
			if key == "unit" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
			if key == "date_time" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
		}
		return extractedEntity
	case *structpb.Value_ListValue:
		list := p.GetListValue()
		if len(list.GetValues()) > 1 {

		}
		extractedEntity = extractDialogflowEntities(list.GetValues()[0])
		return extractedEntity
	default:
		return ""
	}
}
