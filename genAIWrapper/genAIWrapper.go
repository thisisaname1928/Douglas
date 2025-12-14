package genaiwrapper

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"google.golang.org/genai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Secret struct {
	GEMINI_API_KEY string `json:"GEMINI_API_KEY"`
	MAX_RETRY      int    `json:"MAX_RETRY"`
}
var client *genai.Client
var IsInit = false

func InitGenAIWrapper() error {

	b, e := os.ReadFile("secretKey.json")

	if e != nil {
		return e
	}

	e = json.Unmarshal(b, &Secret)

	if e != nil {
		return e
	}

	ctx := context.Background()

	client, e = genai.NewClient(ctx, &genai.ClientConfig{APIKey: Secret.GEMINI_API_KEY,
		Backend: genai.BackendGeminiAPI})

	if e != nil {
		return e
	}

	IsInit = true
	return nil
}

func GeminiGenContent(input string) (string, error) {
	ctx := context.Background()
	var e error = nil
	var res *genai.GenerateContentResponse = nil

	TEMP := float32(0.85)
	TopK := float32(40)
	TopP := float32(0.95)

	for i := 0; i < Secret.MAX_RETRY; i++ {
		res, e = client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(input), &genai.GenerateContentConfig{Temperature: &TEMP, TopP: &TopP, TopK: &TopK})

		if e == nil {
			break
		}

		statusCode, ok := status.FromError(e)

		if !ok {
			return "", e
		}

		switch statusCode.Code() {
		case codes.ResourceExhausted:
			time.Sleep(time.Millisecond * 1000)
			continue
		case codes.Unavailable:
			time.Sleep(time.Millisecond * 500)
			continue
		case codes.InvalidArgument:
			return "", e
		case codes.PermissionDenied:
			return "", e
		}

	}

	return res.Text(), nil
}

func AutoGenContent(input string) (string, error) {
	if !IsInit {
		return "", errors.New("ERR_GEN_AI_NOT_AVAILABLE")
	}

	return GeminiGenContent(input)
}
