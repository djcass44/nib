package main

import (
	"context"
	"github.com/djcass44/go-utils/logging"
	"github.com/djcass44/nib/srv/internal/env"
	"github.com/djcass44/nib/srv/pkg/dotenv"
	"github.com/djcass44/nib/srv/pkg/spa"
	"github.com/go-http-utils/etag"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/autokubeops/serverless"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"path/filepath"
)

type environment struct {
	LogLevel int `split_words:"true"`
	// StaticDir defines the location of your applications
	// static files (required)
	StaticDir string `split_words:"true" required:"true"`
	// DotEnv is the path to the .env file that should
	// be converted into an env-config.js file (optional)
	DotEnv string `split_words:"true"`
	// EnvFile is the name of the JS file created
	// from the DotEnv file
	EnvFile string `split_words:"true" default:"env-config.js"`
	// Port defines the HTTP port to run on
	Port int `envconfig:"PORT" default:"8080"`
}

func main() {
	// read environment
	var e environment
	envconfig.MustProcess("srv", &e)

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(e.LogLevel * -1))
	log, ctx := logging.NewZap(context.TODO(), zc)

	// static dir must be the first set value
	// so that we serve what was set at build
	// time
	staticDir := env.GetFirst("", func(key string) string {
		return e.StaticDir
	})

	if e.DotEnv != "" {
		// dotEnv configuration must be the last set value
		// so that we allow the user to configure it
		// at runtime
		dotEnv := env.GetLast("", func(key string) string {
			return e.DotEnv
		})
		envFile := env.GetLast("", func(key string) string {
			return e.EnvFile
		})
		log.Info("parsing dotenv", "file", dotEnv)
		if err := dotenv.NewReader(ctx, dotEnv, filepath.Join(staticDir, envFile)); err != nil {
			log.Error(err, "failed to configure - this may cause undefined behaviour", "file", envFile)
		}
	}

	// start the file server
	log.Info("serving directory", "path", staticDir)

	// configure routing
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(spa.NewFileSystem(http.Dir(staticDir))))

	// start the http server
	serverless.NewBuilder(handlers.CompressHandler(etag.Handler(router, false))).
		WithPort(e.Port).
		WithLogger(log).
		Run()
}
