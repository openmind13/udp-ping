package aim

type AimConfig struct {
	Name       string `mapstructure:"name"`
	LocalAddr  string `mapstructure:"local_addr"`
	RemoteAddr string `mapstructure:"remote_addr"`
}

func (ac AimConfig) IsValid() bool {
	if ac.Name == "" || ac.LocalAddr == "" || ac.RemoteAddr == "" {
		return false
	}
	return true
}
