// description: Go-ApiDoc 
// 
// @author: xwc1125
// @date: 2020/3/27
package utils

import "os"

// 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
