// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/26
package models

type ApiCall struct {
	Id int

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestUrlParams     map[string]string

	RequestBody  string
	ApiResponse  *ApiResponse
	ResponseBody string
	ResponseCode int
}
