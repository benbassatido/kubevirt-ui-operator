package names

// Operand resource names (Helm/kubevirt-console-ui parity). Not configurable via CR.

const (
	// ResourceName is the required name of the cluster-scoped KubeVirtConsoleUI instance.
	ResourceName = "cluster"

	PriorityClass = "kubevirt-ui-cluster-critical"

	PluginDeployment    = "kubevirt-console-plugin"
	PluginService       = "kubevirt-console-plugin-service"
	PluginServingCert   = "plugin-serving-cert"
	PluginPort          = int32(9443)
	PluginNetworkPolicy = "kubevirt-console-plugin-np"

	ProxyDeployment    = "kubevirt-apiserver-proxy"
	ProxyService       = "kubevirt-apiserver-proxy-service"
	ProxyServingCert   = "console-proxy-serving-cert"
	ProxyPort          = int32(8080)
	ProxyNetworkPolicy = "kubevirt-apiserver-proxy-np"
	ProxyRoute         = "kubevirt-apiserver-proxy"

	ConsolePluginName        = "kubevirt-plugin"
	ConsolePluginDisplayName = "Kubevirt Console Plugin"
	ConsolePluginProxyAlias  = "kubevirt-apiserver-proxy"

	NginxConfigMap        = "nginx-conf"
	UserSettingsConfigMap = "kubevirt-user-settings"
	FeaturesConfigMap     = "kubevirt-ui-features"

	ConfigReaderRole        = "kubevirt-ui-config-reader"
	ConfigReaderRoleBinding = "kubevirt-ui-config-reader-rolebinding"

	PriorityClassDescription = "Priority class for critical KubeVirt and virtualization UI components"
)

const (
	PriorityClassValue int32 = 1000000000
)
