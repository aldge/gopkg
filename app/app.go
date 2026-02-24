package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// app 类型定义
type AppName string

// app 定义列表（简写）
const (
	AppNameDev        AppName = "dev" // 测试开发专用
	AppNameLecture    AppName = "lc"  // 讲座
	AppNameShortDrama AppName = "sd"  // 短剧
	AppNameShortVideo AppName = "sv"  // 小视频
	AppNameMooc       AppName = "mc"  // 慕课
	AppNameMovie      AppName = "mv"  // 影视
)

// 可用的 app 列表
var validApps = map[AppName]bool{
	AppNameDev:        true,
	AppNameLecture:    true,
	AppNameShortDrama: true,
	AppNameShortVideo: true,
	AppNameMooc:       true,
	AppNameMovie:      true,
}

const QueryParamName = "app" // 请求参数名

// AppIsLecture 判断应用是否是讲座
func AppIsLecture(app AppName) bool {
	return app == AppNameLecture
}

// AppIsShortDrama 判断应用是否是短剧
func AppIsShortDrama(app AppName) bool {
	return app == AppNameShortDrama
}

// AppIsShortVideo 判断应用是否是小视频
func AppIsShortVideo(app AppName) bool {
	return app == AppNameShortVideo
}

// AppIsMooc 判断应用是否是慕课
func AppIsMooc(app AppName) bool {
	return app == AppNameMooc
}

// AppIsMovie 判断应用是否是影视
func AppIsMovie(app AppName) bool {
	return app == AppNameMovie
}

// AppIsValid 判断应用是否合法
func AppIsValid(app AppName) bool {
	return validApps[app]
}

// GetAppName 获取应用名称
func GetAppName(ctx *gin.Context) AppName {
	return AppName(ctx.Query(QueryParamName))
}

// CtxAppIsLecture ctx 判断应用是否是讲座
func CtxAppIsLecture(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	return AppIsLecture(app)
}

// CtxAppIsShortDrama ctx 判断应用是否是短剧
func CtxAppIsShortDrama(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	return AppIsShortDrama(app)
}

// CtxAppIsShortVideo ctx 判断应用是否是小视频
func CtxAppIsShortVideo(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	return AppIsShortVideo(app)
}

// CtxAppIsMooc ctx 判断应用是否是慕课
func CtxAppIsMooc(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	return AppIsMooc(app)
}

// CtxAppIsMovie ctx 判断应用是否是影视
func CtxAppIsMovie(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	return AppIsMovie(app)
}

// CtxAppIsValid ctx 判断应用是否合法
func CtxAppIsValid(ctx *gin.Context) bool {
	app := GetAppName(ctx)
	// todo 根据正式域名屏蔽掉 app dev
	return AppIsValid(app)
}

// SetAppQueryParam 设置 HTTP 请求的 app 查询参数
func SetAppQueryParam(req *http.Request, app AppName) {
	if req == nil {
		return
	}
	values := req.URL.Query()
	values.Set(QueryParamName, string(app))
	req.URL.RawQuery = values.Encode()
}
