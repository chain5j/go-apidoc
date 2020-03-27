// description: xblog 
// 
// @author: xwc1125
// @date: 2020/3/26
package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/xwc1125/go-apidoc"
	"github.com/xwc1125/go-apidoc/models"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

/* 32 MB in memory max */
const MaxInMemoryMultipartSize = 32000000

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Trailer":           true,
	"Accept-Encoding":   false,
	"Accept-Language":   false,
	"Cache-Control":     false,
	"Connection":        false,
	"Origin":            false,
	"User-Agent":        false,
}

type YaagHandler struct {
	nextHandler http.Handler
}

func Handle(nextHandler http.Handler) http.Handler {
	return &YaagHandler{nextHandler: nextHandler}
}

func (y *YaagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !apidoc.IsOn() {
		y.nextHandler.ServeHTTP(w, r)
		return
	}
	writer := NewResponseRecorder(w)
	apiCall := models.ApiCall{}
	Before(&apiCall, r)
	y.nextHandler.ServeHTTP(writer, r)
	After(&apiCall, writer, r)
}

func HandleFunc(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !apidoc.IsOn() {
			next(w, r)
			return
		}
		apiCall := models.ApiCall{}
		recorder := NewResponseRecorder(w)
		Before(&apiCall, r)
		next(recorder, r)
		After(&apiCall, recorder, r)
	})
}

func HandleFunc2(w http.ResponseWriter, r *http.Request,
	next func(), recorderFun func() *ResponseRecorder) {
	if !apidoc.IsOn() {
		next()
		return
	}
	apiCall := models.ApiCall{}
	recorder := NewResponseRecorder(w)
	Before(&apiCall, r)
	next()
	if recorder != nil {
		recorder = recorderFun()
	}
	After(&apiCall, recorder, r)
}

func Before(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader = ReadHeaders(req)
	apiCall.RequestUrlParams = ReadQueryParams(req)
	val, ok := apiCall.RequestHeader["Content-Type"]
	if ok {
		ct := strings.TrimSpace(val)
		switch ct {
		case "application/x-www-form-urlencoded":
			fallthrough
		case "application/json, application/x-www-form-urlencoded":
			apiCall.PostForm = ReadPostForm(req)
		case "application/json":
			apiCall.RequestBody = *ReadBody(req)
		default:
			if strings.Contains(ct, "multipart/form-data") {
				handleMultipart(apiCall, req)
			} else {
				apiCall.RequestBody = *ReadBody(req)
			}
		}
	}
}

func ReadQueryParams(req *http.Request) map[string]string {
	params := map[string]string{}
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return params
	}
	for k, v := range u.Query() {
		if len(v) < 1 {
			continue
		}
		// TODO: v is a list, and we should be showing a list of values
		// rather than assuming a single value always, gotta change this
		params[k] = v[0]
	}
	return params
}

func handleMultipart(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader["Content-Type"] = "multipart/form-data"
	req.ParseMultipartForm(MaxInMemoryMultipartSize)
	apiCall.PostForm = ReadMultiPostForm(req.MultipartForm)
}

func ReadMultiPostForm(mpForm *multipart.Form) map[string]string {
	postForm := map[string]string{}
	for key, val := range mpForm.Value {
		postForm[key] = val[0]
	}
	return postForm
}

func ReadPostForm(req *http.Request) map[string]string {
	postForm := map[string]string{}
	for _, param := range strings.Split(*ReadBody(req), "&") {
		value := strings.Split(param, "=")
		postForm[value[0]] = value[1]
	}
	return postForm
}

// 获取头部
func ReadHeaders(req *http.Request) map[string]string {
	b := bytes.NewBuffer([]byte(""))
	err := req.Header.WriteSubset(b, reqWriteExcludeHeaderDump)
	if err != nil {
		return map[string]string{}
	}
	headers := map[string]string{}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if strings.EqualFold(values[0], "") {
			continue
		}
		headers[values[0]] = values[1]
	}
	return headers
}

// 从响应中读取头部
func ReadHeadersFromResponse(header http.Header) map[string]string {
	headers := map[string]string{}
	for k, v := range header {
		headers[k] = strings.Join(v, " ")
	}
	return headers
}

// 获取请求内容
func ReadBody(req *http.Request) *string {
	save := req.Body
	var err error
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return nil
		}
	}
	b := bytes.NewBuffer([]byte(""))
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if req.Body == nil {
		return nil
	}
	var dest io.Writer = b
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
	}
	req.Body = save
	body := b.String()
	return &body
}

func After(apiCall *models.ApiCall, record *ResponseRecorder, r *http.Request) {
	if strings.Contains(r.RequestURI, ".ico") || !apidoc.IsOn() || r.RequestURI == "/" {
		return
	}
	apiCall.MethodType = r.Method
	apiCall.CurrentPath = r.URL.Path
	bodyBytes := record.Body.Bytes()
	var apiResp models.ApiResponse
	err := json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		log.Println("Unmarshal apiResp", "err", err)
	}

	apiCall.ResponseBody = string(bodyBytes)
	apiCall.ApiResponse = &apiResp
	apiCall.ResponseCode = record.Status
	apiCall.ResponseHeader = ReadHeadersFromResponse(record.Header())
	if apidoc.IsStatusCodeValid(record.Status) {
		go apidoc.GenerateJson(apiCall)
	}
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
