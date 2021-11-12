package testcommon

import (
	"fmt"
	"os"
	"strings"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/proto"
)

const CorrelationElementsFixedLength = 14
const EllipsisString = ".."
const FormatColoredString = "\033[%s%s\033[0m[%s] %s %s %s %s\n"
const FormatPlainString = "%s%s %s %s %s %s\n"

func SetConfigurableLoggerFormatter() {
	logger.ClearLogObservers()
	formatter := &ConfigurableFormatter{}
	formatter.SetUp()
	logger.AddLogObserver(os.Stdout, formatter)
}

type ConfigurableFormatter struct {
	Colored          bool
	ShowLevel        bool
	ShowDate         bool
	ShowTime         bool
	ShowLoggerName   bool
	MessageLength    int
	LoggerNameLength int
	BracketStart     string
	BracketEnd       string
	BracketsLength   int
	ShowBrackets     bool
	FormatString     string
	TimestampFormat  string
	// ShowLastLoggerNamePart bool
	// ShowMultilineArgs      bool
}

func (cf *ConfigurableFormatter) SetUp() {
	logger.ToggleLoggerName(cf.ShowLoggerName)
	if cf.MessageLength == 0 {
		cf.MessageLength = 40
	}
	if cf.LoggerNameLength == 0 {
		cf.LoggerNameLength = 20
	}
	if cf.TimestampFormat == "" {
		if cf.ShowDate {
			cf.TimestampFormat = "2006-01-02 15:04:05.000"
		} else {
			cf.TimestampFormat = cf.BracketStart + "15:04:05.000" + cf.BracketEnd
		}
	}
	if cf.FormatString == "" {
		if cf.Colored {
			cf.FormatString = FormatColoredString
		} else {
			cf.FormatString = FormatPlainString
		}
	}
	if !cf.ShowBrackets {
		cf.BracketsLength = 0
		cf.BracketStart = ""
		cf.BracketEnd = ""
	} else {
		if cf.BracketStart == "" && cf.BracketEnd == "" {
			cf.BracketStart = "["
			cf.BracketEnd = "]"
		}
		cf.BracketsLength = len(cf.BracketStart) + len(cf.BracketEnd)
	}
}

// Output converts the provided LogLineHandler into a slice of bytes ready for output
func (cf *ConfigurableFormatter) Output(line logger.LogLineHandler) []byte {
	if line == nil {
		return nil
	}

	var level string
	if cf.ShowLevel {
		level = fmt.Sprintf("%s", logger.LogLevel(line.GetLogLevel()))
	} else {
		level = ""
	}

	var timestamp string
	if cf.ShowTime {
		timestamp = cf.displayTime(line.GetTimestamp())
	} else {
		timestamp = ""
	}

	loggerName := ""
	correlation := ""
	message := cf.formatMessage(line.GetMessage())
	args := cf.formatArgsNoAnsi(line.GetArgs()...)

	if logger.IsEnabledLoggerName() {
		loggerName = cf.formatLoggerName(line.GetLoggerName())
	}

	if logger.IsEnabledCorrelation() {
		correlation = cf.formatCorrelationElements(line.GetCorrelation())
	}

	return []byte(
		fmt.Sprintf(cf.FormatString,
			level,
			timestamp, loggerName, correlation,
			message, args,
		),
	)
}

// formatArgsNoAnsi iterates through the provided arguments displaying the argument name and after that its value
// The arguments must be provided in the following format: "name1", "val1", "name2", "val2" ...
// It ignores odd number of arguments and it does not use ANSI colors
func (cf *ConfigurableFormatter) formatArgsNoAnsi(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	argString := ""
	for index := 1; index < len(args); index += 2 {
		argString += fmt.Sprintf("%s = %s ", args[index-1], args[index])
	}

	return argString
}

// IsInterfaceNil returns true if there is no value under the interface
func (cf *ConfigurableFormatter) IsInterfaceNil() bool {
	return cf == nil
}

func (cf *ConfigurableFormatter) displayTime(timestamp int64) string {
	t := time.Unix(0, timestamp)
	return t.Format(cf.TimestampFormat)
}

func (cf *ConfigurableFormatter) formatMessage(msg string) string {
	return cf.padRight(msg, cf.MessageLength)
}

func (cf *ConfigurableFormatter) padRight(str string, maxLength int) string {
	paddingLength := maxLength - len(str)

	if paddingLength > 0 {
		return str + strings.Repeat(" ", paddingLength)
	}

	return str
}

func (cf *ConfigurableFormatter) formatLoggerName(name string) string {
	name = cf.truncatePrefix(name, cf.LoggerNameLength-cf.BracketsLength)
	formattedName := fmt.Sprintf("%s%s%s", cf.BracketStart, name, cf.BracketEnd)

	return cf.padRight(formattedName, cf.LoggerNameLength)
}

func (cf *ConfigurableFormatter) truncatePrefix(str string, maxLength int) string {
	if len(str) > maxLength {
		startingIndex := len(str) - maxLength + len(EllipsisString)
		return EllipsisString + str[startingIndex:]
	}

	return str
}

func (cf *ConfigurableFormatter) formatCorrelationElements(correlation proto.LogCorrelationMessage) string {
	shard := correlation.GetShard()
	epoch := correlation.GetEpoch()
	round := correlation.GetRound()
	subRound := correlation.GetSubRound()
	formattedElements := fmt.Sprintf("[%s/%d/%d/%s]", shard, epoch, round, subRound)

	return cf.padRight(formattedElements, CorrelationElementsFixedLength)
}
