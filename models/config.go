// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/26
package models

type Config struct {
	On      bool   `json:"-"`
	DocPath string `json:"-"` // 文件路径
	DocName string `json:"-"` // 文件名称,带后缀

	DocTitle    string            `json:"title"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Contact     map[string]string `json:"contact"`
	BaseUrls    map[string]string `json:"baseUrls"`
}
