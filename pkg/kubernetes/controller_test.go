package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/suite"
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

}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(FakeClientSuite))
}
