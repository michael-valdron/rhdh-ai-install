package client

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	apiextclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextfakeclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ClientOptions struct {
	Verbose bool
	DryRun  bool
	Timeout time.Duration

	RhdhOperator   bool
	CurrentContext string
	Namespace      string
	OutputFormat   string
}

type Client struct {
	clientset kubernetes.Interface
	crdclient apiextclientset.Interface
	Namespace string
	Options   ClientOptions
}

type ClientParams struct {
	Namespace string
}

type DryRunRoundTripper struct {
	Wrapped http.RoundTripper
}

func (c *Client) initNamespaceWithParams(params ClientParams) string {
	namespace := c.Namespace

	if params.Namespace != "" {
		namespace = params.Namespace
	}

	return namespace
}

func (c *Client) Deployments(params ClientParams) appsv1.DeploymentInterface {
	return c.clientset.AppsV1().Deployments(c.initNamespaceWithParams(params))
}

func (c *Client) Namespaces() corev1.NamespaceInterface {
	return c.clientset.CoreV1().Namespaces()
}

func (c *Client) ConfigMaps(params ClientParams) corev1.ConfigMapInterface {
	return c.clientset.CoreV1().ConfigMaps(c.initNamespaceWithParams(params))
}

func (c *Client) Secrets(params ClientParams) corev1.SecretInterface {
	return c.clientset.CoreV1().Secrets(c.initNamespaceWithParams(params))
}

func (c *Client) ServiceAccounts(params ClientParams) corev1.ServiceAccountInterface {
	return c.clientset.CoreV1().ServiceAccounts(c.initNamespaceWithParams(params))
}

func (c *Client) ClusterRoles() rbacv1.ClusterRoleInterface {
	return c.clientset.RbacV1().ClusterRoles()
}

func (c *Client) ClusterRoleBindings() rbacv1.ClusterRoleBindingInterface {
	return c.clientset.RbacV1().ClusterRoleBindings()
}

func (c *Client) RoleBindings(params ClientParams) rbacv1.RoleBindingInterface {
	return c.clientset.RbacV1().RoleBindings(c.initNamespaceWithParams(params))
}

func (c *Client) Roles(params ClientParams) rbacv1.RoleInterface {
	return c.clientset.RbacV1().Roles(c.initNamespaceWithParams(params))
}

func (c *Client) CustomResourceDefinitions() apiextv1.CustomResourceDefinitionInterface {
	return c.crdclient.ApiextensionsV1().CustomResourceDefinitions()
}

func (c *Client) GetUnstructured(ctx context.Context, gvr *schema.GroupVersionResource, name string, params ClientParams, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	req := c.clientset.CoreV1().RESTClient().Get()
	spec := new(unstructured.Unstructured)
	var result rest.Result

	if params.Namespace != "" {
		req = req.Namespace(params.Namespace)
	} else {
		req = req.Namespace(c.Namespace)
	}

	req = req.
		AbsPath(fmt.Sprintf("/api/%s/%s", gvr.Group, gvr.Version)).
		Resource(gvr.Resource).
		Name(name)

	result = req.Do(ctx)
	if err := result.Error(); err != nil {
		return nil, err
	}

	if err := result.Into(spec); err != nil {
		return nil, err
	}

	return spec, nil
}

func (c *Client) PatchUnstructured(ctx context.Context, gvr *schema.GroupVersionResource, resourcePatch *unstructured.Unstructured, opts metav1.PatchOptions) error {
	var (
		patchData []byte
		err       error
		req       *rest.Request
	)

	if patchData, err = resourcePatch.MarshalJSON(); err != nil {
		return err
	}

	req = c.clientset.
		CoreV1().
		RESTClient().
		Patch(types.StrategicMergePatchType).
		AbsPath(fmt.Sprintf("/api/%s/%s", gvr.Group, gvr.Version)).
		Resource(gvr.Resource).
		Name(resourcePatch.GetName()).
		Body(patchData)
	if resourcePatch.GetNamespace() != "" {
		req = req.Namespace(resourcePatch.GetNamespace())
	}
	_, err = req.Do(ctx).Raw()

	return err
}

// RoundTrip implements the http.RoundTripper interface.
func (d *DryRunRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only apply dryRun to mutating methods
	if req.Method == http.MethodPost || req.Method == http.MethodPut ||
		req.Method == http.MethodPatch || req.Method == http.MethodDelete {

		// Add dryRun=All query parameter
		query := req.URL.Query()
		query.Set("dryRun", "All")
		req.URL.RawQuery = query.Encode()
	}

	return d.Wrapped.RoundTrip(req)
}

// Creates client config loading rules with the home kubeconfig guaranteed
func createLoadingRules() *clientcmd.ClientConfigLoadingRules {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if home := homedir.HomeDir(); home != "" {
		loadingRules.Precedence = append(loadingRules.Precedence, filepath.Join(home, ".kube", "config"))
	}
	return loadingRules
}

// create client config overrides based on passed options
func createConfigOverrides(opts ClientOptions) *clientcmd.ConfigOverrides {
	overrides := &clientcmd.ConfigOverrides{}

	if opts.CurrentContext != "" {
		overrides.CurrentContext = opts.CurrentContext
	}

	if opts.Namespace != "" {
		overrides.Context.Namespace = opts.Namespace
	}

	if opts.Timeout != time.Duration(0) {
		overrides.Timeout = opts.Timeout.String()
	}

	return overrides
}

// Creates a k8s client entity
func New(opts ClientOptions) (*Client, error) {
	client := &Client{
		Options: opts,
	}

	var (
		kubeconfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(createLoadingRules(), createConfigOverrides(opts))
		config     *rest.Config
		err        error
	)

	if client.Namespace, _, err = kubeconfig.Namespace(); err != nil {
		return nil, fmt.Errorf("error loading kubeconfig namespace: %v", err)
	}

	if config, err = kubeconfig.ClientConfig(); err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %v", err)
	}

	if opts.DryRun {
		if restClient, err := rest.RESTClientFor(config); err != nil {
			return nil, fmt.Errorf("error creating REST client: %v", err)
		} else {
			config.Transport = &DryRunRoundTripper{
				Wrapped: restClient.Client.Transport,
			}
		}
	}

	if client.clientset, err = kubernetes.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("error creating clientset from built config: %v", err)
	}

	if client.crdclient, err = apiextclientset.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("error creating clientset for API extensions (CRDs) from built config: %v", err)
	}

	return client, nil
}

// Creates a mock k8s client entity
func NewMock(opts ClientOptions, objects []runtime.Object, apiExtObjects []runtime.Object) (*Client, error) {
	return &Client{
		clientset: fake.NewClientset(objects...),
		crdclient: apiextfakeclientset.NewClientset(apiExtObjects...),
		Namespace: opts.Namespace,
		Options:   opts,
	}, nil
}
