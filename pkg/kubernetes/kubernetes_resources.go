package kubernetes

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesApp struct {
	app *kubernetes.Clientset
}

func NewKubernetesClient() (*KubernetesApp, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubernetesApp{clientset}, nil
}

func int32Ptr(i int32) *int32 { return &i }

// deploymentを作成する処理
func (k *KubernetesApp) CreateDeployment(app string, deploymentName string, registryName string, envVars []EnvVar) error {
	//deploymentの定義
	deployment := DeploymentDefinition(app, deploymentName, registryName, envVars)

	//k8sに送信
	deploymentsClient := k.app.AppsV1().Deployments(apiv1.NamespaceDefault)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// serviceを作成する処理
func (k *KubernetesApp) CreateService(app string, serviceName string, targetPort int) error {
	//serviceの定義
	service := ServiceDefinition(app, serviceName, targetPort)

	//k8sに送信
	serviceClient := k.app.CoreV1().Services(apiv1.NamespaceDefault)
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}
	fmt.Printf("Created service%q.\n", result.GetObjectMeta().GetName())
	return nil
}

// ingressを作成する処理
func (k *KubernetesApp) CreateIngress(ingressName string, hostName string, serviceName string) error {
	//ingressの定義
	ingress := IngressDefinition(ingressName, hostName, serviceName)

	//k8sに送信
	ingressClient := k.app.NetworkingV1().Ingresses("default")
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ingress: %v", err)
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// kanikoのjobを生成する処理
func (k *KubernetesApp) CreateJob(githubUrl string, appName string, registryName string, envVars []EnvVar) (string, string, error) {
	//jobの定義
	job := JobDefinition(githubUrl, appName, registryName, envVars)

	//k8sに送信
	jobClient := k.app.BatchV1().Jobs("default")
	result, err := jobClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to create job: %v", err)
	}
	uid := string(result.UID)
	name := result.Name

	return name, uid, nil
}

// pvcを作成する処理
func (k *KubernetesApp) CreatePvc(jobName string, jobUid string, appName string) error {
	//pvcの定義
	pvc := PvcDefinition(jobName, jobUid, appName)

	// PVCを作成
	pvcClient := k.app.CoreV1().PersistentVolumeClaims("default")
	result, err := pvcClient.Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created PVC %q.\n", result.GetObjectMeta().GetName())
	return err
}

// jobを監視する処理
func (k *KubernetesApp) CheckJobCompletion(jobName string) error {

	for {
		job, err := k.app.BatchV1().Jobs("default").Get(context.Background(), jobName, metav1.GetOptions{})
		if err != nil {
			panic(fmt.Errorf("failed to get job status: %v", err))
		}
		if job.Status.Succeeded > 0 {
			fmt.Println("Job completed successfully!")
			break
		} else if job.Status.Failed > 0 {
			// ジョブのポッドのログを取得してエラーメッセージを出力する

			fmt.Printf("Error")
			return fmt.Errorf("job %s failed with %d failed pods", jobName, job.Status.Failed)
		}
		fmt.Println("Job is still running...")
		time.Sleep(15 * time.Second)
	}
	return nil
}

// pod内のlogを取得する処理
func (k *KubernetesApp) GetPodLog(podName string) (string, error) {
	// k8sの初期化処理
	namespace := "default" // Specify the namespace
	fmt.Println(podName)

	// podのLogを取得
	podLogOpts := apiv1.PodLogOptions{}
	req := k.app.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer podLogs.Close()
	// podのLogを読み出して、stringに帰る
	var sb strings.Builder
	buf := make([]byte, 2000)
	for {
		numBytes, err := podLogs.Read(buf)
		if numBytes == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		sb.Write(buf[:numBytes])
	}
	logOutput := sb.String()
	fmt.Println(logOutput)
	return logOutput, nil
}

// deploymentを削除する処理
func (k *KubernetesApp) DeleteDeployment(deploymentName string) error {

	deploymentClient := k.app.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting deployment...")
	if err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete deployment: %v", err)
	}
	fmt.Println("Deleted deployment.")
	return nil
}

// serviceを削除する処理
func (k *KubernetesApp) DeleteService(serviceName string) error {

	serviceClient := k.app.CoreV1().Services(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting service...")
	if err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}
	fmt.Println("Deleted service.")
	return nil
}

// ingressを削除する処理
func (k *KubernetesApp) DeleteIngress(ingressName string) error {

	ingressClient := k.app.NetworkingV1().Ingresses("default")
	deletePolicy := metav1.DeletePropagationForeground

	fmt.Println("Deleting ingress...")
	if err := ingressClient.Delete(context.TODO(), ingressName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("failed to delete ingress: %v", err)
	}
	fmt.Println("Deleted ingress.")
	return nil
}

// deployment名からpodのステータスを確認する処理
func (k *KubernetesApp) GetDeploymentStatus(deploymentName string) (string, error) {

	namespace := "default"

	//デプロイメントの取得
	deployment, err := k.app.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	replicaSets, err := k.app.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return "", err
	}

	if len(replicaSets.Items) == 0 {
		log.Fatalf("No ReplicaSets found for deployment %s", deploymentName)
	}

	var latestReplicaSet *appsv1.ReplicaSet
	for _, rs := range replicaSets.Items {
		if latestReplicaSet == nil || rs.CreationTimestamp.After(latestReplicaSet.CreationTimestamp.Time) {
			latestReplicaSet = &rs
		}
	}

	if latestReplicaSet == nil {
		return "", err
	}

	podList, err := k.app.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: latestReplicaSet.Spec.Template.Labels,
		}),
	})

	if err != nil {
		return "", err
	}

	if len(podList.Items) > 0 {
		// 最初のPodのステータスを表示
		pod := podList.Items[0]
		fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
		return string(pod.Status.Phase), nil
	} else {
		fmt.Println("No pods found for the deployment")
		return "No pods found for the deployment", nil
	}

}

func (k *KubernetesApp) GetJobsStatus(jobName string) (string, error) {

	namespace := "default"
	response := "AppCreating"

	job, err := k.app.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	labelSelector := metav1.FormatLabelSelector(job.Spec.Selector)
	pods, err := k.app.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return "", err
	}
	fmt.Print(pods)

	pod := pods.Items[0]
	if pod.Status.Phase == "Running" {
		return response, nil
	} else {
		response = "Unknown"
		return response, nil
	}
}
