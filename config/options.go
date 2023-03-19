package config

type Options struct {
	MapUsersAndGroups bool `mapstructure:"map_users_and_groups" yaml:"map_users_and_groups"`
	BindAll           bool `mapstructure:"bind_all" yaml:"bind_all"`
}
