package plugins

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUnmarshalConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    corev1.ConfigMap
		expected *DynamicPluginConfig
	}{
		{
			name: "Test Unmarsaling YAML into DynamicPluginConfig",
			input: corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-dynamic-plugins",
				},
				Data: map[string]string{
					DefaultDynamicPluginsKey: `includes:
  - test.yaml
plugins:
  - disabled: false
    package: testpackage
`,
				},
			},
			expected: &DynamicPluginConfig{
				Includes: []string{"test.yaml"},
				Plugins: DynamicPluginList{
					{
						Disabled: false,
						Package:  "testpackage",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			result, err := UnmarshalConfig(&test.input)

			if err != nil {
				tt.Fatal(err)
			} else if !reflect.DeepEqual(result, test.expected) {
				tt.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
