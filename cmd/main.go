package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	smtpmock "github.com/mocktools/go-smtp-mock"
)

var signals, logFatalf = make(chan os.Signal, 1), log.Fatalf

// Main entrypoint
func main() {
	if err := run(os.Args); err != nil {
		logFatalf("%s\n", err)
	}
}

// SMTP mock server life cycle runner
func run(args []string, options ...flag.ErrorHandling) error {
	failureScenario := flag.ExitOnError
	if len(options) > 0 {
		failureScenario = options[0]
	}

	configAttr, err := configurationAttrFromCommandLine(args, failureScenario)
	if err != nil {
		return err
	}

	server := smtpmock.New(*configAttr)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if err := server.Start(); err != nil {
		return err
	}

	<-signals

	return server.Stop()
}

// Converts string separated by commas to slice
func toSlice(str string) []string {
	return strings.Split(str, ",")
}

// Creates pointer to ConfigurationAttr based on passed command line arguments
func configurationAttrFromCommandLine(args []string, options ...flag.ErrorHandling) (*smtpmock.ConfigurationAttr, error) {
	failureScenario := flag.ExitOnError
	if len(options) > 0 {
		failureScenario = options[0]
	}

	flags := flag.NewFlagSet(args[0], failureScenario)
	var (
		host                          = flags.String("host", "", "Host address where smtpmock will run. It's equal to 127.0.0.1 by default")
		port                          = flags.Int("port", 0, "Server port number. If not specified it will be assigned dynamically")
		log                           = flags.Bool("log", false, "Enables log server activity. Disabled by default")
		sessionTimeout                = flags.Int("sessionTimeout", 0, "Session timeout in seconds. It's equal to 30 seconds by default")
		failFast                      = flags.Bool("failFast", false, "Enables fail fast scenario. Disabled by default")
		blacklistedHeloDomains        = flags.String("blacklistedHeloDomains", "", "Blacklisted HELO domains, separated by commas")
		blacklistedMailfromEmails     = flags.String("blacklistedMailfromEmails", "", "Blacklisted MAIL FROM emails, separated by commas")
		blacklistedRcpttoEmails       = flags.String("blacklistedRcpttoEmails", "", "Blacklisted RCPT TO emails, separated by commas")
		notRegisteredEmails           = flags.String("notRegisteredEmails", "", "Not registered (non-existent) RCPT TO emails, separated by commas")
		msgSizeLimit                  = flags.Int("msgSizeLimit", 0, "Message body size limit in bytes. It's equal to 10485760 bytes")
		msgGreeting                   = flags.String("msgGreeting", "", "Custom server greeting message")
		msgInvalidCmd                 = flags.String("msgInvalidCmd", "", "Custom invalid command message")
		msgInvalidCmdHeloSequence     = flags.String("msgInvalidCmdHeloSequence", "", "Custom invalid command HELO sequence message")
		msgInvalidCmdHeloArg          = flags.String("msgInvalidCmdHeloArg", "", "Custom invalid command HELO argument message")
		msgHeloBlacklistedDomain      = flags.String("msgHeloBlacklistedDomain", "", "Custom HELO blacklisted domain message")
		msgHeloReceived               = flags.String("msgHeloReceived", "", "Custom HELO received message")
		msgInvalidCmdMailfromSequence = flags.String("msgInvalidCmdMailfromSequence", "", "Custom invalid command MAIL FROM sequence message")
		msgInvalidCmdMailfromArg      = flags.String("msgInvalidCmdMailfromArg", "", "Custom invalid command MAIL FROM argument message")
		msgMailfromBlacklistedEmail   = flags.String("msgMailfromBlacklistedEmail", "", "Custom MAIL FROM blacklisted email message")
		msgMailfromReceived           = flags.String("msgMailfromReceived", "", "Custom MAIL FROM received message")
		msgInvalidCmdRcpttoSequence   = flags.String("msgInvalidCmdRcpttoSequence", "", "Custom invalid command RCPT TO sequence message")
		msgInvalidCmdRcpttoArg        = flags.String("msgInvalidCmdRcpttoArg", "", "Custom invalid command RCPT TO argument message")
		msgRcpttoNotRegisteredEmail   = flags.String("msgRcpttoNotRegisteredEmail", "", "Custom RCPT TO not registered email message")
		msgRcpttoBlacklistedEmail     = flags.String("msgRcpttoBlacklistedEmail", "", "Custom RCPT TO blacklisted email message")
		msgRcpttoReceived             = flags.String("msgRcpttoReceived", "", "Custom RCPT TO received message")
		msgInvalidCmdDataSequence     = flags.String("msgInvalidCmdDataSequence", "", "Custom invalid command DATA sequence message")
		msgDataReceived               = flags.String("msgDataReceived", "", "Custom DATA received message")
		msgMsgSizeIsTooBig            = flags.String("msgMsgSizeIsTooBig", "", "Custom size is too big message")
		msgMsgReceived                = flags.String("msgMsgReceived", "", "Custom received message body message")
		msgQuitCmd                    = flags.String("msgQuitCmd", "", "Custom quit command message")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	return &smtpmock.ConfigurationAttr{
		HostAddress:                   *host,
		PortNumber:                    *port,
		LogToStdout:                   *log,
		LogServerActivity:             *log,
		SessionTimeout:                *sessionTimeout,
		IsCmdFailFast:                 *failFast,
		BlacklistedHeloDomains:        toSlice(*blacklistedHeloDomains),
		BlacklistedMailfromEmails:     toSlice(*blacklistedMailfromEmails),
		BlacklistedRcpttoEmails:       toSlice(*blacklistedRcpttoEmails),
		NotRegisteredEmails:           toSlice(*notRegisteredEmails),
		MsgSizeLimit:                  *msgSizeLimit,
		MsgGreeting:                   *msgGreeting,
		MsgInvalidCmd:                 *msgInvalidCmd,
		MsgInvalidCmdHeloSequence:     *msgInvalidCmdHeloSequence,
		MsgInvalidCmdHeloArg:          *msgInvalidCmdHeloArg,
		MsgHeloBlacklistedDomain:      *msgHeloBlacklistedDomain,
		MsgHeloReceived:               *msgHeloReceived,
		MsgInvalidCmdMailfromSequence: *msgInvalidCmdMailfromSequence,
		MsgInvalidCmdMailfromArg:      *msgInvalidCmdMailfromArg,
		MsgMailfromBlacklistedEmail:   *msgMailfromBlacklistedEmail,
		MsgMailfromReceived:           *msgMailfromReceived,
		MsgInvalidCmdRcpttoSequence:   *msgInvalidCmdRcpttoSequence,
		MsgInvalidCmdRcpttoArg:        *msgInvalidCmdRcpttoArg,
		MsgRcpttoNotRegisteredEmail:   *msgRcpttoNotRegisteredEmail,
		MsgRcpttoBlacklistedEmail:     *msgRcpttoBlacklistedEmail,
		MsgRcpttoReceived:             *msgRcpttoReceived,
		MsgInvalidCmdDataSequence:     *msgInvalidCmdDataSequence,
		MsgDataReceived:               *msgDataReceived,
		MsgMsgSizeIsTooBig:            *msgMsgSizeIsTooBig,
		MsgMsgReceived:                *msgMsgReceived,
		MsgQuitCmd:                    *msgQuitCmd,
	}, nil
}