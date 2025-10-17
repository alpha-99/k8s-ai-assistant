package models

type APIKey struct {
	// 用于在访问外部API服务时进行认证的令牌名称
	Name string `yaml:"name"`

	// 用于在访问外部API服务时进行认证的令牌的值
	Value string `yaml:"value"`

	// 在访问外部API服务时进行认证的令牌是放在header中还是query中，如果api没令牌，填none
	In string `yaml:"in"`
}

type APIProvider struct {
	APIKey APIKey `yaml:"apiKey"`
}

type APIConfig struct {
	APIProvider APIProvider `yaml:"apiProvider"`
	// 工具的openapi文档
	API string `yaml:"api"`
}

type Config struct {
	Apis              APIConfig `yaml:"apis"`
	Instruction       string    `yaml:"instruction"`
	MaxIterationSteps int       `yaml:"max_iteration_steps"`
}
