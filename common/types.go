type AppConfig struct {
    ListenLocation string

    DatabaseHost string    
    DatabasePort uint16
    DatabaseName string
    DatabaseUser string
    DatabasePassword string
    DatabaseMaxConnections int

    SmtpActive bool
    SmtpServer string
    SmtpPort string
    SmtpFromAddress string
    SmtpUser string
    SmtpPassword string
    SmtpRootUrl string

    HashAlgorithm string
}