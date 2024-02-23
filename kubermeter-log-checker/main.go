package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
	"k8s.io/api/core/v1" // Corrected import for PodLogOptions

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
    "k8s.io/client-go/tools/remotecommand"
    "k8s.io/client-go/kubernetes/scheme"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error getting in-cluster config: %s\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %s\n", err)
		os.Exit(1)
	}

	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		fmt.Printf("Error reading namespace: %s\n", err)
		os.Exit(1)
	}
	ns := strings.TrimSpace(string(namespace))
    print("The current namespace is " + ns + "\n")

  //pattern := regexp.MustCompile(`Configuring remote engine`)
	pattern := regexp.MustCompile(`Finished the test on host jmeter-[0-9]{1,2}`)

	for {
        print("Checking if pattern is found in pod logs...\n")
		found := checkPods(clientset, ns, pattern)
		if found {
			fmt.Println("Pattern found, exiting...")
            executeCommandInJMeterMaster(clientset, ns, config)
			break
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}

func checkPods(clientset *kubernetes.Clientset, namespace string, pattern *regexp.Regexp) bool {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app in (jmeter, jmeter-master)",
	})
	if err != nil {
		fmt.Printf("Error listing pods: %s\n", err)
		return false
	}

	for _, pod := range pods.Items {
		logOptions := &v1.PodLogOptions{}
        // get logs from pod
        
		podLogs, err := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, logOptions).Stream(context.TODO())
		if err != nil {
			fmt.Printf("Error getting logs for pod %s: %s\n", pod.Name, err)
			continue
		}
        // print pod name and logs
        print("Pod name: " + pod.Name + "\n")
        // ./podLogs.String ined (type io.ReadCloser has no field or method String
        //print(podLogs)

		defer podLogs.Close()

		byteLogs, err := ioutil.ReadAll(podLogs)
		if err != nil {
			fmt.Printf("Error reading logs for pod %s: %s\n", pod.Name, err)
			continue
		}

		logContent := string(byteLogs)
		if pattern.MatchString(logContent) {
			fmt.Printf("Pattern found in pod %s\n", pod.Name)
			return true
		}

        print(logContent)
	}
	return false
}

func executeCommandInJMeterMaster(clientset *kubernetes.Clientset, namespace string, config *rest.Config) {
	labelSelector := "app=jmeter-master" // Label selector to identify the pod
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
    
    podName := pods.Items[0].Name
    command := []string{"/bin/sh", "/opt/apache-jmeter-5.5/bin/stoptest.sh"}

	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
        SubResource("exec").
        Param("container", "jmeter").
        VersionedParams(&v1.PodExecOptions{
            Command: command,
            Stdin:   true,
            Stdout:  true,
            Stderr:  true,
            TTY:     true,
        }, scheme.ParameterCodec) // Corrected typo ParameterCodec

    


	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Printf("Error creating executor: %s\n", err)
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		fmt.Printf("Error executing command in jmeter-master pod: %s\n", err)
	}
}

