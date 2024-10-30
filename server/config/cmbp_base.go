package config

type CMBPBase struct {
	JupyterLabUrl       string `mapstructure:"jupyter-lab-url" json:"jupyter-lab-url" yaml:"jupyter-lab-url"`
	CmbpUrl             string `mapstructure:"cmbp-url" json:"cmbp-url" yaml:"cmbp-url"`
	OssPath             string `mapstructure:"oss-path" json:"oss-path" yaml:"oss-path"`
	ModelMarketMedia    string `mapstructure:"model-market-media" json:"model-market-media" yaml:"model-market-media"`
	ModelWareHouseMedia string `mapstructure:"model-warehouse-media" json:"model-warehouse-media" yaml:"model-warehouse-media"`
	OssModelPath        string `mapstructure:"oss-model-path" json:"oss-model-path" yaml:"oss-model-path"`
	DockerRegistry      string `mapstructure:"docker-registry" json:"docker-registry" yaml:"docker-registry"`
	ModelWareHouse      string `mapstructure:"model-warehouse" json:"model-warehouse" yaml:"model-warehouse"`
	ModelPath           string `mapstructure:"model-path" json:"model-path" yaml:"model-path"`
	OssModelMedia       string `mapstructure:"oss-model-media" json:"oss-model-media" yaml:"oss-model-media"`
	OssMarketModelPath  string `mapstructure:"oss-market-model-path" json:"oss-market-model-path" yaml:"oss-market-model-path"`
	WorkFlowUrl         string `mapstructure:"workflow-url" json:"workflow-url" yaml:"workflow-url"`
	WorkFlowAppName     string `mapstructure:"workflow-app-name" json:"workflow-app-name" yaml:"workflow-app-name"`
	WorkFlowAppSk       string `mapstructure:"workflow-app-sk" json:"workflow-app-sk" yaml:"workflow-app-sk"`
	OssRuntimeLibrary   string `mapstructure:"oss-runtime-library" json:"oss-runtime-library" yaml:"oss-runtime-library"`
	RuntimeDownRole     string `mapstructure:"runtime-down-role" json:"runtime-down-role" yaml:"runtime-down-role"`
	OssMode             string `mapstructure:"oss-mode" json:"oss-mode" yaml:"oss-mode"`
	OssExpireSeconds    int    `mapstructure:"oss-expire-seconds" json:"oss-expire-seconds" yaml:"oss-expire-seconds"`
}

type CMBPModelCfg struct {
	CythonPath             string `mapstructure:"cython-path" json:"cython-path" yaml:"cython-path"`
	CythonSrcLibInclude    string `mapstructure:"cython-src-lib-include" json:"cython-src-lib-include" yaml:"cython-src-lib-include"`
	CythonSrcLibIncludeArm string `mapstructure:"cython-src-lib-include-arm" json:"cython-src-lib-include-arm" yaml:"cython-src-lib-include-arm"`
}
