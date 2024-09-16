package main

import (
	"context"
	"errors"
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
	"os"
	"path/filepath"
)

type environment struct {
	LogLevel int `split_words:"true"`
	// DataPath defines the location of your applications
	// static files (required)
	DataPath string `split_words:"true" required:"true"`
	Env      struct {
		// File is the name of the JS file created
		// from the DotEnv file
		File string `split_words:"true" default:"env-config.js"`
		// Dir is the subdirectory to place the File in.
		// Used to mount an 'emptyDir' volume so that the
		// container can set 'readOnlyRootFilesystem'
		Dir string `split_words:"true"`
	}
	// Port defines the HTTP port to run on
	Port int `envconfig:"PORT" default:"8080"`
}

func main() {
	// read environment
	var e environment
	envconfig.MustProcess("nib", &e)

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(e.LogLevel * -1))
	log, ctx := logging.NewZap(context.Background(), zc)

	// static dir must be the first set value
	// so that we serve what was set at build
	// time
	staticDir := env.GetFirst("", func(key string) string {
		return filepath.Clean(e.DataPath)
	})

	var hasDotEnv bool
	dotPath := filepath.Join(e.DataPath, ".env")
	if _, err := os.Stat(dotPath); !errors.Is(err, os.ErrNotExist) {
		log.Info("detected .env file")
		hasDotEnv = true
	}

	if hasDotEnv {
		// dotEnv configuration must be the last set value
		// so that we allow the user to configure it
		// at runtime
		envFile := env.GetLast("", func(key string) string {
			return e.Env.File
		})
		log.Info("parsing dotenv", "file", dotPath)
		envPath := filepath.Join(staticDir, e.Env.Dir, envFile)
		if e.Env.Dir == "" {
			envPath = filepath.Join(staticDir, envFile)
		}
		if err := dotenv.NewReader(ctx, dotPath, envPath); err != nil {
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
