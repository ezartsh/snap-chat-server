package cmd

import (
	// "database/sql"

	"net/http"
	"snap_chat_server/config"
	"snap_chat_server/database"
	"snap_chat_server/logger"
	"snap_chat_server/rest"
	"snap_chat_server/websockets"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webSocketCmd)
}

var webSocketCmd = &cobra.Command{
	Use:              "serve",
	TraverseChildren: true,
	Short:            "Start api and websocket server",
	Long:             `Start api and websocket server.`,
	Run: func(cmd *cobra.Command, args []string) {

		db := database.OpenDbConnection()

		hub := websockets.NewHub()

		go hub.Run()

		rest.HealthCheck()
		rest.AccountRegister(db)
		rest.AccountLogin(db)
		rest.GetContact(db)
		rest.AddContact(db)
		rest.RemoveContact(db)
		rest.GetGroup(db)
		rest.CreateGroup(db)
		rest.JoinGroup(db)
		rest.LeaveGroup(db)

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			websockets.ServeWs(hub, db, w, r)
		})

		logger.AppLog.Debugf("Listening and serve on port :%s ", config.Env.AppPort)
		err := http.ListenAndServe(":"+config.Env.AppPort, nil)

		if err != nil {
			logger.AppLog.Fatalf(err, "Listening and serve on port :%s", config.Env.AppPort)
		}
	},
}
