package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/SawitProRecruitment/UserService/core/service/usersvc"
	sawithttp "github.com/SawitProRecruitment/UserService/handler/http"
	"github.com/SawitProRecruitment/UserService/lib/locker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := initConfig()
		ctx := context.Background()
		logger := initLogger()

		repo, err := initPostgres(ctx, cfg.DB)
		if err != nil {
			log.Fatalf("error init postgres repo: %v", err)
		}

		libLocker := locker.New(cfg.AES.SecretKey)
		authSvc, err := initAuthSvc(cfg.Auth, repo)
		if err != nil {
			log.Fatalf("error init auth service: %v", err)
		}

		userSvc := usersvc.New(repo)

		// HTTP handler based on api.yml
		handler := sawithttp.NewHandler(logger, libLocker, userSvc, authSvc)
		s := initServer(cfg.Server, handler)

		// running http server
		lock := make(chan error)
		go func(lock chan error) {
			lock <- s.ListenAndServe()
		}(lock)

		logger.Info(fmt.Sprintf("running at %s", s.Addr))
		err = <-lock
		if err != nil {
			_ = s.Close()
			log.Fatal(err)
		}
	},
}
