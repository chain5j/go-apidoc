// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/25
package apidoc

import (
	"encoding/json"
	"github.com/xwc1125/go-apidoc/models"
	"github.com/xwc1125/go-apidoc/utils"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var config *models.Config

// Initial empty spec
var spec *models.Spec = &models.Spec{}

func IsOn() bool {
	return config.On
}

func Init(conf *models.Config) {
	if conf.DocPath == "" {
		conf.DocPath = "./"
	} else if !strings.HasSuffix(conf.DocPath, "/") {
		conf.DocPath = conf.DocPath + "/"
	}
	config = conf
	// load the config file
	if conf.DocName == "" {
		conf.DocName = "apidoc.json"
	}

	filePath, err := filepath.Abs(conf.DocPath + conf.DocName)
	exists, err := utils.PathExists(filePath)
	if !exists {
		os.MkdirAll(conf.DocPath, os.ModePerm)
	}
	dataFile, err := os.Open(filePath)
	defer dataFile.Close()
	if err == nil {
		json.NewDecoder(io.Reader(dataFile)).Decode(spec)

		sort.Sort(spec)
		generateHtml()
	}
}

func GenerateJson(apiCall *models.ApiCall) {
	shouldAddPathSpec := true
	for k, apiSpec := range spec.ApiSpecs {
		// 该地址已经存在
		if apiSpec.Path == apiCall.CurrentPath {
			shouldAddPathSpec = false
			deleteCommonHeaders(apiCall)
			avoid := false
			//avoid2 := true
			// 判断请求方式是否已经存在
			if _, ok := apiSpec.MethodsMap[apiCall.MethodType]; ok {
				// Method 已经存在

			} else {
				// Method 不存在
				apiSpec.MethodsMap[apiCall.MethodType] = apiCall.MethodType
				apiSpec.HttpVerb = append(apiSpec.HttpVerb, apiCall.MethodType)
			}
			isExist := false
			for _, currentApiCall := range apiSpec.Calls {
				if apiCall.ResponseCode == currentApiCall.ResponseCode {
					// 请求数据，和响应的code一致
					if apiCall.ApiResponse != nil && currentApiCall.ApiResponse != nil && apiCall.ApiResponse.Code == currentApiCall.ApiResponse.Code {
						// 最终的返回结果一致，此时不需要重新添加
						isExist = true
					}
				}
			}

			if isExist {
				avoid = true
			}
			if !avoid {
				// 添加新的数据
				spec.ApiSpecs[k].Calls = append(apiSpec.Calls, *apiCall)
			}
		}
	}

	// 新建
	if shouldAddPathSpec {
		apiSpec := models.ApiSpec{
			MethodsMap: make(map[string]string),
			Path:       apiCall.CurrentPath,
		}
		apiSpec.MethodsMap[apiCall.MethodType] = apiCall.MethodType
		apiSpec.HttpVerb = append(apiSpec.HttpVerb, apiCall.MethodType)
		deleteCommonHeaders(apiCall)
		apiSpec.Calls = append(apiSpec.Calls, *apiCall)
		spec.ApiSpecs = append(spec.ApiSpecs, apiSpec)
		spec.ApiInfo = config
	}

	filePath, err := filepath.Abs(config.DocPath + config.DocName)
	dataFile, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer dataFile.Close()

	sort.Sort(spec)
	data, err := json.Marshal(spec)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dataFile.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
	generateHtml()
}

func generateHtml() {
	funcs := template.FuncMap{}
	t := template.New(config.DocTitle).Funcs(funcs)
	htmlString := Template
	t, err := t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	filePath, err := filepath.Abs(config.DocPath + config.DocName + ".html")
	if err != nil {
		panic("Error while creating file path : " + err.Error())
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()
	if err != nil {
		panic("Error while creating documentation file : " + err.Error())
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{
		"Title":    config.DocTitle,
		"baseUrls": config.BaseUrls,
		"array":    spec.ApiSpecs,
	})
}

// 删除通用的Headers
func deleteCommonHeaders(call *models.ApiCall) {
	delete(call.RequestHeader, "Accept")
	delete(call.RequestHeader, "Accept-Encoding")
	delete(call.RequestHeader, "Accept-Language")
	delete(call.RequestHeader, "Cache-Control")
	delete(call.RequestHeader, "Connection")
	delete(call.RequestHeader, "Cookie")
	delete(call.RequestHeader, "Origin")
	delete(call.RequestHeader, "User-Agent")
}

func IsStatusCodeValid(code int) bool {
	if code >= 200 && code < 300 {
		return true
	} else {
		return false
	}
}
