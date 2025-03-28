package configuration

const (
	DEFAULT_LENGUAJE = "es"
)

const (
	MAX_LEN_PASSWORD       = 5
	ACCES_ADMIN            = 3
	MAX_FILE_SIZE          = 5 << 20 //3 mgbyttes
	EXPIRATION_TOKEN       = 24      //HOURS
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
	TIBIA_MONSTER    = "https://www.tibia.com/library/?subtopic=creatures"
	TIBIA_BOSSES     = "https://www.tibia.com/library/?subtopic=boostablebosses"
	TWITCH_CLIPS     = "https://clips.twitch.tv/embed?clip=%s&parent=%s&muted=true"
	TWITCH_CANAL     = "https://player.twitch.tv/?channel=%s&parent=%s&muted=true"
)

// PRocesing ERROR WEBSCRAPING
const (
	COLLECTOR_EMPTY         = "NO TIENE CONTENIDO"
	COLLECTOR_NOT_COMPLETED = "No tiene informacion suficiente"
)

const (
	//Path controller
	IMAGEN_PATH = "data/image/todays"
)

const (
	PATH_CONFIG           = "config.yml"
	PATH_LENGUAJE_CLIENT  = "data/lenguaje.csv"
	PATH_LENGUAJE_SERVER  = "data/errorServer.csv"
	PATH_STATIC_CREATURES = "creatures/%s"
	PATH_STATIC_BOSSES    = "bosses/%s"
	PATH_STATIC_PUBLIC    = "data/image/"
	PATH_RECOVERY_EMAIL   = "data/email/recovery.html"
)
const (
	PATH_WEB_CREATURES_ID = "static/creatures_id.html"
	PATH_WEB_CREATURES    = "static/creatures.html"
	PATH_WEB_BOSSES       = "static/bosses.html"
	PATH_WEB_ERROR        = "static/error404.html"
	PATH_WEB_ADMIN        = "static/admin.html"
	PATH_WEB_INDEX        = "static/index.html"
	PATH_WEB_MY_PROFILE   = "static/my_profile.html"
	PATH_WEB_TODAYS       = "static/todays.html"
	PATH_WEB_TODAYS_POST  = "static/todays_post.html"
	PATH_WEB_UPLOAD_FILES = "static/upload_files.html"
	PATH_WEB_SHARED_LOOT  = "static/shared_loot.html"
	PATH_WEB_TOOLS        = "static/tools.html"
	PATH_WEB_SHARED_EXP   = "static/shared_exp.html"
	PATH_WEB_MY_FAVO_PIC  = "static/my_favorite_pictures.html"
	PATH_WEB_MY_PLAYERS   = "static/my_players.html"
	PATH_WEB_RECOVERY     = "static/recovery.html"
)

const (
	//router web
	ROUTER_INDEX               = "/"
	ROUTER_TODAYS_POST         = "/todays_post/%d"
	ROUTER_TODAYS_POST_APROVED = "/auth/todays_aproved/%d"
	ROUTER_MY_PLAYERS          = "/auth/my_players"
	ROUTER_UPLOAD_IMAGES       = "/auth/upload_image"
	ROUTER_RECOVERY_CODE       = "http://%s/recovery/%s"
	ROUTER_STREAMER_ID         = "/auth/streamer/%d"
	ROUTER_NEWS_TICKET_ID      = "/auth/newsticket/%d"
	ROUTER_NEWS_TICKET         = "/auth/newsticket"
	ROUTER_CREATURE_ID         = "/creatures/%d"
	ROUTER_BOSSES_ID           = "/bosses/%d"
)

const (
	//Var Multilenguaje
	ErrorDefault = iota + 1
	ErrorPassword
	ErrorEmptyField
	ErrorPrivileges
	ErrorInternal
	ErrorError
	ErrorCode
	ErrorPolicies
	NotAuthorized
	ErrorMaxFileSize
	ErrorExpireSession
	ProcessedSuccess
	SavedSuccess
	GenerateSucces
	ChangePassword
	ErrorExpiredProcess
	AccountRecovery
)
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)
