package config

type Values struct {
	Log              LogConfig
	Server           ServerConfig
	Environment      string
	ProfilingEnabled bool
	Datastores       Datastores
}

type LogConfig struct {
	Level string
}

type ServerConfig struct {
	Host string
	Port int
}

type Datastores struct {
	TestDB MongoDB `mapstructure:"testDB"`
}

type MongoDB struct {
	Hosts      string       `mapstructure:"hosts"`
	Port       int          `mapstructure:"port"`
	User       string       `mapstructure:"user"`
	Password   string       `mapstructure:"password"`
	Database   string       `mapstructure:"database"`
	AuthSource string       `mapstructure:"authSource"`
	ReplicaSet string       `mapstructure:"replicaSet"`
	AppName    string       `mapstructure:"appName"`
	Options    MongoOptions `mapstructure:"options"`
}

type MongoOptions struct {
	MaxPoolSize       int `mapstructure:"maxPoolSize"`
	MinPoolSize       int `mapstructure:"minPoolSize"`
	IdleTimeout       int `mapstructure:"idleTimeout"`
	ConnectionTimeout int `mapstructure:"connectionTimeout"`
}
