package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fission/fission-workflows/pkg/apiserver/httpclient"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type client struct {
	Admin      *httpclient.AdminAPI
	Workflow   *httpclient.WorkflowAPI
	Invocation *httpclient.InvocationAPI
}

func getClient(ctx Context) client {

	url := ctx.GlobalString("url")

	// fetch the FISSION_URL env variable. If not set, port-forward to controller.
	if len(url) == 0 {
		fissionURL := os.Getenv("FISSION_URL")
		if len(fissionURL) == 0 {
			fissionNamespace := getFissionNamespace()
			kubeConfig := getKubeConfigPath()
			localPort := setupPortForward(kubeConfig, fissionNamespace, "application=fission-api")
			url = "http://127.0.0.1:" + localPort
			logrus.Debugf("Forwarded Fission API to %s.", url)
		} else {
			url = fissionURL
		}
	}
	path := ctx.GlobalString("path-prefix")
	if path[0] != '/' {
		path = "/" + path
	}
	url = url + strings.TrimSuffix(path, "/")
	httpClient := http.Client{}
	return client{
		Admin:      httpclient.NewAdminAPI(url, httpClient),
		Workflow:   httpclient.NewWorkflowAPI(url, httpClient),
		Invocation: httpclient.NewInvocationAPI(url, httpClient),
	}
}

func getFissionNamespace() string {
	fissionNamespace := os.Getenv("FISSION_NAMESPACE")
	return fissionNamespace
}

func getKubeConfigPath() string {
	kubeConfig := os.Getenv("KUBECONFIG")
	if len(kubeConfig) == 0 {
		home := os.Getenv("HOME")
		kubeConfig = filepath.Join(home, ".kube", "config")

		if _, err := os.Stat(kubeConfig); os.IsNotExist(err) {
			panic("Couldn't find kubeconfig file. Set the KUBECONFIG environment variable to your kubeconfig's path.")
		}
	}
	return kubeConfig
}

func findFreePort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		return "", nil
	}

	err = listener.Close()
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return port, nil
}

// runPortForward creates a local port forward to the specified pod
func runPortForward(kubeConfig string, labelSelector string, localPort string, fissionNamespace string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Kubernetes: %s", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Kubernetes: %s", err))
	}

	// if fission namespace is unset, try to find a fission pod in any namespace
	if len(fissionNamespace) == 0 {
		fissionNamespace = meta_v1.NamespaceAll
	}

	// get the pod; if there is more than one, ask the user to disambiguate
	podList, err := clientset.CoreV1().Pods(fissionNamespace).List(meta_v1.ListOptions{LabelSelector: labelSelector})
	if err != nil || len(podList.Items) == 0 {
		panic("Error getting controller pod for port-forwarding")
	}

	// make a useful error message if there is more than one install
	if len(podList.Items) > 1 {
		namespaces := make([]string, 0)
		for _, p := range podList.Items {
			namespaces = append(namespaces, p.Namespace)
		}
		panic(fmt.Sprintf("Found %v fission installs, set FISSION_NAMESPACE to one of: %v",
			len(podList.Items), strings.Join(namespaces, " ")))
	}

	// pick the first pod
	podName := podList.Items[0].Name
	podNameSpace := podList.Items[0].Namespace

	// get the service and the target port
	svcs, err := clientset.CoreV1().Services(podNameSpace).
		List(meta_v1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		panic(fmt.Sprintf("Error getting %v service :%v", labelSelector, err.Error()))
	}
	if len(svcs.Items) == 0 {
		panic(fmt.Sprintf("Service %v not found", labelSelector))
	}
	service := &svcs.Items[0]

	var targetPort string
	for _, servicePort := range service.Spec.Ports {
		targetPort = servicePort.TargetPort.String()
	}

	stopChannel := make(chan struct{}, 1)
	readyChannel := make(chan struct{})

	// create request URL
	req := clientset.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(podNameSpace).Name(podName).SubResource("portforward")
	url := req.URL()

	// create ports slice
	portCombo := localPort + ":" + targetPort
	ports := []string{portCombo}

	// actually start the port-forwarding process here
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		msg := fmt.Sprintf("newexecutor errored out :%v", err.Error())
		panic(msg)
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)

	fw, err := portforward.New(dialer, ports, stopChannel, readyChannel, nil, os.Stderr)
	if err != nil {
		msg := fmt.Sprintf("portforward.new errored out :%v", err.Error())
		panic(msg)
	}

	return fw.ForwardPorts()
}

// Port forward a free local port to a pod on the cluster. The pod is
// found in the specified namespace by labelSelector. The pod's port
// is found by looking for a service in the same namespace and using
// its targetPort. Once the port forward is started, wait for it to
// start accepting connections before returning.
func setupPortForward(kubeConfig, namespace, labelSelector string) string {
	localPort, err := findFreePort()
	if err != nil {
		panic(fmt.Sprintf("Error finding unused port :%v", err.Error()))
	}

	for {
		conn, _ := net.DialTimeout("tcp",
			net.JoinHostPort("", localPort), time.Millisecond)
		if conn != nil {
			conn.Close()
		} else {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	go func() {
		err := runPortForward(kubeConfig, labelSelector, localPort, namespace)
		if err != nil {
			panic(fmt.Sprintf("Error forwarding to controller port: %s", err.Error()))
		}
	}()

	for {
		conn, _ := net.DialTimeout("tcp",
			net.JoinHostPort("", localPort), time.Millisecond)
		if conn != nil {
			conn.Close()
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	return localPort
}
