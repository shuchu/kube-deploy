package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func setupKubeAPI() *kubernetes.Clientset {

	var kubeconfig string
	var homeDir string

	if homeDir = os.Getenv("HOME"); homeDir == "" {
		fmt.Println("=> Oh no! Couldn't figure out what your homedir is, please set the environment variable $HOME.")
		os.Exit(1)
	}

	kubeconfig = filepath.Join(homeDir, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func kubeAPIGetSingleDeployment(name string) *v1beta1.Deployment {
	deployment, _ := repoConfig.KubeAPIClientSet.
		ExtensionsV1beta1().Deployments(repoConfig.Namespace).
		Get(name, metav1.GetOptions{})
	// Return even if nil
	return deployment
}

func kubeAPIUpdateDeployment(deployment *v1beta1.Deployment) *v1beta1.Deployment {
	newDeployment, err := repoConfig.KubeAPIClientSet.
		ExtensionsV1beta1().Deployments(repoConfig.Namespace).
		Update(deployment)
	if err != nil {
		panic(err.Error())
	}
	if runFlags.Bool("debug") {
		fmt.Printf("=> => Updated deployment %s.\n", deployment.Name)
	}
	return newDeployment

}

func kubeAPIAddDeploymentLabel(deployment *v1beta1.Deployment, key string, value string) {
	existingLabels := deployment.GetLabels()
	existingLabels[key] = value
}

func kubeAPIRemoveDeploymentLabel(deployment *v1beta1.Deployment, key string) {
	existingLabels := deployment.GetLabels()
	delete(existingLabels, key)
}

func kubeAPIDeleteDeployment(deployment *v1beta1.Deployment) {
	deletePolicy := metav1.DeletePropagationForeground

	if err := repoConfig.KubeAPIClientSet.
		ExtensionsV1beta1().Deployments(repoConfig.Namespace).
		Delete(deployment.Name, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
		panic(err.Error())
	}
}

func kubeAPIDeleteService(service *v1.Service) {
	if err := repoConfig.KubeAPIClientSet.CoreV1().Services(repoConfig.Namespace).
		Delete(service.Name, nil); err != nil {
		panic(err.Error())
	}
}

func kubeAPIDeleteSecret(secret *v1.Secret) {
	if err := repoConfig.KubeAPIClientSet.CoreV1().Secrets(repoConfig.Namespace).
		Delete(secret.Name, nil); err != nil {
		panic(err.Error())
	}
}

func kubeAPIDeleteIngress(ingress *v1beta1.Ingress) {
	if err := repoConfig.KubeAPIClientSet.ExtensionsV1beta1().Ingresses(repoConfig.Namespace).
		Delete(ingress.Name, nil); err != nil {
		panic(err.Error())
	}
}

func kubeAPIListDeployments(labelFilter map[string]string) *v1beta1.DeploymentList {

	label := labels.Set(labelFilter)

	deployments, err := repoConfig.KubeAPIClientSet.
		ExtensionsV1beta1().Deployments(repoConfig.Namespace).
		List(metav1.ListOptions{LabelSelector: label.String()})

	if err != nil {
		panic(err.Error())
	}

	return deployments
	// // Examples for error handling:
	// // - Use helper functions like e.g. errors.IsNotFound()
	// // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	// _, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
	// if errors.IsNotFound(err) {
	// 	fmt.Printf("Pod not found\n")
	// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	// } else if err != nil {
	// 	panic(err.Error())
	// } else {
	// 	fmt.Printf("Found pod\n")
	// }

	// time.Sleep(10 * time.Second)

}
