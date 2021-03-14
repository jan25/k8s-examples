package main

import (
	"os"
	"os/signal"
	"syscall"
	"sync"
	"fmt"
	"log"
	"time"
	"context"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	var wg sync.WaitGroup
	stopCh := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Interruption caught.")
		close(stopCh)
	}()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	wg.Add(1)
	go func(){
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				fmt.Println("Cleaning pod watcher.")
				return
			case <-ticker.C:
				watchPods(clientset)	
			}
		}
	}()

	wg.Wait()
	fmt.Println("Exited.")
}

func watchPods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("I see %d pods.\n", len(pods.Items));
}