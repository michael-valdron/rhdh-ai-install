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

// ClientOptions holds the configuration options for creating a new client instance.
type ClientOptions struct {
	Verbose        bool          // Enable verbose logging
	DryRun         bool          // Run the client in dry-run mode
	Timeout        time.Duration // Timeout for API calls
	RhdhOperator   bool          // Indicates whether the RHDH operator is being used for deployments
	CurrentContext string        // Current kubeconfig context to use
	Namespace      string        // Namespace to operate in, if not specified it will use the one from the kubeconfig
}

// Client represents a Kubernetes client entity. It contains all necessary components for interacting with the Kubernetes API.
type Client struct {
	clientset kubernetes.Interface      // The main client interface for interacting with Kubernetes resources
	crdclient apiextclientset.Interface // Interface for working with Custom Resource Definitions (CRDs)
	Namespace string                    // The namespace to operate in, as specified by the user or derived from kubeconfig
	Options   ClientOptions             // Configuration options for the client instance
}

// ClientParams holds the override parameters for individual client commands.
type ClientParams struct {
	Namespace string // Namespace to target when executing a command
}

// DryRunRoundTripper is a custom http.RoundTripper that adds 'dryRun=All' query parameter to mutating requests (POST, PUT, PATCH, DELETE) for dry-run mode.
type DryRunRoundTripper struct {
	Wrapped http.RoundTripper // The underlying RoundTripper to use for request execution
}

// Provides namespace for the client to use when executing commands. Uses namespace override if provided otherwise use client
// set namespace.
func (c *Client) initNamespaceWithParams(params ClientParams) string {
	namespace := c.Namespace

	if params.Namespace != "" {
		namespace = params.Namespace
	}

	return namespace
}

// Returns an interface to manage Kubernetes Deployment resources in the specified namespace.
func (c *Client) Deployments(params ClientParams) appsv1.DeploymentInterface {
	return c.clientset.AppsV1().Deployments(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Kubernetes Namespace resources.
func (c *Client) Namespaces() corev1.NamespaceInterface {
	return c.clientset.CoreV1().Namespaces()
}

// Returns an interface to manage Kubernetes ConfigMap resources in the specified namespace.
func (c *Client) ConfigMaps(params ClientParams) corev1.ConfigMapInterface {
	return c.clientset.CoreV1().ConfigMaps(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Kubernetes Secret resources in the specified namespace.
func (c *Client) Secrets(params ClientParams) corev1.SecretInterface {
	return c.clientset.CoreV1().Secrets(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Kubernetes ServiceAccount resources in the specified namespace.
func (c *Client) ServiceAccounts(params ClientParams) corev1.ServiceAccountInterface {
	return c.clientset.CoreV1().ServiceAccounts(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Kubernetes ClusterRole resources.
func (c *Client) ClusterRoles() rbacv1.ClusterRoleInterface {
	return c.clientset.RbacV1().ClusterRoles()
}

// Returns an interface to manage Kubernetes ClusterRoleBinding resources.
func (c *Client) ClusterRoleBindings() rbacv1.ClusterRoleBindingInterface {
	return c.clientset.RbacV1().ClusterRoleBindings()
}

// Returns an interface to manage Kubernetes RoleBinding resources in the specified namespace.
func (c *Client) RoleBindings(params ClientParams) rbacv1.RoleBindingInterface {
	return c.clientset.RbacV1().RoleBindings(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Kubernetes Role resources in the specified namespace.
func (c *Client) Roles(params ClientParams) rbacv1.RoleInterface {
	return c.clientset.RbacV1().Roles(c.initNamespaceWithParams(params))
}

// Returns an interface to manage Custom Resource Definitions (CRDs).
func (c *Client) CustomResourceDefinitions() apiextv1.CustomResourceDefinitionInterface {
	return c.crdclient.ApiextensionsV1().CustomResourceDefinitions()
}

// GetUnstructured retrieves an Unstructured resource object from the Kubernetes API using the provided context, GroupVersionResource, name, and option parameters.
func (c *Client) GetUnstructured(ctx context.Context, gvr *schema.GroupVersionResource, name string, params ClientParams, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	// Initialize a new request object for GET operation using the clientset's RESTClient.
	req := c.clientset.CoreV1().RESTClient().Get()

	// Create an Unstructured object to hold the retrieved resource data.
	spec := new(unstructured.Unstructured)

	var result rest.Result

	// Use the namespace override if provided, otherwise use the client set namespace.
	if params.Namespace != "" {
		req = req.Namespace(params.Namespace) // Apply namespace if specified in params.
	} else {
		req = req.Namespace(c.Namespace) // Otherwise, use the client's namespace.
	}

	// Construct the request URL using the GroupVersionResource and name of the resource to fetch.
	req = req.AbsPath(fmt.Sprintf("/api/%s/%s", gvr.Group, gvr.Version)).
		Resource(gvr.Resource).Name(name)

	// Execute the GET request with the provided context.
	result = req.Do(ctx)

	// Check if there was an error during the API call and return it if any.
	if err := result.Error(); err != nil {
		return nil, err // Return the error if encountered.
	}

	// Decode the response into the Unstructured object.
	if err := result.Into(spec); err != nil {
		return nil, err // Return the decoding error if any.
	}

	return spec, nil // Return the retrieved Unstructured resource object on success.
}

// PatchUnstructured applies a strategic merge patch to an Unstructured resource using the Kubernetes clientset.
// It sends a PATCH request with the specified options and returns any encountered error during the process.
func (c *Client) PatchUnstructured(ctx context.Context, gvr *schema.GroupVersionResource, resourcePatch *unstructured.Unstructured, opts metav1.PatchOptions) error {
	// Initialize variables for patch data and request.
	var (
		patchData []byte
		err       error
		req       *rest.Request
	)

	// Marshal the Unstructured resource patch into JSON format.
	if patchData, err = resourcePatch.MarshalJSON(); err != nil {
		return err // Return the marshaling error if any.
	}

	// Construct a PATCH request using the clientset's RESTClient.
	req = c.clientset.
		CoreV1().
		RESTClient().
		Patch(types.StrategicMergePatchType).                       // Specify strategic merge patch type
		AbsPath(fmt.Sprintf("/api/%s/%s", gvr.Group, gvr.Version)). // Set the API group and version path
		Resource(gvr.Resource).                                     // Set the resource type from the provided GroupVersionResource
		Name(resourcePatch.GetName()).                              // Set the name of the resource to patch
		Body(patchData)                                             // Attach the marshalled patch data

	// Apply namespace if specified in the resource patch.
	if resourcePatch.GetNamespace() != "" {
		req = req.Namespace(resourcePatch.GetNamespace()) // Add namespace to the request if provided
	}

	// Execute the PATCH request with the given context and return any error encountered during execution.
	_, err = req.Do(ctx).Raw()

	return err // Return the error from the request execution if any.
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

// Creates a Kubernetes client entity based on provided configuration options.
func New(opts ClientOptions) (*Client, error) {
	// Initialize a new instance of Client struct with the given options.
	client := &Client{
		Namespace: opts.Namespace,
		Options:   opts,
	}

	var (
		// Initialize loading rules and config overrides using non-interactive deferred loading for kubeconfig.
		kubeconfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(createLoadingRules(), createConfigOverrides(opts))
		config     *rest.Config
		err        error
	)

	// If no namespace is given in the options, attempt to load the namespace from the kubeconfig.
	if client.Namespace == "" {
		if client.Namespace, _, err = kubeconfig.Namespace(); err != nil {
			return nil, fmt.Errorf("error loading kubeconfig namespace: %v", err)
		}
	}

	// Load the rest config from the kubeconfig.
	if config, err = kubeconfig.ClientConfig(); err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %v", err)
	}

	// If dry-run mode is enabled, modify the transport to add 'dryRun=All' query parameter for mutating requests (POST, PUT, PATCH, DELETE).
	if opts.DryRun {
		if restClient, err := rest.RESTClientFor(config); err != nil {
			return nil, fmt.Errorf("error creating REST client: %v", err)
		} else {
			config.Transport = &DryRunRoundTripper{
				Wrapped: restClient.Client.Transport,
			}
		}
	}

	// Create a new Kubernetes clientset using the built config.
	if client.clientset, err = kubernetes.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("error creating clientset from built config: %v", err)
	}

	// Create a new API extensions (CRDs) clientset using the built config.
	if client.crdclient, err = apiextclientset.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("error creating clientset for API extensions (CRDs) from built config: %v", err)
	}

	return client, nil
}

// Creates a mock k8s client entity with provided objects and API extensions objects.
func NewMock(opts ClientOptions, objects []runtime.Object, apiExtObjects []runtime.Object) (*Client, error) {
	// Initialize a new instance of Client struct with the given options.
	client := &Client{
		Namespace: opts.Namespace,
		Options:   opts,
	}

	// Create a fake clientset using the provided objects for core Kubernetes resources.
	client.clientset = fake.NewClientset(objects...)

	// Create a fake API extensions (CRDs) clientset using the provided API extension objects.
	client.crdclient = apiextfakeclientset.NewClientset(apiExtObjects...)

	return client, nil
}
