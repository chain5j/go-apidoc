// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/26
package models

type ApiInfo struct {
	Title       string            `json:"title"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Contact     map[string]string `json:"contact"`
}
