package types

type Config struct {
	Hedera Hedera `yaml:"hedera"`
}

type Hedera struct {
	Mirror Mirror `yaml:"mirror"`
}

type Mirror struct {
	Rosetta Rosetta `yaml:"rosetta"`
}

type Rosetta struct {
	Network string `yaml:"network"`

	Db DB `yaml:"db"`

	Port string `yaml:"port"`

	Shard string `yaml:"shard"`
}

type DB struct {
	Host string `yaml:"host"`

	Name string `yaml:"name"`

	Password string `yaml:"password"`

	Port string `yaml:"port"`

	Username string `yaml:"username"`
}
