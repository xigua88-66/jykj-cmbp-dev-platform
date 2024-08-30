package config

type CMBPBase struct {
	JupyterLabUrl       string `mapstructure:"jupyter-lab-url" json:"jupyter-lab-url" yaml:"jupyter-lab-url"`
	CmbpUrl             string `mapstructure:"cmbp-url" json:"cmbp-url" yaml:"cmbp-url"`
	OssPath             string `mapstructure:"oss-path" json:"oss-path" yaml:"oss-path"`
	ModelMarketMedia    string `mapstructure:"model-market-media" json:"model-market-media" yaml:"model-market-media"`
	ModelWareHouseMedia string `mapstructure:"model-warehouse-media" json:"model-warehouse-media" yaml:"model-warehouse-media"`
	OssModelPath        string `mapstructure:"oss-model-path" json:"oss-model-path" yaml:"oss-model-path"`
	DockerRegistry      string `mapstructure:"docker-registry" json:"docker-registry" yaml:"docker-registry"`
}
