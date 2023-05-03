package dotenv

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"os"
	"strings"
)

type Reader struct{}

func NewReader(ctx context.Context, src, dst string) error {
	r := new(Reader)
	lines, err := r.GetLines(ctx, src)
	if err != nil {
		return err
	}
	data := r.Parse(ctx, lines)
	return r.Write(ctx, data, dst)
}

func (*Reader) GetLines(ctx context.Context, path string) ([]string, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("path", path)
	log.V(1).Info("reading file")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err, "failed to read file")
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func (*Reader) Parse(ctx context.Context, lines []string) string {
	log := logr.FromContextOrDiscard(ctx)
	data := strings.Builder{}
	data.WriteString("window._env_ = {")
	var envCount int
	var count int
	var total int
	for _, l := range lines {
		total++
		bits := strings.SplitN(l, "=", 2)
		if len(bits) != 2 {
			log.Info("failed to parse line as key=value", "line", l)
			continue
		}
		count++
		key := bits[0]
		value := bits[1]
		envValue := os.Getenv(key)
		if envValue == "" {
			log.V(1).Info("failed to locate value in environment - using fallback", "key", key)
			envValue = value
		} else {
			envCount++
		}
		data.WriteString(fmt.Sprintf("\n\t%s: \"%s\",", key, envValue))
	}
	data.WriteString("\n};")
	log.Info("loaded variables from the environment", "count", count, "total", total, "fromEnv", envCount)
	return data.String()
}

func (*Reader) Write(ctx context.Context, data, path string) error {
	log := logr.FromContextOrDiscard(ctx).WithValues("path", path)
	log.Info("writing file")
	// make sure we write with matching user and group permissions
	// for OpenShift compat
	err := os.WriteFile(path, []byte(data), 0660) //nolint:gosec
	if err != nil {
		log.Error(err, "failed to write to file")
		return err
	}
	return nil
}
