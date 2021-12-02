package shadow

import (
	"errors"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"v.marlon.life/toolkit/util"
)

var (
	ErrNodeUnavailable = errors.New("get avalible shadowsocks server err")
	ErrServerInit      = errors.New("server init error")
)

var (
	proxyClient *http.Client
	latestUsed  = time.Now()
	refreshLock sync.Mutex

	dlerSubscribeList = "{your clash subscribe link}"
)

func init() {
	initProxyClient()
}

// GetProxyClient returns a proxy client
func GetProxyClient() *http.Client {
	refreshLock.Lock()
	defer refreshLock.Unlock()

	err := util.Retry(3, time.Second, func() error {
		// 增加了url检测, 所以检测时间, 所以
		if time.Since(latestUsed) > time.Hour*12 {
			killPlugin()
			initProxyClient()
		}

		return util.Retry(3, time.Second, func() error {
			if checkClientAvalible(proxyClient) {
				return nil
			}
			logger.Println(ErrServerInit.Error())
			return ErrServerInit
		})
	})

	if err != nil {
		panic(err)
	}

	latestUsed = time.Now()
	return proxyClient
}

// init defualt proxy client
func initProxyClient() {
	list := getDlerList(dlerSubscribeList)
	var node *ShadowConfig
	// telnet address
	for _, v := range list {
		_, err := net.DialTimeout("tcp", v.Addr, time.Second)
		if err != nil {
			continue
		}
		node = v
		break
	}

	if node == nil {
		panic(ErrNodeUnavailable)
	}

	dialer, _ := GetProxyDialer(node)
	proxyClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			Dial:                  dialer.Dial,
		},
	}
}

func checkClientAvalible(c *http.Client) bool {
	resp, err := c.Get("https://google.com")
	if err != nil {
		return false
	}
	if resp.StatusCode < 100 || resp.StatusCode >= 400 {
		return false
	}
	return true
}
