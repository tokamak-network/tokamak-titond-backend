package kubernetes

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type FakeClientSuite struct {
	suite.Suite
	fakeKubernetes *Kubernetes
}

func (fcs *FakeClientSuite) SetupSuite() {
	fakeClient := fake.NewSimpleClientset()
	fcs.fakeKubernetes = &Kubernetes{fakeClient}
}

func (fcs *FakeClientSuite) TestApplyResource() {
	fcs.fakeKubernetes.CreateL2Geth("testNamespace", "testName")
	res, _ := fcs.fakeKubernetes.client.AppsV1().StatefulSets("testNamespace").Get(context.TODO(), "testName", metav1.GetOptions{})
	fmt.Println(res.Name + "!!!!!!!!!!!!!!")
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(FakeClientSuite))
}
