package aiagent

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// EncouragerPromptTemplate 创建程序员鼓励师的对话模板
func EncouragerPromptTemplate() *prompt.PromptTemplate {
	return prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("问题: {question}"),
	)
}

// RunEncourager 启动程序员鼓励师并处理用户问题
func RunEncourager(ctx context.Context) {
	// 创建 ChatModel（使用 OpenAI）
	chatModel, err := NewChatModel(ctx, "ollama")
	if err != nil {
		log.Fatal(err)
	}

	// 创建对话模板
	template := EncouragerPromptTemplate()

	// 使用模板生成消息
	messages, err := template.Format(ctx, map[string]interface{}{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 运行 ChatModel 并获取结果
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println("鼓励师回复:")
	fmt.Println(result.Content)
}
