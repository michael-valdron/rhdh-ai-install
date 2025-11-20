package serviceaccount

import (
	"context"
	"fmt"

	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/client"
	"github.com/redhat-ai-dev/rhdh-ai-install/pkg/util"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This struct holds references to various Kubernetes resources required for setting up the model catalog service account with associated roles,
// role bindings, cluster roles, and cluster role bindings. It also includes a reference to a secret used for token generation.
type ServiceAccountSuite struct {
	ServiceAccount       *corev1.ServiceAccount     // Represents the service account resource
	ClusterRole          *rbacv1.ClusterRole        // Represents the cluster role resource
	ClusterRoleBinding   *rbacv1.ClusterRoleBinding // Represents the cluster role binding resource
	RoleBinding          *rbacv1.RoleBinding        // Represents the role binding resource for regular namespaces
	DashboardRoleBinding *rbacv1.RoleBinding        // Represents the role binding resource for RHOAI dashboard namespace
	Role                 *rbacv1.Role               // Represents the role resource
	Secret               *corev1.Secret             // Represents the secret resource used for token generation
}

// Create creates all necessary Kubernetes resources for setting up the model catalog service account with associated roles, role bindings, cluster roles, and cluster role bindings.
// It also includes creating a secret resource used for token generation.
func (sa *ServiceAccountSuite) Create(ctx context.Context, c *client.Client, opts metav1.CreateOptions) error {
	// Create the service account resource
	if _, err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Create(ctx, sa.ServiceAccount, opts); err != nil {
		return err
	}

	// Create the cluster role resource
	if _, err := c.ClusterRoles().Create(ctx, sa.ClusterRole, opts); err != nil {
		return err
	}

	// Create the cluster role binding resource
	if _, err := c.ClusterRoleBindings().Create(ctx, sa.ClusterRoleBinding, opts); err != nil {
		return err
	}

	// Create the role binding resource for regular namespaces
	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Create(ctx, sa.RoleBinding, opts); err != nil {
		return err
	}

	// Create the role binding resource for RHOAI dashboard namespace
	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Create(ctx, sa.DashboardRoleBinding, opts); err != nil {
		return err
	}

	// Create the role resource
	if _, err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Create(ctx, sa.Role, opts); err != nil {
		return err
	}

	// Create the secret resource used for token generation
	if _, err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Create(ctx, sa.Secret, opts); err != nil {
		return err
	}

	return nil
}

// Delete removes all necessary Kubernetes resources for setting up the model catalog service account with associated roles, role bindings, cluster roles, and cluster role bindings.
// It also includes deleting a secret resource used for token generation.
func (sa *ServiceAccountSuite) Delete(ctx context.Context, c *client.Client, opts metav1.DeleteOptions) error {
	// Delete the secret resource used for token generation
	if err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Delete(ctx, sa.Secret.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the role resource
	if err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Delete(ctx, sa.Role.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the role binding resource for RHOAI dashboard namespace
	if err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Delete(ctx, sa.DashboardRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the role binding resource for regular namespaces
	if err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Delete(ctx, sa.RoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the cluster role binding resource
	if err := c.ClusterRoleBindings().Delete(ctx, sa.ClusterRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the cluster role resource
	if err := c.ClusterRoles().Delete(ctx, sa.ClusterRole.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	// Delete the service account resource
	if err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Delete(ctx, sa.ServiceAccount.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

// Check verifies that there are no conflicts with setting up the model catalog service account including the tied roles, role bindings, cluster roles, and cluster role bindings for the service account.
// It also checks if a secret resource used for token generation does not exist.
func (sa *ServiceAccountSuite) Check(ctx context.Context, c *client.Client, opts metav1.GetOptions) error {
	errs := []any{}

	// Verify that the model catalog service account does not already exist and is ready for install.
	if _, err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Get(ctx, sa.ServiceAccount.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account '%s' already exists and is not ready for install", sa.ServiceAccount.Name))
	}

	// Verify that the model catalog cluster role does not already exist and is ready for install.
	if _, err := c.ClusterRoles().Get(ctx, sa.ClusterRole.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account cluster role already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account cluster role '%s' already exists and is not ready for install", sa.ClusterRole.Name))
	}

	// Verify that the model catalog cluster role binding does not already exist and is ready for install.
	if _, err := c.ClusterRoleBindings().Get(ctx, sa.ClusterRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account cluster role binding already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account cluster role binding '%s' already exists and is not ready for install", sa.ClusterRoleBinding))
	}

	// Verify that the model catalog role binding for regular namespaces does not already exist and is ready for install.
	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Get(ctx, sa.RoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account role binding already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account role binding '%s' already exists and is not ready for install", sa.RoleBinding.Name))
	}

	// Verify that the model catalog service account token secret does not already exist and is ready for install.
	if _, err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Get(ctx, sa.Secret.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account token secret already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account token secret '%s' already exists and is not ready for install", sa.Secret.Name))
	}

	// Verify that the model catalog service account role does not already exist and is ready for install.
	if _, err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Get(ctx, sa.Role.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The model catalog service account role already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("model catalog service account role '%s' already exists and is not ready for install", sa.Role.Name))
	}

	// Verify that the RHOAI dashboard role binding does not already exist and is ready for install.
	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Get(ctx, sa.DashboardRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		// The RHOAI dashboard role binding already exists and is not ready for install.
		errs = append(errs, fmt.Errorf("RHOAI dashboard role binding '%s' already exists and is not ready for install", sa.DashboardRoleBinding.Name))
	}

	if len(errs) != 0 {
		if c.Options.Verbose {
			return fmt.Errorf("model catalog service account resources are not ready, see the following errors (note: remember uninstall first if there is a previous install):\n%s", util.Join(errs, "\n"))
		} else {
			return fmt.Errorf("model catalog service account resources are not ready")
		}
	} else {
		return nil
	}
}

// GetServiceAccountTokenName generates a token secret name for a given service account name.
func GetServiceAccountTokenName(serviceAccountName string) string {
	return fmt.Sprintf("%s-token", serviceAccountName)
}

// New initializes and returns a new ServiceAccountSuite instance with the specified namespaces.
func New(rhdhNamespace string, rhoaiNamespace string) *ServiceAccountSuite {
	metadata := metav1.ObjectMeta{
		Name:      ServiceAccountName,
		Namespace: rhdhNamespace,
	}
	dashboardMetadata := metav1.ObjectMeta{
		Name:      DashboardRoleBindingName,
		Namespace: rhoaiNamespace,
	}

	return &ServiceAccountSuite{
		ServiceAccount: &corev1.ServiceAccount{
			ObjectMeta: metadata,
		},
		ClusterRole: &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: metadata.Name,
				Annotations: map[string]string{
					"argocd.argoproj.io/sync-wave": "0",
				},
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"apiextensions.k8s.io"},
					Resources: []string{"customresourcedefinitions"},
					Verbs:     []string{"get"},
				},
				{
					APIGroups: []string{"route.openshift.io"},
					Resources: []string{"routes"},
					Verbs:     []string{"get", "list", "watch"},
				},
				{
					APIGroups: []string{""},
					Resources: []string{"serviceaccounts", "services"},
					Verbs:     []string{"get", "list", "watch"},
				},
				{
					APIGroups: []string{"serving.kserve.io"},
					Resources: []string{"inferenceservices"},
					Verbs:     []string{"get", "list", "watch"},
				},
			},
		},
		ClusterRoleBinding: &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: metadata.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     metadata.Name,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      metadata.Name,
					Namespace: metadata.Namespace,
				},
			},
		},
		RoleBinding: &rbacv1.RoleBinding{
			ObjectMeta: metadata,
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     metadata.Name,
			},
		},
		DashboardRoleBinding: &rbacv1.RoleBinding{
			ObjectMeta: dashboardMetadata,
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     DashboardRoleName,
			},
			Subjects: []rbacv1.Subject{
				{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     rbacv1.GroupKind,
					Name:     fmt.Sprintf("system:serviceaccounts:%s", metadata.Namespace),
				},
			},
		},
		Role: &rbacv1.Role{
			ObjectMeta: metadata,
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"configmaps"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
				},
			},
		},
		Secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      GetServiceAccountTokenName(metadata.Name),
				Namespace: metadata.Namespace,
				Annotations: map[string]string{
					corev1.ServiceAccountNameKey: metadata.Name,
				},
			},
			Type: corev1.SecretTypeServiceAccountToken,
		},
	}
}
