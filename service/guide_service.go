package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type guideResponse struct {
	Code int       `json:"code"`
	Data guideData `json:"data"`
}

type guideData struct {
	App struct {
		Name     string `json:"name"`
		Tagline  string `json:"tagline"`
		Audience string `json:"audience"`
	} `json:"app"`
	Hero struct {
		Title    string   `json:"title"`
		Subtitle string   `json:"subtitle"`
		Chips    []string `json:"chips"`
	} `json:"hero"`
	Phases    []guidePhase    `json:"phases"`
	Checklist []guideChecklist `json:"checklist"`
	Tips      []guideTip      `json:"tips"`
	FAQ       []guideFAQ      `json:"faq"`
	Footer    guideFooter     `json:"footer"`
}

type guidePhase struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Duration    string   `json:"duration"`
	Description string   `json:"description"`
	Tasks       []string `json:"tasks"`
}

type guideChecklist struct {
	Title string   `json:"title"`
	Items []string `json:"items"`
}

type guideTip struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type guideFAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type guideFooter struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

// HomeHandler renders the static frontend page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("./index.html")
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, string(data))
}

// GuideHandler returns the mini-program development guide content.
func GuideHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res := guideResponse{
		Code: 0,
		Data: buildGuideData(),
	}

	body, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(body))
}

func buildGuideData() guideData {
	var data guideData

	data.App.Name = "小程序开发流程卡片"
	data.App.Tagline = "像看微信小程序一样，轻松看懂从 0 到上线"
	data.App.Audience = "送给正在了解小程序开发的她"

	data.Hero.Title = "把小程序开发拆成 4 个清晰阶段"
	data.Hero.Subtitle = "先想清楚做什么，再搭页面、接接口、测试上线。每一步都不复杂，只要顺着流程做。"
	data.Hero.Chips = []string{
		"需求梳理",
		"页面搭建",
		"接口联调",
		"测试上线",
	}

	data.Phases = []guidePhase{
		{
			ID:          "01",
			Title:       "确认需求",
			Duration:    "第 1 步",
			Description: "先决定这个小程序解决什么问题，用户是谁，最核心的功能只留 2 到 3 个。",
			Tasks: []string{
				"列出首页、列表页、详情页这些基础页面",
				"画一个很粗的页面草图，先不用精修视觉",
				"确定登录、表单、支付这些是否真的需要",
			},
		},
		{
			ID:          "02",
			Title:       "搭前端页面",
			Duration:    "第 2 步",
			Description: "用小程序常见的卡片、分组标题和底部按钮来组织内容，先把浏览流程跑通。",
			Tasks: []string{
				"写页面结构和样式，优先保证手机上清晰好读",
				"给重要按钮留出明显的点击区域",
				"先用假数据把界面填满，别一开始就卡在接口",
			},
		},
		{
			ID:          "03",
			Title:       "接后端接口",
			Duration:    "第 3 步",
			Description: "当页面结构稳定后，再把静态内容替换成真实接口，处理加载、失败和刷新状态。",
			Tasks: []string{
				"定义接口返回的数据结构",
				"前端通过请求拿到列表、详情、提示文案",
				"把空状态和错误提示一起补上",
			},
		},
		{
			ID:          "04",
			Title:       "测试和上线",
			Duration:    "第 4 步",
			Description: "最后检查功能、文案和交互，再提交体验版或正式版。",
			Tasks: []string{
				"真机上检查字体、间距、滚动和按钮反馈",
				"把容易出错的路径都点一遍",
				"整理上线前的图标、名称、介绍和截图",
			},
		},
	}

	data.Checklist = []guideChecklist{
		{
			Title: "开始前准备",
			Items: []string{
				"想清楚目标用户和核心功能",
				"准备页面草图或参考图",
				"列出至少 1 个主流程",
			},
		},
		{
			Title: "联调时关注",
			Items: []string{
				"接口字段名是否统一",
				"加载中、失败、空数据是否可见",
				"按钮点击后是否给到明确反馈",
			},
		},
		{
			Title: "上线前检查",
			Items: []string{
				"文案有没有错别字",
				"图片和封面是否清晰",
				"重要页面是否都能返回上一层",
			},
		},
	}

	data.Tips = []guideTip{
		{
			Title:   "先做能跑通的版本",
			Content: "不要一开始就追求功能很多。能让用户从首页顺利走到结果页，已经是一个很好的第一版。",
		},
		{
			Title:   "把复杂问题拆小",
			Content: "页面像积木，接口像数据管道。一个页面一个页面做，一个接口一个接口接，进度会很稳。",
		},
		{
			Title:   "先对，再漂亮",
			Content: "排版和动画可以后补，但信息层级、按钮位置和主流程一定要先做对。",
		},
	}

	data.FAQ = []guideFAQ{
		{
			Question: "小程序开发一定要先写后端吗？",
			Answer:   "不一定。更高效的做法通常是先把前端页面和交互原型搭出来，再按真实页面需要去定义接口。",
		},
		{
			Question: "页面很多，会不会很难开始？",
			Answer:   "先只做一个最核心的主流程，比如首页到详情页。能跑通一个闭环后，剩下页面基本就是同样的方法扩展。",
		},
		{
			Question: "怎么判断项目已经可以上线？",
			Answer:   "当主流程稳定、关键页面真机可用、文案和状态完整，而且没有明显阻塞 bug，就可以先发体验版验证。",
		},
	}

	data.Footer = guideFooter{
		Signature: "给最可爱的产品体验官",
		Message:   "这页不是在讲很难的技术，而是在告诉你：小程序开发其实就是把想法一步一步落地。我把流程都整理好了，剩下的我们慢慢做。",
	}

	return data
}
