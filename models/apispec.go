package models

type ApiSpec struct {
	HttpVerb   []string
	MethodsMap map[string]string
	Path       string
	Calls      []ApiCall
}
