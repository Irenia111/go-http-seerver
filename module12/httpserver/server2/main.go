package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"module12/httpserver/metrics"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	filename       string
	lastModifyTime int64
	Http           Http `yaml:"http"`
	Log            Log  `yaml:"log"`
}

type Http struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Log struct {
	Level string `yaml:"level"`
}

var (
	log        *logrus.Logger
	config     *Config
	configLock = new(sync.RWMutex)
)

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func webRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("entering web root handler")
	// 为 HTTPServer 添加 0-2 秒的随机延时
	delay := randInt(0, 2000)
	time.Sleep(time.Microsecond * time.Duration(delay))
	io.WriteString(w, "===================arrive server1, invoke server3============\n")

	// 创建新请求
	newReq, err := http.NewRequest("GET", "http://service3/", nil)

	if err != nil {
		fmt.Printf("%s", err)
	}
	// 创建新请求头
	lowerCaseHeader := make(http.Header)
	// 将当前请求中的请求头复制到新请求头中
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	log.Info("headers:", lowerCaseHeader)
	// 为新请求设置请求头
	newReq.Header = lowerCaseHeader
	// 创建 http 客户端
	client := &http.Client{}
	// 调用 service3
	newResp, err := client.Do(newReq)
	if err != nil {
		log.Info("HTTP get failed with error: ", "error", err)
	} else {
		log.Info("HTTP get succeeded")
	}
	// 将 service2 的响应通过当前请求输出
	if newResp != nil {
		io.WriteString(w, "===================server2 print response from server3============\n")
		newResp.Write(w)
	}
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	fmt.Println("结束退出...")
	os.Exit(0)
}

func startHttpServer(listenServer string) {
	os.Setenv("VERSION", "1")
	// 注册一个 prometheus 指标采集器
	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/", webRoot)
	// 增加 prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(listenServer, mux); err != nil {
		log.Fatalf("start http server failed, error: %s\n", err.Error())
	}
}

func LogInit(logLevel string) {
	log = logrus.New()
	log.Out = os.Stdout
	level := logrus.InfoLevel
	switch {
	case logLevel == "debug":
		level = logrus.DebugLevel
	case logLevel == "info":
		level = logrus.InfoLevel
	case logLevel == "warn":
		level = logrus.WarnLevel
	case logLevel == "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("Get current path fail\n", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func (config *Config) reload() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		func() {
			f, err := os.Open(GetConfig().filename)
			if err != nil {
				log.Fatalf("open file error:%s\n", err)
				return
			}
			defer f.Close()

			fileInfo, err := f.Stat()
			if err != nil {
				log.Fatalf("stat file error:%s\n", err)
				return
			}
			curModifyTime := fileInfo.ModTime().Unix()
			if curModifyTime > GetConfig().lastModifyTime {
				log.Info("cfg change, load new cfg ...")
				loadConfig()
				GetConfig().lastModifyTime = curModifyTime
			}
		}()
	}
}

// 记录响应信息的结构体
type RespRecoder struct {
	http.ResponseWriter
	StatusCode int
}

// 重写 http.ResponseWriter WriteHeader 记录 statusCode
func (r *RespRecoder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func loadConfig() bool {
	log.Println("Load cfg ...")
	f, err := ioutil.ReadFile(config.filename)
	if err != nil {
		log.Fatalln("load config error: ", err)
		return false
	}
	temp := new(Config)
	err = yaml.Unmarshal(f, &temp)
	if err != nil {
		log.Fatalln("Para config failed: ", err)
		return false
	}
	temp.filename = GetConfig().filename
	temp.lastModifyTime = GetConfig().lastModifyTime
	log.Debugf("now cfg:%#v\n", temp)
	configLock.Lock()
	config = temp
	configLock.Unlock()
	return true
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func IsIpv4(ipv4 string) bool {
	address := net.ParseIP(ipv4)
	if address != nil {
		log.Infof("%s is a legal ipv4 address\n", ipv4)
		return true
	} else {
		log.Infof("%s is not a legal ipv4 address\n", ipv4)
		return false
	}
}

func CheckPortRange(port int) bool {
	if 1 <= port && port <= 65535 {
		return true
	}
	return false
}

func CheckConfig() (listenServer, logLevel string) {
	var allConfig = GetConfig()
	var httpConfig = allConfig.Http
	var httpPort = httpConfig.Port
	var httpHost = httpConfig.Host
	if !IsIpv4(httpHost) {
		httpHost = "0.0.0.0"
	}
	if port, err := strconv.Atoi(httpPort); err == nil {
		if !CheckPortRange(port) {
			httpPort = "8080"
		}
	}
	listenServer = httpHost + ":" + httpPort
	logLevel = allConfig.Log.Level
	return listenServer, logLevel
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func init() {
	LogInit("info")
	config = new(Config)
	pwd := getCurrentDirectory()
	config.filename = pwd + "/conf/config.yaml"
	if !loadConfig() {
		os.Exit(1)
	}
	go config.reload()
}

func main() {
	listenServer, logLevel := CheckConfig()
	LogInit(logLevel)
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGTERM:
				log.Info("捕获 SIGTERM 不退出", s)
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT:
				log.Info("退出", s)
				ExitFunc()
			default:
				fmt.Println("other", s)
			}
		}
	}()

	startHttpServer(listenServer)
}
