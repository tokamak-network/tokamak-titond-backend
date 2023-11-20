package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
	"k8s.io/apimachinery/pkg/runtime"
)

type TitondAPI struct {
	k8s         *kubernetes.Kubernetes
	db          db.Client
	fileManager services.IFIleManager
}

func NewTitondAPI(k8s *kubernetes.Kubernetes, db db.Client, fileManager services.IFIleManager) *TitondAPI {
	return &TitondAPI{
		k8s,
		db,
		fileManager,
	}
}

func (t *TitondAPI) CreateNetwork(data *model.Network) *model.Network {
	// t.fileManager.UploadContent("File_name_9", " New Content 9 ")
	result, _ := t.db.CreateNetwork(data)
	status, _ := t.k8s.GetPodStatus("default", "l2geth-0")
	fmt.Println(status)

	return result
}

func (t *TitondAPI) CreateL2Geth(data *model.Component) {
	// TODO : deal with DB
	t.db.CreateComponent()

	namespace := "default"
	createL2Geth(t.k8s, namespace)
}

func createL2Geth(client *kubernetes.Kubernetes, namespace string) {
	obj := getObject("l2geth", "configMap")
	client.CreateConfigMap(namespace, obj)

	obj = getObject("l2geth", "pvc")
	client.CreatePersistentVolumeClaim(namespace, obj)

	obj = getObject("l2geth", "service")
	client.CreateService(namespace, obj)

	obj = getObject("l2geth", "statefulset")
	client.CreateStatefulSet(namespace, obj)

	obj = getObject("l2geth", "ingress")
	client.CreateIngress(namespace, obj)
}

func getObject(dir, file string) runtime.Object {
	f := kubernetes.GetYAMLfile(dir, file)
	return kubernetes.ConvertBytestoObject(f)
}
