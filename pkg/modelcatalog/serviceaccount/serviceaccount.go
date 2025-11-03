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

type ServiceAccountSuite struct {
	ServiceAccount       *corev1.ServiceAccount
	ClusterRole          *rbacv1.ClusterRole
	ClusterRoleBinding   *rbacv1.ClusterRoleBinding
	RoleBinding          *rbacv1.RoleBinding
	DashboardRoleBinding *rbacv1.RoleBinding
	Role                 *rbacv1.Role
	Secret               *corev1.Secret
}

func (sa *ServiceAccountSuite) Create(ctx context.Context, c *client.Client, opts metav1.CreateOptions) error {
	if _, err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Create(ctx, sa.ServiceAccount, opts); err != nil {
		return err
	}

	if _, err := c.ClusterRoles().Create(ctx, sa.ClusterRole, opts); err != nil {
		return err
	}

	if _, err := c.ClusterRoleBindings().Create(ctx, sa.ClusterRoleBinding, opts); err != nil {
		return err
	}

	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Create(ctx, sa.RoleBinding, opts); err != nil {
		return err
	}

	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Create(ctx, sa.DashboardRoleBinding, opts); err != nil {
		return err
	}

	if _, err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Create(ctx, sa.Role, opts); err != nil {
		return err
	}

	if _, err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Create(ctx, sa.Secret, opts); err != nil {
		return err
	}

	return nil
}

func (sa *ServiceAccountSuite) Delete(ctx context.Context, c *client.Client, opts metav1.DeleteOptions) error {
	if err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Delete(ctx, sa.Secret.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Delete(ctx, sa.Role.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Delete(ctx, sa.DashboardRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Delete(ctx, sa.RoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.ClusterRoleBindings().Delete(ctx, sa.ClusterRoleBinding.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.ClusterRoles().Delete(ctx, sa.ClusterRole.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Delete(ctx, sa.ServiceAccount.Name, opts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (sa *ServiceAccountSuite) Check(ctx context.Context, c *client.Client, opts metav1.GetOptions) error {
	errs := []any{}

	if _, err := c.ServiceAccounts(client.ClientParams{Namespace: sa.ServiceAccount.Namespace}).Get(ctx, sa.ServiceAccount.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("model catalog service account '%s' already exists and is not ready for install", sa.ServiceAccount.Name))
	}

	if _, err := c.ClusterRoles().Get(ctx, sa.ClusterRole.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("model catalog service account cluster role '%s' already exists and is not ready for install", sa.ClusterRole.Name))
	}

	if _, err := c.ClusterRoleBindings().Get(ctx, sa.ClusterRoleBinding.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("model catalog service account cluster role binding '%s' already exists and is not ready for install", sa.ClusterRoleBinding))
	}

	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.RoleBinding.Namespace}).Get(ctx, sa.RoleBinding.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("model catalog service account role binding '%s' already exists and is not ready for install", sa.RoleBinding.Name))
	}

	if _, err := c.RoleBindings(client.ClientParams{Namespace: sa.DashboardRoleBinding.Namespace}).Get(ctx, sa.DashboardRoleBinding.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("RHOAI dashboard role binding '%s' already exists and is not ready for install", sa.DashboardRoleBinding.Name))
	}

	if _, err := c.Roles(client.ClientParams{Namespace: sa.Role.Namespace}).Get(ctx, sa.Role.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err == nil {
		errs = append(errs, fmt.Errorf("model catalog service account role '%s' already exists and is not ready for install", sa.Role.Name))
	}

	if _, err := c.Secrets(client.ClientParams{Namespace: sa.Secret.Namespace}).Get(ctx, sa.Secret.Name, opts); !errors.IsNotFound(err) {
		errs = append(errs, err)
	} else if err != nil {
		errs = append(errs, fmt.Errorf("model catalog service account token secret '%s' already exists and is not ready for install", sa.Secret.Name))
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

func GetServiceAccountTokenName(serviceAccountName string) string {
	return fmt.Sprintf("%s-token", serviceAccountName)
}

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
