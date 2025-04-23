package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// TODO 这一部分后续可以考虑怎么跟Gee结合一下

// 这里用如果是前缀的话个人觉得用prefix更好理解一点
const defaultBasePath = "/_geecache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.

type HTTPPool struct {
	// this peer's base URL, e.g "https://exaple.net:8000"
	self     string
	basePath string
}

// NewHTTPPoll initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServerHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断访问的路径前缀是否是basePath， 不是则抛出错误panic
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// 访问路径格式为 /<basepath>/<groupname>/<key>
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 通过groupName拿到group实例
	groupName := parts[0]
	key := parts[1]

	// 获取缓存数据
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	// 将缓存之作为httpResponse 的 body返回
	w.Write(view.ByteSlice())
}
