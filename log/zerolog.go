package log

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type ZeroLogger struct {
	logger zerolog.Logger
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

	return &ZeroLogger{logger: zctx.Logger().Level(level)}
}

func (zl *ZeroLogger) Debug(msg string, kv ...any) {
	zl.logger.Debug().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Info(msg string, kv ...any) {
	zl.logger.Info().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Warn(msg string, kv ...any) {
	zl.logger.Warn().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) Error(err error, msg string, kv ...any) {
	zl.logger.Error().Fields(kv).Err(err).Msg(msg)
}

func (zl *ZeroLogger) Fatal(msg string, kv ...any) {
	zl.logger.Fatal().Fields(kv).Msg(msg)
}

func (zl *ZeroLogger) With(kv ...any) Logger {
	log := zl.logger.With().Fields(kv).Logger()
	return &ZeroLogger{logger: log}
}

func (zl *ZeroLogger) Slog() *slog.Logger {
	l := zl.logger.With().CallerWithSkipFrameCount(2).Logger()
	return slog.New(slogzerolog.Option{Logger: &l}.NewZerologHandler())
}
