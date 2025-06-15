# 中国象棋命令行游戏

这是一个用Go语言实现的命令行版中国象棋游戏，支持基本的棋子移动规则和游戏流程。

## 功能特性

- 9x10标准棋盘，使用Unicode字符渲染
- 支持所有中国象棋棋子（将/帅、士、象、马、车、炮、兵/卒）
- 实现完整的棋子移动规则验证
- 支持悔棋功能
- 命令行交互界面
- 新增Web界面支持

## 安装要求

- Go 1.20+
- 支持Unicode和ANSI颜色的终端
- 现代浏览器（用于Web界面）

## 编译运行

```bash
# 进入项目目录
cd china-xiangqi

# 编译命令行版本
go build -o chinese-chess main.go

# 运行命令行版本
./chinese-chess

# 编译Web版本
go build -o web-chess web/main.go

# 运行Web版本
./web-chess
```

## 操作说明

### 命令行模式
1. 游戏启动后会显示棋盘和当前回合方
2. 输入移动指令进行游戏，格式为`fromXfromYtoXtoY`
   - 示例: `e3e5` 表示将e3位置的棋子移动到e5
3. 特殊指令:
   - `undo`: 悔棋一步
   - `exit`: 退出游戏

### Web模式
1. 启动服务器后访问 http://localhost:8080
2. 点击棋子选择要移动的棋子
3. 再次点击目标位置完成移动
4. 页面会自动刷新并显示最新棋盘状态

## 目录结构

```
china-xiangqi/
├── main.go          # 命令行入口
├── web/
│   ├── main.go      # Web入口
│   ├── router.go    # Web路由处理
│   └── template.html# HTML模板
├── board/           # 棋盘相关代码
│   ├── board.go     # 棋盘数据结构
│   └── renderer.go  # 棱镜渲染器
├── pieces/          # 棋子相关代码
│   ├── piece.go     # 棋子基类及规则
│   └── *.go         # 各种具体棋子类型（预留扩展）
└── game/            # 游戏引擎
    ├── engine.go    # 游戏核心逻辑
    └── move_validator.go # 移动验证逻辑（待扩展）
```

## 后续开发计划

- 实现将军检测和将死判定
- 添加基础AI对手
- 扩展更多游戏功能（存档/读档等）
- 改进用户交互体验