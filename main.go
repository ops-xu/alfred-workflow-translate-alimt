package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"github.com/deanishe/awgo"
	"log"
	"unicode"
)

var wf *aw.Workflow
var subTitle = "回车自动复制"

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

func main() {

	// Wrap your entry point with Run() to catch and log panics and

	// show an error in Alfred instead of silently dying

	wf.Run(run)
}

// Your workflow starts here

func run() {
	args := wf.Args()
	if len(args) <= 0 {
		return
	}
	trans(args[0])
}

func trans(text string) {
	config := wf.Config

	accessKey := config.GetString("access_key")
	secret := config.GetString("secret")

	log.Println(accessKey)
	log.Println(secret)

	if accessKey == "" || secret == "" {
		wf.NewWarningItem("配置尚未完成", "请在环境变量中设置自己的access_key以及secret")
		wf.SendFeedback()
		return
	}

	// 创建ecsClient实例
	alimtClient, err := alimt.NewClientWithAccessKey(
		"cn-hangzhou", // 地域ID
		accessKey,     // 您的Access Key ID
		secret)        // 您的Access Key Secret
	if err != nil {
		// 异常处理
		panic(err)
	}
	// 创建API请求并设置参数
	request := alimt.CreateTranslateECommerceRequest()
	// 等价于 request.PageSize = "10"
	request.Method = "POST"         //设置请求
	request.FormatType = "text"     //翻译文本的格式
	request.SourceLanguage = "auto" //源语言
	request.SourceText = text       //原文
	if isHan(text) {
		request.TargetLanguage = "en"
	} else {
		request.TargetLanguage = "zh" //目标语言
	}
	request.Scene = "general" //目标语言
	// 发起请求并处理异常
	response, err := alimtClient.TranslateECommerce(request)
	log.Println(response)
	if err != nil {
		// 异常处理
		panic(err)
	}

	translated := response.Data.Translated
	// Add a "Script Filter" result

	wf.NewItem(translated).Subtitle(subTitle).Arg(translated).Valid(true)

	// Send results to Alfred

	wf.SendFeedback()
}

func isHan(text string) bool {
	for _, item := range text {
		if unicode.Is(unicode.Han, item) {
			return true
		}
	}
	return false
}
