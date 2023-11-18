package configuration

const (
	EXPIRATION_TOKEN       = 24 //HOURS
	NAME_SESSION           = "Authorization"
	ERROR_SERVICE_SECURITY = "Service error in Security"
	ERROR_PRIVILEGES_GEN   = "No Tienes los privilegios para acceder"
	ERROR_SERVICE_ACCOUNT  = "Service error in account"
	ERROR_DATABASE_GET     = "No se ha encontrado ningun motor de basededatos para procesar"
)

// Captured WEB
const (
	TIBIA_NEWS       = "https://www.tibia.com/news/?subtopic=latestnews"
	TIBIA_CHARS      = "https://www.tibia.com/community/?name=%s"
	TIBIA_WORDS      = "https://www.tibia.com/community/?subtopic=worlds"
	TIBIA_HIGHSCORES = "https://www.tibia.com/community/?subtopic=highscores"
	TIBIA_GUILD      = "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=%s"
)

// PRocesing ERROR WEBSCRAPING
const (
	COLLECTOR_EMPTY         = "NO TIENE CONTENIDO"
	COLLECTOR_NOT_COMPLETED = "No tiene informacion suficiente"
)

const (
	IMAGEN_PATH = "data/todays"
)
