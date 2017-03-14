package config

import (
    "strconv"
    "fmt"
    "path/filepath"
    "strings"
    
    
    "github.com/vaughan0/go-ini"

    "github.com/SCKelemen/Cassius/common"
)

func LoadConfigFromFile(path string) (common.AppConfig, error) {
    conf, err := OpenConfigFile(path)
	if err != nil {
        return common.AppConfig{}, err
	}

    appConfig := common.AppConfig{}

    address, _ := conf.Get("server", "address")
    port, _ := conf.Get("server", "port")
    appConfig.ListenLocation = address + ":" + port

    appConfig.DatabaseHost, _ = conf.Get("database", "host")
    appConfig.DatabaseName, _ = conf.Get("database", "name")
    appConfig.DatabaseUser, _ = conf.Get("database", "user")
    appConfig.DatabasePassword, _ = conf.Get("database", "password")
    
    dbPortStr, _ := conf.Get("database", "port")
    dbPort, _ := strconv.Atoi(dbPortStr)
    appConfig.DatabasePort = uint16(dbPort)

    maxConnStr, _ := conf.Get("database", "max_conn")
    maxConn, _ := strconv.Atoi(maxConnStr)
    appConfig.DatabaseMaxConnections = maxConn

    appConfig.SmtpServer, appConfig.SmtpActive = conf.Get("smtp", "server")
    if appConfig.SmtpActive {
        appConfig.SmtpRootUrl, _ = conf.Get("smtp", "root_url") 
        appConfig.SmtpFromAddress, _ = conf.Get("smtp", "from_address") 
        appConfig.SmtpUser, _ = conf.Get("smtp", "username") 
        appConfig.SmtpPassword, _ = conf.Get("smtp", "password") 
        appConfig.SmtpPort, _ = conf.Get("smtp", "port")
    }

    hashalg, _ := conf.Get("security", "hash_algorithm")
    switch strings.ToLower(hashalg) {
        case "bcrypt":
        case "b-crypt":
            appConfig.HashAlgorithm = "bcrypt"
        default:
            appConfig.HashAlgorithm = "scrypt"
    }
  
    return appConfig, nil
}

func OpenConfigFile(path string) (ini.File, error) {
    path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("Invalid config path: %v", err)
	}

	file, err := ini.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config file: %v", err)
	}

	return file, nil
}
