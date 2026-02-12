package middleware

import (
	"net/http"
)

// MethodHandler 根据 HTTP 方法路由到不同的处理函数
func MethodHandler(methodMap map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求方法
		method := r.Method

		// 查找对应的处理函数
		handler, ok := methodMap[method]
		if !ok {
			// 如果请求方法不支持，返回 405 Method Not Allowed
			w.Header().Set("Allow", getAllowedMethods(methodMap))
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// 调用对应的处理函数
		handler(w, r)
	}
}

// getAllowedMethods 获取所有允许的 HTTP 方法
func getAllowedMethods(methodMap map[string]http.HandlerFunc) string {
	allowedMethods := ""
	for method := range methodMap {
		if allowedMethods != "" {
			allowedMethods += ", "
		}
		allowedMethods += method
	}
	return allowedMethods
}
