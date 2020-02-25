package dialogflow

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

type DialogFlowAgent struct {
	ProjectID        string
	AuthJSONFilePath string
	Lang             string
	TimeZone         string
	AgentsClient     *dialogflow.AgentsClient
	ctx              context.Context
}

type DialogFlowContent struct {
	ProjectID        string
	AuthJSONFilePath string
	Lang             string
	TimeZone         string
	ContextsClient   *dialogflow.ContextsClient
	ctx              context.Context
}

type DialogFlowIntents struct {
	ProjectID        string
	AuthJSONFilePath string
	Lang             string
	TimeZone         string
	IntentsClient    *dialogflow.IntentsClient
	ctx              context.Context
}

type NLPResponse struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

var dp DialogFlowProcessor

func (dp *DialogFlowProcessor) Init(ProjectID, AuthJSONFilePath, Lang, TimeZone string) error {
	dp.ProjectID = ProjectID
	dp.AuthJSONFilePath = AuthJSONFilePath
	dp.Lang = Lang
	dp.TimeZone = TimeZone

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.AuthJSONFilePath))

	if err != nil {
		log.Fatal("Error in auth with DialogFlow")
		return err
	}

	dp.SessionClient = sessionClient
	return nil
}

func (dp *DialogFlowProcessor) ProcessNLP(rewMessage, username string) (r NLPResponse) {
	sessionID := username

	request := dialogflowpb.DetectIntentRequest{
		Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.ProjectID, sessionID),
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

func (da *DialogFlowAgent) InitAgen(ProjectID, AuthJSONFilePath, Lang, TimeZone string) error {
	da.ProjectID = ProjectID
	da.AuthJSONFilePath = AuthJSONFilePath
	da.Lang = Lang
	da.TimeZone = TimeZone
	da.ctx = context.Background()
	var err error
	da.AgentsClient, err = dialogflow.NewAgentsClient(da.ctx, option.WithCredentialsFile(da.AuthJSONFilePath))
	if err != nil {
		return err
	}

	return nil
}

func (da *DialogFlowAgent) GetAgent(p string) (*dialogflowpb.Agent, error) {
	req := &dialogflowpb.GetAgentRequest{
		Parent: fmt.Sprintf("projects/%s", p),
	}
	resp, err := da.AgentsClient.GetAgent(da.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (da *DialogFlowAgent) SetAgent(req *dialogflowpb.SetAgentRequest) (*dialogflowpb.Agent, error) {
	// req := &dialogflowpb.SetAgentRequest{}
	resp, err := da.AgentsClient.SetAgent(da.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (dc *DialogFlowContent) InitContent(ProjectID, AuthJSONFilePath, Lang, TimeZone string) error {
	dc.ProjectID = ProjectID
	dc.AuthJSONFilePath = AuthJSONFilePath
	dc.Lang = Lang
	dc.TimeZone = TimeZone
	dc.ctx = context.Background()
	var err error
	dc.ContextsClient, err = dialogflow.NewContextsClient(dc.ctx, option.WithCredentialsFile(dc.AuthJSONFilePath))
	if err != nil {
		return err
	}
	return nil
}

func (dc *DialogFlowContent) GetListContexts(sessionID string) {
	req := &dialogflowpb.ListContextsRequest{Parent: fmt.Sprintf("projects/%s/agent/sessions/%s", dc.ProjectID, sessionID)}
	fmt.Println("err", req)
	resp := dc.ContextsClient.ListContexts(dc.ctx, req)
	// dc.ContextsClient.UpdateContext
	fmt.Println("=", resp)
}

func (dc *DialogFlowContent) CreateContent() {
	// m := structpb.Value("m")
	// m.Descriptor
	// a:=m.GetNullValue()
	// req := &dialogflowpb.CreateContextRequest{
	// 	Context: &dialogflowpb.Context{
	// 		Name:          "test",
	// 		LifespanCount: 1,
	// 		// Parameters:    &_struct.Struct{Fields: map[string]{}},
	// 	},
	// }
}

func (dc *DialogFlowContent) UpdateContext(sessionID string, c *dialogflowpb.Context) {
	req := &dialogflowpb.UpdateContextRequest{Context: c}
	fmt.Println("err", req)
	// resp, err := dc.ContextsClient.UpdateContext(dc.ctx, req)
	// fmt.Println("=", resp)
}

func (dc *DialogFlowContent) DeleteContext(sessionID, name string) {
	req := &dialogflowpb.DeleteContextRequest{Name: name}
	err := dc.ContextsClient.DeleteContext(dc.ctx, req)
	fmt.Println(err)
}

func (di *DialogFlowIntents) Init() error {
	var err error

	di.IntentsClient, err = dialogflow.NewIntentsClient(di.ctx, option.WithCredentialsFile(di.AuthJSONFilePath))
	if err != nil {
		return err
	}
	return nil
}

func (di *DialogFlowIntents) GetIntents(name, languageCode string, intentView dialogflowpb.IntentView) (*dialogflowpb.Intent, error) {
	req := &dialogflowpb.GetIntentRequest{Name: name, LanguageCode: languageCode, IntentView: intentView}
	resp, err := di.IntentsClient.GetIntent(di.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// func (di *DialogFlowIntents) ListIntent(projectID, languageCode, pageToken string, pageSize int32, intentView dialogflowpb.IntentView) *dialogflowpb.IntentIterator {
// 	req := &dialogflowpb.ListIntentsRequest{Parent: fmt.Sprintf("projects/%s/agent", projectID), LanguageCode: languageCode, PageSize: pageSize, PageToken: pageToken, IntentView: intentView}
// 	resp := di.IntentsClient.ListIntents(di.ctx, req)
// 	return resp
// }
