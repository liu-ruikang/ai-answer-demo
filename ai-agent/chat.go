package aiagent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
)

// ChatModelType 定义支持的 ChatModel 类型
type ChatModelType string

const (
	OpenAIModel ChatModelType = "openai"
	OllamaModel ChatModelType = "ollama"
)

// NewChatModel 根据类型创建对应的 ChatModel
func NewChatModel(ctx context.Context, modelType ChatModelType) (schema.ChatModel, error) {
	switch modelType {
	case OpenAIModel:
		return openai.NewChatModel(ctx, &openai.ChatModelConfig{
			Model:  "gpt-4",
			APIKey: os.Getenv("OPENAI_API_KEY"),
		})
	case OllamaModel:
		return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: "http://localhost:11434",
			Model:   "llama2",
		})
	default:
		return nil, fmt.Errorf("unsupported chat model type: %s", modelType)
	}
}

// RunChatModelGenerate 使用 Generate 模式运行 ChatModel
func RunChatModelGenerate(chatModel schema.ChatModel, messages []*schema.Message) (string, error) {
	result, err := chatModel.Generate(context.Background(), messages)
	if err != nil {
		return "", err
	}
	return result.Content, nil
}

// RunChatModelStream 使用 Stream 模式运行 ChatModel
func RunChatModelStream(chatModel schema.ChatModel, messages []*schema.Message, streamCallback func(string)) error {
	streamResult, err := chatModel.Stream(context.Background(), messages)
	if err != nil {
		return err
	}

	reader := streamResult.Stream()
	defer reader.Close()

	for {
		message, err := reader.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		streamCallback(message.Content)
	}

	return nil
}
