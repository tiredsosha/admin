package configurator

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/tiredsosha/admin/tools/logger"
	"gopkg.in/yaml.v3"
)

// структура конфига
type conf struct {
	Broker   string `yaml:"broker"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"httpPort"`
	MqttOn   bool   `yaml:"mqttActive"`
}

// поиск конфига на диске
func getConf(file string, cnf any) error {
	yamlFile, err := os.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(yamlFile, cnf)
	}

	return err
}

// валидация конфига, валидация проходится если все заполнено
func validateConf(cfg *conf) error {
	var err error
	v := reflect.ValueOf(*cfg)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Name
		value := v.Field(i).Interface()
		if value == "" {
			err = errors.New("config: " + strings.ToLower(field) + " field is emtpy/nonexist")
			break
		}
		err = nil
	}
	return err
}

// создание дефолтного конфига и запись его в файл
func confFile() *conf {
	logger.Warn.Println("making a default config")
	confDef := conf{
		Broker:   "127.0.0.1",
		Username: "admin",
		Password: "password",
		MqttOn:   false,
		Port:     8080,
	}

	yamlData, _ := yaml.Marshal(confDef)

	if err := os.WriteFile("./configs/config.yaml", yamlData, 0644); err != nil {
		logger.Error.Fatal("can't to write default conf into the file")
	}
	return &confDef
}

func ConfInit() *conf {
	// создаем пустую версию конфига, или ссылку на него
	cfg := &conf{}
	// если у нас нет файла конфиг, то мы создаем дефотный конфиг
	if err := getConf("./configs/config.yaml", cfg); err != nil {
		logger.Error.Println(err)
		cfg = confFile()
	}
	// если у нас есть конфиг и он не проходит валидацию, то вы выходим из приложения
	if err := validateConf(cfg); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	return cfg
}

func ConfSubInit() {
	// если наше приложение не видит побочных конфигов
	if err := configRelay(); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
	if err := configPC(); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
	if err := configPJ(); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
}
