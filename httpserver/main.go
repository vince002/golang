package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vince002/golang/httpserver/metrics"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	v := os.Getenv("v")
	flag.Set("v", v)
	flag.Parse()
	glog.V(2).Info("Starting http server v=2...")

	glog.V(4).Info("Starting http server v=4...")

	//module10: 1、注册metrics指标
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)

	//module10: 1、指定Handler
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	// 优雅终止，1分钟后执行kill
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	<-done
	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
	w.WriteHeader(http.StatusOK)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {


	fmt.Println("entering root handler")

	glog.V(2).Info("entering root handler v=2...")

	glog.V(4).Info("entering root handler v=4...")


	version := os.Getenv("VERSION")
	w.Header().Set("Version", version)

	//module10：2、输出指标
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	user := r.URL.Query().Get("user")

	//module10:1、为 HTTPServer 添加 0-2 秒的随机延时
	delay := randInt(10,2000)
	time.Sleep(time.Millisecond*time.Duration(delay))

	req, err := http.NewRequest("GET", "http://service1", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get failed with error: ", "error", err)
	} else {
		glog.Info("HTTP get succeeded, http://service1")
	}
	if resp != nil {
		resp.Write(w)
	}


	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}


	for k, v := range r.Header {
		for _, value := range v {
			w.Header().Add(string(k), string(value))
		}
	}

	ipAddress, err := GetIP(r)
	fmt.Println("Host=\n", ipAddress)

	if err != nil {
		fmt.Println("err=\n", err)
	}

	status := http.StatusOK
	w.WriteHeader(status)

	fmt.Println("http.Status =", status)
	fmt.Println("delay =", delay)
	glog.V(4).Infof("Respond in %d ms, v=4...", delay)

}

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

