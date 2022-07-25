# K8s-gatekeeper
[TOC]
## 1.Overview
### 1.1 What is K8s-gatekeeper
K8s-gatekeeper is an admission webhook for k8s, using [Casbin](https://casbin.org/docs/en/overview) to apply arbitrary user-defined access control rules to help prevent any operation on k8s which administrator doesn't want.

Casbin is a powerful and efficient open-source access control library. It provides support for enforcing authorization based on various access control models. For more detail about Casbin, see <https://casbin.org/docs/en/overview>.

Admission webhooks in K8s are HTTP callbacks that receive 'admission requests' and do something with them. In particular, K8s-gatekeeper is a special type of admission webhoook: 'ValidatingAdmissionWebhook', which can decide whether to accept or reject this admission request or not. As for admission requests, they are HTTP requests describing an operation on specified resources of K8s (for example, creating/deleting a deployment). For more about admission webhooks, see K8s official doc <https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#what-are-admission-webhooks>

### 1.2 An example illustrating how it works.

![](/doc/overview.png)

For example, when somebody wants to create a deployment containing a pod running nginx (using kubectl or k8s clients), K8s will generate an admission request, which (if translated into yaml format) can be something like this.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.1
        ports:
        - containerPort: 80
```

This request will go through the process of all the middleware shown in the picture, including our K8s-gatekeeper. K8s-gatekeeper can detected all the Casbin enforcers stored in K8s's etcd, which is created and maintained by user(via kubectl or go-client we provide). Each enforcer contains a Casbin model and a Casbin policy. The admission request will be processed by every enforcer, one by one, and only by passing all enforcers can a request be accepted by this K8s-gatekeeper.

(If you do not understand what is Casbin enforcer, model or policy, see this document <https://casbin.org/docs/en/get-started>)

For example, for some reason, the administrator want to forbid the apperance of image 'nginx:1.14.1' while allowing  'nginx:1.3.1', an enforcer containing the following rule and policy can be created: (We will explain how to create an enforcer, what these models and policies and how to write them in following chapters.)

model:
```
[request_definition]
r =  obj

[policy_definition]
p =  obj,eft

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="deployments" && \
access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Image") == p.obj

```

policy: 
```
p, "nginx:1.13.1",allow
p, "nginx:1.14.1",deny
```

By creating an enforcer containg model and policy above, the previous admission request will be reject by this enforcer, which means K8s won't create this deployment.


## 2 Install K8s-gatekeeper

Three methods are provided for installing K8s-gatekeeper: External webhook, Internal webhook and helm.

*Note: these methods are only for user to try K8s-gatekeeper, and it is not secure. If you want to use it in productive environment, please make sure you read Chapter 5. Advanced setting and make modifications accordingly when necessary before installation  *

### 2.1 Install K8s-gatekeeper via helm

*Will be worked out soon*

### 2.2 Internal webhook
#### 2.2.1 Step 1: Build image

Internal webhook means the webhook itself will be implmented as a service inside k8s. Creating a service as well as deployment requires a image of K8s-gatekeeper. You can choose to build your own image, or use the pre-built image we provide.

##### 2.2.1.1 Using Prebuild image
*Will be worked out soon*

##### 2.2.1.2 Build K8s-gatekeeper image 
Run 
```shell
docker build --target webhook -t k8s-gatekeeper .
```
Then there will be a local image called 'k8s-gatekeeper:latest'.

*Note: if you are using minikube, please execute `eval $(minikube -p minikube docker-env)` before running docker build*

#### 2.2.2 Step 2: Set up services and deployments for K8s-gatekeeper
Run following commands
```shell
kubectl apply -f config/rbac.yaml
kubectl apply -f config/webhook_deployment.yaml 
kubectl apply -f config/webhook_internal.yaml 
```
Soon K8s-gatekeeper should be running, and you can use `kubectl get pods` to confirm that.

#### 2.2.3 Step3: Install Crd Resources for K8s-gatekeeper
Run following commands
```shell
kubectl apply -f config/auth.casbin.org_casbinmodels.yaml 
kubectl apply -f config/auth.casbin.org_casbinpolicies.yaml
```

### 2.3 External webhook
External webhook means K8s-gatekeeper will be running outside the K8s, and K8s will visit K8s-gatekeeper like visiting a ordinary website. K8s has mandatory requirement that admission webhook must be  HTTPS. For the sake of user's experience in trying  K8s-gatekeeper, we have provided you a set of certificate as well as private key (though it is not secure). If you prefer to use your own certificate, please refer to Chapter 5. Advanced setting to make adjustments to the certificate and private key.

The certificate we provide is issued for 'webhook.domain.local', so please modify the host (like /etc/hosts), point webhook.domain.local to the ip address on which K8s-gatekeeper is running.

Then execute
```
go mod tidy
go mod vendor
go run cmd/webhook/main.go
kubectl apply -f config/auth.casbin.org_casbinmodels.yaml 
kubectl apply -f config/auth.casbin.org_casbinpolicies.yaml
kubectl apply -f config/webhook_external.yaml 
```
## 3. Try K8s-gatekeeper

### 3.1 Create Casbin Model and Policy
You have 2 methods to create a model and policy: via kubectl or via go-client we provide.

#### 3.1.1 Create/Update Casbin Model and Policy via kubectl
In K8s-gatekeeper, Casbin model is stored in a kind of CRD Resource called 'CasbinModel'. Its definition is located in config/auth.casbin.org_casbinmodels.yaml

There are examples in example/allowed_repo/model.yaml. You are supposed to pay attention to the following fields:
- metadata.name: name of the model. This name MUST be same with the name of CasbinPolicy object related to this model, so that K8s-gatekeeper can pair them and create an enforcer.
- spec.enable: if this field is set to "false", this model(as well as CasbinPolicy object related to this model) will be ignored.
- spec.modelText: a string which contains the model text of a casbin model. 

Casbin Policy is stored in another kind of CRD Resource called 'CasbinPolicy', whose definition can be found in config/auth.casbin.org_casbinpolicies.yaml

There are examples in example/allowed_repo/policy.yaml. You are supposed to pay attention to the following fields:
- metadata.name: name of the policy. This name MUST be same with the name of CasbinModel object related to this policy, so that K8s-gatekeeper can pair them and create an enforcer.
- spec.policyItem: a string which contains the policy text of a casbin model.
  
After creating your own CasbinModel and CasbinPolicy files, use 
```shell
kubectl apply -f <filename>
```
to put them into effect.

Once a pair of CasbinModel and CasbinPolicy is created, within at most 5 seconds K8s-gatekeeper will be able to detect it.

#### 3.1.2 Create /Updata Casbin Model and Policy via go-client we provide
It has been taken into consideration that there may be situation that it is not convenient to use shell to execute command directly on a node of K8s cluster, for example, when you are building a automatic cloud platform for your corporation. Therefore we have developed a go-client to create maintain CasbinModel and CasbinPolicy.

The go-client library is located in pkg/client. 

In client.go we provide function to create a client.
```
func NewK8sGateKeeperClient(externalClient bool) (*K8sGateKeeperClient, error) 
```
parameter externalClient means whether K8s-gatekeeper is running inside the K8s cluster or not.


In model.go we provide various functions to create/delete/modify CasbinModel. You can find out how to use there interfaces in model_test.go.

In policy.go we provide various functions to create/delete/modify CasbiPolicy. You can find out how to use there interfaces in policy_test.go.

### 3.1.2 Try whether K8s-gatekeeper works

Suppose you have already created exactly the model and policy in example/allowed_repo, now try this 
```
kubectl apply -f example/allowed_repo/testcase/reject_1.yaml
```

you are supposed to find that k8s will reject this request, and mentioning that this webhook was the reason why this request is rejected. However, when you tries to apply example/allowed_repo/testcase/approve_2.yaml, it will be accepted.

## 4. How to write Model and Policy K8s-gatekeeper
First of all, you are supposed to know the basic grammar of Casbin Models and Policies. If you haven't acknowledged it, please read <https://casbin.org/docs/en/get-started> first. In this chapter we will assume that you have known what are Casbin Models and Policies.

### 4.1 Request Definition of Model
When K8s-gatekeeper is authorizing a request, the input is always one object: the go object of the Admission Request. Which means the enforcer will always be used like this
```golang
  ok, err := enforcer.Enforce(admission)
```
in which admission is an `AdmissionReview` object defined by K8s's official go api `"k8s.io/api/admission/v1"`. You can see the definition of this struct is this repository <https://github.com/kubernetes/api/blob/master/admission/v1/types.go>. Or see <https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#webhook-request-and-response> for more information

Therefore for any model used by K8s-gatekeeper, the definitiion of request_definition should always be like this
```
    [request_definition]
    r =  obj
```

Name 'obj' is not mandatory, as long as the name is consistent with the name used in `[matchers]` part.

### 4.2 Matchers of Model
You are supposed to use the ABAC feature of Casbin to write down your rule. However, the expression evaluator integrated in Casbin supports neither indexing in masp or arrays(slices), nor the expansion of array. Therefore K8s-gatekeeper provide various 'Casbin functions' as extension to impelement these features. If you still find that your demand cannot be fulfilled by these extensions, it is welcomed to start a issue, or pr directly.

If you don't know what is casbin funtion, see <https://casbin.org/docs/en/function> for more information.

Here are the extension functions
### 4.2.1 Externsion functions
#### 4.2.1.1 access
Access is used to solve the problem that Casbin doesn't support indexing in map or array. `example/allowed_repo/model.yaml` is the example of this function
```
[matchers]
m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="deployments" && \
access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Image") == p.obj
```
In this matcher, access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Image") is equal to `r.obj.Request.Object.Object.Spec.Template.Spec.Containers[0].Image`, in which `r.obj.Request.Object.Object.Spec.Template.Spec.Containers` is obviously a slice.

Access is also able to call simple funtion which has not parameters and one single return value. `example/container_resource_limit/model.yaml` is an example.

```
[matchers]
  m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="deployments" && \
  parseFloat(access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Resources","Limits","cpu","Value")) >= parseFloat(p.cpu) && \
  parseFloat(access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Resources","Limits","memory","Value")) >= parseFloat(p.memory)
```

In this matcher, `access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Resources","Limits","cpu","Value")` is equal to `r.obj.Request.Object.Object.Spec.Template.Spec.Containers[0].Resources.Limits["cpu"].Value()`, where `r.obj.Request.Object.Object.Spec.Template.Spec.Containers[0].Resources.Limits` is a map, and `Value()` is a simple funtion which has not parameters and one single return value.

#### 4.2.1.2 accessWithWildcard
Sometimes it is natural to have demand like this: all elements in an array must have prefix "aaa". However, Casbin doesn't support `for` loop. However with `accessWithWildcard` and the "map/slice expansion" feature, such demand can be easily implemented. 

For example, suppose `a.b.c` is an array `[aaa,bbb,ccc,ddd,eee]`, then result of `accessWithWildcard(a,"b","c","*")` will be a slice `[aaa,bbb,ccc,ddd,eee]`. See? with wildcard `*` this slice is expanded.


Similarly, wildcard can be used more than once. For example, result of `accessWithWildcard(a,"b","c","*","*")` will be `[a.b.c[0][0], a.b.c[0][1]... a.b.c[1][0], a.b.c[1][1]...]`

### 4.2.1.3 Functions Supporting Variable-length Argument
In the expression evaluator of Casbin, when a parameter is an array, it will be automatically expanded  as the variable-length argument. Utilizing this feature to support the array/slice/map expansion, we also integrated serveral functions accepting an array/slice as parameter.

- contain(), accept multiple parameters, and returns whether there is an parameter other than the last parameter equals the last parameter 
- split(a,b,c...,sep,index) it returns a slice which contains`[splits(a,sep)[index], splits(b,sep)[index], splits(a,sep)[index]...]`
- len() return the length of the variable-length argument
- matchRegex(a,b,c...regex) return whether a,b,c... all of them matches the given regex

Here is an example in `example/disallowed_tag/model.yaml`
```
    [matchers]
    m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="deployments" && \
    contain(split(accessWithWildcard(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , "*", "Image"),":",1) , p.obj)
```

Assume `accessWithWildcard(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , "*", "Image")`returns `["a:b", "c:d", "e:f", "g:h"]` then because splits supports variable-length argument, and splits operation is applied on each element, and eventualy element whose index is 1 will be selected and return, so `split(accessWithWildcard(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , "*", "Image"),":",1)` returns `["b","d","f","h"]`. And `contain(split(accessWithWildcard(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , "*", "Image"),":",1) , p.obj)` returns whether`p.obj` is contained in `["b","d","f","h"]`

#### 4.2.1.2 Type conversion functions
- ParseFloat(): parse an integer to a float. (It is because that any number in comparsion must be converted into float).
- ToString(): convert an object to string. This object must have a basic type of string. (for example, an object of type `XXX` when there is a statement `type XXX string`)
- IsNil(): return whether the parameter is nil

## 5. Advanced Settings
To be continued




