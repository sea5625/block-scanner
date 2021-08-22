package main

import (
	"log"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/handlers/alerting"
	"motherbear/backend/handlers/auth"
	"motherbear/backend/handlers/blocks"
	"motherbear/backend/handlers/channels"
	"motherbear/backend/handlers/nodes"
	"motherbear/backend/handlers/nodetype"
	"motherbear/backend/handlers/resources"
	"motherbear/backend/handlers/settings"
	"motherbear/backend/handlers/symptom"
	"motherbear/backend/handlers/txs"
	"motherbear/backend/handlers/users"
	"motherbear/backend/polarbear"
	"motherbear/backend/prometheus/prom_crawler"
	"motherbear/backend/utility"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"net/http"
	"os"
	"path"

	. "motherbear/backend/configuration"
	"motherbear/backend/logger"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "motherbear/docs"
)

var frontPageList = []string{"block", "blockInfo", "txInfo", "users", "nodes", "settingEtc", "settingNodes", "settingChannels", "blockList", "txList"}

func init() {
	// Initialize yaml settings
	if _, err := os.Stat(constants.ConfigFolderName); os.IsNotExist(err) {
		os.MkdirAll(constants.ConfigFolderName, os.ModePerm)
	}

	// Read the configuration file.
	InitConfigData(constants.ConfigFolderName + "/" + constants.ConfigFileName)

	// Set Log Level (1=Info, 2=Debug)
	logger.SetLogLevel(Conf().ETC.LogLevel)

	// Initialize yaml settings
	if _, err := os.Stat(constants.DBFolderName); os.IsNotExist(err) {
		os.MkdirAll(constants.DBFolderName, os.ModePerm)
	}

	// Init ISSAC database
	if Conf().ETC.DB[0].DBType == constants.DBTypeMysql {
		// mysql param is (dbType, id, pass, database)
		db.InitDB(constants.DBTypeMysql, Conf().ETC.DB[0].Id, Conf().ETC.DB[0].Pass, Conf().ETC.DB[0].Database)
	} else if Conf().ETC.DB[0].DBType == constants.DBTypeSqlite3 {
		// sqlite param is (dbType, path)
		db.InitDB(constants.DBTypeSqlite3, Conf().ETC.DB[0].DBPath)
	}

	// Init Block Crawling database
	if Conf().Blockchain.DB[0].DBType == constants.DBTypeMysql {
		// mysql param is (dbType, id, pass, database)
		if err := polarbear.Init(constants.DBTypeMysql, Conf().Blockchain.DB[0].Id,
			Conf().Blockchain.DB[0].Pass, Conf().Blockchain.DB[0].Database); os.IsExist(err) {
			panic(err)
		}
	} else if Conf().Blockchain.DB[0].DBType == constants.DBTypeSqlite3 {
		// sqlite param is (dbType, path)
		if err := polarbear.Init(constants.DBTypeSqlite3, Conf().Blockchain.DB[0].DBPath); os.IsExist(err) {
			panic(err)
		}
	}

	// Run crawling data from prometheus.
	_, err := prom_crawler.BeginToCrawl()
	if err != nil {
		panic(err)
	}

	// Run crawling block data from loopchain.
	polarbear.BeginToCrawl()
}

// @title ISAAC
// @version RC0.2.14a
// @description The management system of loopchain.
// @termsOfService https://www.iconloop.com/

// @contact.name API Support
// @contact.url https://www.iconloop.com/
// @contact.email haspori@icon.foundation

// @host localhost:6553
// @BasePath /api/v1
func main() {

	// Use 4 CPU cores as default.
	runtime.GOMAXPROCS(4)

	// Set the router as the default one shipped with Gin.
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Serve the frontend file.
	dir, err := os.Getwd()
	if err != nil {
		logger.Fatalln(err)
	}
	var staticFilesPath = path.Join(dir, "frontend", "build")

	// Serve static file.
	logger.Info(staticFilesPath)
	router.Use(static.Serve("/", static.LocalFile(staticFilesPath, true)))

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		pathSplit := strings.Split(path, constants.URLPathSeparator)

		if len(pathSplit) <= 1 {
			c.AbortWithStatus(http.StatusNotFound)
		}

		// If url is front-end page, go to current page.
		//  - If go to index.html go to current page when use React route.
		// Otherwise, return to 404 error.
		if utility.IsExistValueInList(pathSplit[1], frontPageList) {
			c.File(staticFilesPath + "/index.html")
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	// Move swagger document under /api/v1.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := router.Group(constants.APIVersionURL)

	{
		apiV1.Use(auth.AuthentificateMiddleware())
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// /api/v1/users
		apiV1.GET(constants.UsersGetListAPIURL, users.GetHandler)
		apiV1.GET(constants.UsersGetAPIURL, users.GetHandler)
		apiV1.POST(constants.UsersPostAPIURL, users.PostHandler)
		apiV1.PUT(constants.UsersPutAPIURL, users.PutHandler)
		apiV1.DELETE(constants.UsersDeleteAPIURL, users.DeleteHandler)

		// /api/v1/channels
		apiV1.GET(constants.ChannelsGetListAPIURL, channels.GetHandler)
		apiV1.GET(constants.ChannelsGetAPIURL, channels.GetHandler)
		apiV1.POST(constants.ChannelsPostAPIURL, channels.PostHandler)
		apiV1.PUT(constants.ChannelsPutAPIURL, channels.PutHandler)
		apiV1.DELETE(constants.ChannelsDeleteAPIURL, channels.DeleteHandler)

		// /api/v1/nodes
		apiV1.GET(constants.NodesGetListAPIURL, nodes.GetHandler)
		apiV1.GET(constants.NodesGetAPIURL, nodes.GetHandler)
		apiV1.POST(constants.NodesPostAPIURL, nodes.PostHandler)
		apiV1.PUT(constants.NodesPutAPIURL, nodes.PutHandler)
		apiV1.DELETE(constants.NodesDeleteAPIURL, nodes.DeleteHandler)

		// /api/v1/alerting
		apiV1.GET(constants.AlertingGetListAPIURL, alerting.GetHandler)
		apiV1.PUT(constants.AlertingPutAPIURL, alerting.PutHandler)

		// /api/v1/settings
		apiV1.GET(constants.SettingsGetListAPIURL, settings.GetHandler)
		apiV1.PUT(constants.SettingsPutAPIURL, settings.PutHandler)

		// /api/v1/auth
		apiV1.POST(constants.AuthLoginAPIURL, auth.Login)
		apiV1.PUT(constants.AuthReissueTokenAPIURL, auth.Token)

		// /api/v1/block
		apiV1.GET(constants.BlockGETListAPIURL, blocks.GetHandlerList)
		apiV1.GET(constants.BlockGETAPIURL, blocks.GetHandler)

		// /api/v1/txs
		apiV1.GET(constants.TxGETListAPIURL, txs.GetHandlerList)
		apiV1.GET(constants.TxGETAPIURL, txs.GetHandler)

		// /api/v1/resources
		apiV1.GET(constants.ResourcesGETAPIURL, resources.GetHandler)

		// /api/v1/symptom
		apiV1.GET(constants.SymptomGETAPIURL, symptom.GetHandlerList)

		// /api/v1/prometheus
		apiV1.GET(constants.PrometheusGETAPIURL, nodetype.GetHandler)
	}

	////////////////////
	// Signal handling
	go handleSignals()

	////////////////////
	// Start the app
	logger.Info("ISAAC backend running!")
	router.Run(constants.ServerPort)
}

func Quit() {
	// Delete ISAAC db instance
	db.CloseDBInstance()

	// Delete Block Crawling db instance
	polarbear.CloseDBInstance()

	// Stop crawling data from prometheus.
	prom_crawler.StopToCrawl()

	// Stop crawling block data from loopchain.
	polarbear.StopToCrawl()

	os.Exit(0)
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println(sig)
			logger.Infof("ISAAC is exited, signal : %s", sig)
			Quit()
			return
		}
	}
}

// CORSMiddleware handles CORS issue from web front-end.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Pragma")
		// To-Do: Should allow for crendential with JWT.
		//c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", c.ClientIP()+":6553")
		c.Header("Access-Control-Allow-Origin", c.ClientIP()+":3000")
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Expose-Headers", "Content-Range")
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
