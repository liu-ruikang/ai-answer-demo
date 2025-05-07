package aiagent

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/components/tool/utils"
)

// AddTodoParams 定义添加 Todo 的参数
type AddTodoParams struct {
	Content   string `json:"content" jsonschema:"description=content of the todo"`
	StartedAt *int64 `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64 `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
}

// AddTodoFunc 添加 Todo 的处理函数
func AddTodoFunc(_ context.Context, params *AddTodoParams) (string, error) {
	return `{"msg": "add todo success"}`, nil
}

// GetAddTodoTool 使用 NewTool 构建 add_todo 工具
func GetAddTodoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "Add a todo item",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}

	return utils.NewTool(info, AddTodoFunc)
}

// UpdateTodoParams 定义更新 Todo 的参数
type UpdateTodoParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

// UpdateTodoFunc 更新 Todo 的处理函数
func UpdateTodoFunc(_ context.Context, params *UpdateTodoParams) (string, error) {
	return `{"msg": "update todo success"}`, nil
}

// GetUpdateTodoTool 使用 InferTool 构建 update_todo 工具
func GetUpdateTodoTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"update_todo",
		"Update a todo item, eg: content,deadline...",
		UpdateTodoFunc,
	)
}

// ListTodoTool 实现 Tool 接口
type ListTodoTool struct{}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "List all todo items",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return `{
		"todos": [
			{
				"id": "1",
				"content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作",
				"started_at": 1717401600,
				"deadline": 1717488000,
				"done": false
			}
		]
	}`, nil
}
