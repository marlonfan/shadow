package shadow

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ShadowConfig shadowsocks 配置
type ShadowConfig struct {
	Addr        string
	Cipher      string
	Password    string
	Plugin      string
	PluginParam string
}

// getDlerList 获取 dler 列表
func getDlerList(addr string) []*ShadowConfig {
	list := make([]*ShadowConfig, 0)
	resp, err := http.Get(addr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	raw, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		panic(err)
	}
	for _, v := range strings.Split(string(raw), "\n") {
		if !strings.HasPrefix(v, "ss://") {
			continue
		}
		u, err := url.Parse(v)
		if err != nil {
			continue
		}
		if u.User == nil {
			continue
		}
		r, err := base64.RawURLEncoding.DecodeString(u.User.Username())
		if err != nil {
			panic(err)
		}
		var cipher, password, plugin, pluginParam string
		key := strings.Split(string(r), ":")
		cipher, password = key[0], key[1]
		plugins := strings.Split(u.Query().Get("plugin"), ";")
		if len(plugins) > 2 {
			plugin = plugins[0]
			pluginParam = plugins[1] + ";" + plugins[2]
		}
		list = append(list, &ShadowConfig{Addr: u.Host, Cipher: cipher, Password: password, Plugin: plugin, PluginParam: pluginParam})
	}
	return list
}
