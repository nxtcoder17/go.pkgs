package log

import (
	"os"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	log zerolog.Logger
}

func New(opts Options) Logger {
	exclude := make([]string, 2)

	if !opts.ShowTimestamp {
		exclude = append(exclude, zerolog.TimestampFieldName)
	}

	if !opts.ShowLogLevel {
		exclude = append(exclude, zerolog.LevelFieldName)
	}

	if opts.Writer == nil {
		opts.Writer = os.Stderr
	}

	output := zerolog.ConsoleWriter{Out: opts.Writer, PartsExclude: exclude}

	// output.FormatLevel = func(level any) string {
	// 	return strings.ToUpper(fmt.Sprintf("%-6s", level))
	// }

	// output.FormatCaller = func(caller any) string {
	// 	return fmt.Sprintf("| %s |", caller)
	// }

	zerolog.ErrorFieldName = "err"
	zerolog.FormattedLevels = map[zerolog.Level]string{
		zerolog.TraceLevel: "TRACE",
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO ",
		zerolog.WarnLevel:  "WARN ",
		zerolog.ErrorLevel: "ERROR",
		zerolog.FatalLevel: "FATAL",
	}

	// output.FormatMessage = func(i interface{}) string {
	// 	return fmt.Sprintf(" %s", i)
	// }
	//
	level := zerolog.InfoLevel
	if opts.ShowDebugLogs {
		level = zerolog.DebugLevel
	}

	zctx := zerolog.New(output).With()

	if opts.ShowCaller {
		zctx = zctx.CallerWithSkipFrameCount(3)
	}

	return &ZeroLogger{log: zctx.Logger().Level(level)}
}

func (zl *ZeroLogger) Debug(msg string, kv ...any) {
	zl.log.Debug().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Info(msg string, kv ...any) {
	zl.log.Info().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Warn(msg string, kv ...any) {
	zl.log.Warn().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Error(err error, msg string, kv ...any) {
	zl.log.Error().Fields(kv).Err(err).Msg(msg)
}

func (zl *ZeroLogger) Fatal(msg string, kv ...any) {
	zl.log.Fatal().Fields(kv).Msg(msg)
}

// With implements Logger.
func (zl *ZeroLogger) With(kv ...any) Logger {
	log := zl.log.With().Fields(kv).Logger()
	return &ZeroLogger{log: log}
}
