// Package main

// import "fmt"

// func main() {
// 	// TODO: Add test paper generation logic here
// }

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// 定义试题结构体
type Question struct {
	Number    int
	Content   string
	Answer    string
	IsApplied bool // 是否为应用题
}

// 定义知识点结构体
type KnowledgePoint struct {
	Title   string
	Content string
}

// 定义题目解析结构体
type QuestionAnalysis struct {
	Question Question
	Analysis string
}

// 启动预览服务器
func startPreviewServer() {
	// 提供静态文件服务
	http.Handle("/shijuan/", http.StripPrefix("/shijuan/", http.FileServer(http.Dir("shijuan"))))
	
	// 启动HTTP服务器
	fmt.Println("正在本地启动试卷预览服务，访问 http://localhost:8080/shijuan/test_paper.html")
	http.ListenAndServe(":8080", nil)
}

func main() {
	// 创建试卷文件
	createTestPaper("shijuan/test_paper.html", "shijuan/answers.txt")
	
	// 创建课件文件
	os.MkdirAll("shijuan", os.ModePerm) // 强制创建目标目录，避免路径错误
	createLessonPlan("shijuan/lesson_plan.html")
	
	// 启动预览服务器
	startPreviewServer()
}

// 创建试卷和答案文件
func createTestPaper(htmlPath, answerPath string) {
	// 生成20道题目（8填空，8选择，4应用）
	questions := []Question{
		// 填空题（每题5分）
		{Number: 1, Content: "3 + 5 = ___", Answer: "8"},
		{Number: 2, Content: "10 - 7 = ___", Answer: "3"},
		{Number: 3, Content: "4 × 2 = ___", Answer: "8"},
		{Number: 4, Content: "12 ÷ 3 = ___", Answer: "4"},
		{Number: 5, Content: "正方形有___条边", Answer: "4"},
		{Number: 6, Content: "1米=___厘米", Answer: "100"},
		{Number: 7, Content: "3点15分时针指向___", Answer: "3"},
		{Number: 8, Content: "长方形的对边___", Answer: "相等"},
		
		// 选择题（每题5分）
		{Number: 9, Content: "下列哪个是质数？A.4 B.5 C.6 D.8", Answer: "B"},
		{Number: 10, Content: "三位数最大值是？A.100 B.999 C.99 D.1000", Answer: "B"},
		{Number: 11, Content: "1小时=？分钟 A.30 B.60 C.90 D.120", Answer: "B"},
		{Number: 12, Content: "平行四边形有几条对称轴？A.0 B.1 C.2 D.4", Answer: "A"},
		{Number: 13, Content: "π的近似值是？A.3 B.3.1 C.3.14 D.3.141", Answer: "C"},
		{Number: 14, Content: "等边三角形每个角是？A.60° B.90° C.120° D.180°", Answer: "A"},
		{Number: 15, Content: "最小的两位数是？A.10 B.11 C.99 D.100", Answer: "A"},
		{Number: 16, Content: "闰年有多少天？A.364 B.365 C.366 D.367", Answer: "C"},
		
		// 应用题（每题5分）
		{Number: 17, Content: "小明有10个苹果，吃掉了3个，请问还剩几个？", Answer: "7个", IsApplied: true},
		{Number: 18, Content: "正方形的边长为5cm，求周长？", Answer: "20cm", IsApplied: true},
		{Number: 19, Content: "汽车每小时行驶60公里，3小时行驶多少公里？", Answer: "180公里", IsApplied: true},
		{Number: 20, Content: "一本书有100页，每天读5页，几天读完？", Answer: "20天", IsApplied: true},
	}
	// 创建HTML模板
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>数学练习题</title>
	<style>
		body { font-family: 微软雅黑; margin: 2cm auto; width: 21cm; background-color: #fff; color: #000; }
		.header { text-align: center; margin-bottom: 2cm; }
		.questions { line-height: 2em; }
		.applied { margin-bottom: 3cm; }
	</style>
</head>
<body>
	<div class="header">
		<h1>{{if eq .Grade 1}}一年级{{else if eq .Grade 2}}二年级{{else}}六年级{{end}}数学专项测试题</h1>
		<p>日期：__________ 姓名：__________</p>
		<div class="score">满分：100分</div>
	</div>
	<ol class="questions">
	{{range .Questions}}
		<li>
			{{.Content}}
			{{if .IsApplied}}
			<div class="applied"></div>
			{{end}}
		</li>
	{{end}}
	</ol>
</body>
</html>`

	// 解析模板并写入文件
	tmpl, _ := template.New("testPaper").Parse(htmlTemplate)
	htmlFile, _ := os.Create(htmlPath)
	defer htmlFile.Close()
	tmpl.Execute(htmlFile, map[string]interface{}{
		"Grade":      2, // 年级参数
		"Questions": questions,
	})

	// 写入答案文件
	answerFile, _ := os.Create(answerPath)
	defer answerFile.Close()
	for _, q := range questions {
		answerFile.WriteString(q.Content + " 答案：" + q.Answer + "\n")
	}
}

// 创建课件文件
func createLessonPlan(filePath string) {
	// 获取绝对路径
	parsedPath, _ := filepath.Abs(filePath)
	
	// 确保目录存在
	dir := filepath.Dir(parsedPath)
	os.MkdirAll(dir, os.ModePerm)
	
	// 定义知识点内容
	knowledgePoints := []struct {
		Title   string
		Content string
	}{
		{"四则运算", "加减乘除的基本运算规则和优先级"},
		{"几何基础", "常见图形的性质和计算公式"},
		{"时间与单位", "时间单位换算和基本计量单位"},
	}

	// 定义重点题型解析
	analysis := []struct {
		QuestionNum int
		Question    string
		Answer      string
		Analysis    string
	}{
		{1, "3 + 5 = ___", "8", "加法运算：将3和5相加得到和8，注意保持位数对齐"},
		{17, "小明有10个苹果，吃掉了3个，请问还剩几个？", "7个", "减法应用：总数10个减去吃掉的3个，剩余7个。注意理解'吃掉'表示减少的概念"},
	}

	// 构建完整的HTML内容
	htmlContent := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>知识讲解</title>
	<style>
		body { font-family: 微软雅黑; margin: 2cm auto; width: 21cm; background-color: #fff; color: #000; }
		.header { text-align: center; margin-bottom: 2cm; }
		.section { margin-bottom: 2cm; }
		.title { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 5px; }
	</style>
</head>
<body>
	<div class="header">
		<h1>数学知识讲解与重点分析</h1>
	</div>

	<div class="section">
		<h2 class="title">核心知识点</h2>`

	// 添加知识点部分
	for _, kp := range knowledgePoints {
		htmlContent += fmt.Sprintf(`
		<div>
			<h3>%s</h3>
			<p>%s</p>
		</div>`, kp.Title, kp.Content)
	}

	// 添加题型解析部分
	htmlContent += `
	</div>

	<div class="section">
		<h2 class="title">重点题型解析</h2>`

	for _, a := range analysis {
		htmlContent += fmt.Sprintf(`
		<div>
			<h3>第%d题：%s</h3>
			<p>答案：%s</p>
			<p>解析：%s</p>
		</div>`, a.QuestionNum, a.Question, a.Answer, a.Analysis)
	}

	// 添加HTML结尾
	htmlContent += `
	</div>
</body>
</html>`

	// 写入文件
	err := os.WriteFile(parsedPath, []byte(htmlContent), os.ModePerm)
	if err != nil {
		fmt.Printf("文件写入错误: %v\n", err)
	}
}
