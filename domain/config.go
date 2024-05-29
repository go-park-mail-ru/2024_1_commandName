package domain

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	App struct {
		IsDebug    bool   `yaml:"isDebug"`
		InMemory   bool   `yaml:"inMemory"`
		AvatarPath string `yaml:"avatarPath"`
	} `yaml:"app"`
	Yandex struct {
		TranslateKey string `yaml:"trKey"`
		Url          string `yaml:"url"`
		FolderID     string `yaml:"folderId"`
		Header       string `yaml:"header"`
		Method       string `yaml:"method"`
	} `yaml:"yandex"`
	Gpt struct {
		TrKey    string `yaml:"trKey"`
		Url      string `yaml:"url"`
		FolderID string `yaml:"folderId"`
		Header   string `yaml:"header"`
		Method   string `yaml:"method"`
	}
}

type YandexConfig struct {
	TranslateKey string
	FolderID     string
	Url          string
	Method       string
	Header       string
}
