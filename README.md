# ApiDoc

This Project is base on yaag! https://github.com/betacraft/yaag

Thanks betacraft very much!

## ApiDoc for net.http or Gorilla Mux

- 1、Init ApiDoc

```
	apidoc.Init(&models.Config{
		On:       true,          // 是否开启
		DocTitle: "ApiDoc",      // 文档名称
		Version:  "v1.0.0",      // 版本号
		DocPath:  "./",          // 路径
		DocName:  "apidoc.json", // 文件名称，带json
		BaseUrls: baseUrls,      // baseUrls
	})
```

- 2、Add middleware handler

```
http.HandleFunc("/", middleware.HandleFunc(handler))
```

## ApiDoc for iris

- 1、Init ApiDoc

```
	apidoc.Init(&models.Config{
		On:       true,          // 是否开启
		DocTitle: "ApiDoc",      // 文档名称
		Version:  "v1.0.0",      // 版本号
		DocPath:  "./",          // 路径
		DocName:  "apidoc.json", // 文件名称，带json
		BaseUrls: baseUrls,      // baseUrls
	})
```

- 2、Create iris context.Handler

```
func New() context.Handler {
	return func(ctx context.Context) {
		middleware.HandleFunc2(ctx.ResponseWriter(), ctx.Request(), func() {
			ctx.Record()
			ctx.Next()
		}, func() *middleware.ResponseRecorder {
			r := middleware.NewResponseRecorder(ctx.Recorder().Naive())
			r.Body = bytes.NewBuffer(ctx.Recorder().Body())
			r.Status = ctx.Recorder().StatusCode()
			return r
		})
	}
}
```

- 3、Add middleware handler

```
var app *iris.Application
app.Use(New()) //irisyaag记录响应主体并向apidoc提供所有必要的信息
```

