// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/26
package models

type Spec struct {
	ApiInfo  *Config
	ApiSpecs []ApiSpec
}

// Len is the number of elements in the collection.
func (s *Spec) Len() int {
	return len(s.ApiSpecs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s *Spec) Less(i, j int) bool {
	if s.ApiSpecs[i].Path < s.ApiSpecs[j].Path {
		return true
	}
	return false
}

// Swap swaps the elements with indexes i and j.
func (s *Spec) Swap(i, j int) {
	s.ApiSpecs[i], s.ApiSpecs[j] = s.ApiSpecs[j], s.ApiSpecs[i]
}
