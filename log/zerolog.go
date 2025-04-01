package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	slogcommon "github.com/samber/slog-common"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type ZeroLogger struct {
	logger     *zerolog.Logger
	skipFrames int
}

func newZeroLogger(options ...Options) *ZeroLogger {
	exclude := make([]string, 2)

	opts := defaultOptions

	if len(options) >= 1 {
		opts = options[0]
	}

	if !opts.ShowTimestamp {
		exclude = append(exclude, zerolog.TimestampFieldName)
	}

	if !opts.ShowLogLevel {
		exclude = append(exclude, zerolog.LevelFieldName)
	}

	if opts.Writer == nil {
		opts.Writer = os.Stderr
	}

	writer := func() io.Writer {
		if !opts.JSONFormat {
			return zerolog.ConsoleWriter{
				Out:          opts.Writer,
				PartsExclude: exclude,
				TimeFormat:   time.RFC822,
				FormatCaller: func(caller any) string {
					switch c := caller.(type) {
					case string:
						return fmt.Sprintf("| %s |", trimPath(c, 2))
					default:
						return "« invalid caller type » " + fmt.Sprintf("%v", c)
					}
				},
			}
		}

		return os.Stderr
	}()

	zerolog.TimeFieldFormat = time.RFC822
	zerolog.ErrorFieldName = "err"
	zerolog.FormattedLevels = map[zerolog.Level]string{
		zerolog.TraceLevel: "TRACE",
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO ",
		zerolog.WarnLevel:  "WARN ",
		zerolog.ErrorLevel: "ERROR",
		zerolog.FatalLevel: "FATAL",
	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", trimPath(file, 2), line)
	}

	level := zerolog.InfoLevel
	if opts.ShowDebugLogs {
		level = zerolog.DebugLevel
	}

	zctx := zerolog.New(writer).With().Timestamp()

	skipFrames := 2 + 1 // 2 because of zerolog and 1 because i am creating a helper function

	if opts.ShowCaller {
		zctx = zctx.CallerWithSkipFrameCount(skipFrames)
	}

	logger := zctx.Logger().Level(level)
	return &ZeroLogger{
		logger:     &logger,
		skipFrames: skipFrames,
	}
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

func (zl *ZeroLogger) SkipFrames(skip int) Logger {
	l := zl.logger.With().CallerWithSkipFrameCount(zl.skipFrames + skip).Logger()
	return &ZeroLogger{logger: &l, skipFrames: zl.skipFrames + skip}
}

func (zl *ZeroLogger) With(kv ...any) Logger {
	log := zl.logger.With().Fields(kv).Logger()
	return &ZeroLogger{logger: &log}
}

func (zl *ZeroLogger) Slog() *slog.Logger {
	// l := zl.logger.With().CallerWithSkipFrameCount(3).Logger()
	l := zl.logger.With().CallerWithSkipFrameCount(zl.skipFrames + 2).Logger()
	slogzerolog.ErrorKeys = []string{"error", "err"}
	return slog.New(slogzerolog.Option{
		Logger: &l,

		NoTimestamp: true,

		// [source](https://github.com/samber/slog-zerolog/blob/11dde13940f914ffbdb501c69b3765069914abf1/converter.go#L14-L30)
		Converter: func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
			// WHY? just makes things more like original logger

			// aggregate all attributes
			attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

			// handler formatter
			output := slogcommon.AttrsToMap(attrs...)

			return output
		},
	}.NewZerologHandler())
}
