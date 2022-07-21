package crdadaptor

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/util"
	k8sauthzv1 "github.com/casbin/k8s-gatekeeper/pkg/apis/k8sauthz/v1"
	"github.com/casbin/k8s-gatekeeper/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sCRDAdaptor struct {
	namespace string
	modelName string
	clientset *versioned.Clientset
}

func NewK8sAdaptor(namespace string, modelName string, isExternalClient bool) (*K8sCRDAdaptor, error) {
	var res = &K8sCRDAdaptor{
		namespace: namespace,
		modelName: modelName,
	}
	var err error
	if isExternalClient {
		err = res.establishExternalClient()
	} else {
		err = res.establishInternalClient()
	}
	return res, err
}

func (k *K8sCRDAdaptor) LoadPolicy(model model.Model) error {
	policyObj, err := k.getPolicyObject()
	if err != nil {
		return err
	}
	splits := strings.Split(policyObj.Spec.PolicyItem, "\n")
	for _, line := range splits {
		if line != "" {
			persist.LoadPolicyLine(line, model)
		}
	}
	return nil
}

func (k *K8sCRDAdaptor) SavePolicy(model model.Model) error {
	var buffer strings.Builder
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			buffer.WriteString(fmt.Sprintf("%s,%s\n", ptype, util.ArrayToString(rule)))
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			buffer.WriteString(fmt.Sprintf("%s,%s\n", ptype, util.ArrayToString(rule)))
		}
	}
	err := k.updatePoliyObject(buffer.String())
	return err
}

func (k *K8sCRDAdaptor) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (k *K8sCRDAdaptor) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (k *K8sCRDAdaptor) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

func (k *K8sCRDAdaptor) establishInternalClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *K8sCRDAdaptor) establishExternalClient() error {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *K8sCRDAdaptor) getPolicyObject() (*k8sauthzv1.CasbinPolicy, error) {
	obj, err := k.clientset.AuthV1().CasbinPolicies(k.namespace).Get(
		context.TODO(),
		k.modelName,
		metav1.GetOptions{},
	)
	return obj, err
}

func (k *K8sCRDAdaptor) updatePoliyObject(policy string) error {
	oldObj, err := k.clientset.AuthV1().CasbinPolicies(k.namespace).Get(
		context.TODO(),
		k.modelName,
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}
	oldObj.Spec.PolicyItem = policy
	_, err = k.clientset.AuthV1().CasbinPolicies(k.namespace).Update(
		context.TODO(),
		oldObj,
		metav1.UpdateOptions{},
	)
	return err
}
