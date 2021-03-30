package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	pods, _ := clientset.CoreV1().Pods("footstone-common").List(metav1.ListOptions{})
	for _, p := range pods.Items {
		fmt.Println(p.GetName())
	}

	// 初始化pod所在的corev1资源组，发送请求
	// PodExecOptions struct 包括Container stdout stdout  Command 等结构
	// scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的



	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name("api-gateway-ng-2-5cb8fb9fd-6ppq8").
		Namespace("footstone-common").
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: []string{"bash"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	// remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请求
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	// 检查是不是终端
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		fmt.Errorf("stdin/stdout should be terminal")
	}

	// 这个应该是处理Ctrl + C 这种特殊键位
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Println(err)
	}
	defer terminal.Restore(0, oldState)

	// 用IO读写替换 os stdout
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	// 建立链接之后从请求的sream中发送、读取数据
	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    true,
	}); err != nil {
		fmt.Print(err)
	}
}
