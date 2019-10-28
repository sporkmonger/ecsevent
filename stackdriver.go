package ecsevent

import (
	"strings"
)

var (
	fieldStackdriverTimestamp               = "timestamp"
	fieldStackdriverSeverity                = "severity"
	fieldStackdriverHTTPRequestMethod       = "httpRequest.requestMethod"
	fieldStackdriverHTTPRequestURL          = "httpRequest.requestUrl"
	fieldStackdriverHTTPRequestSize         = "httpRequest.requestSize"
	fieldStackdriverHTTPRequestStatus       = "httpRequest.status"       // response
	fieldStackdriverHTTPRequestResponseSize = "httpRequest.responseSize" // response
	fieldStackdriverHTTPRequestUserAgent    = "httpRequest.userAgent"
	fieldStackdriverHTTPRequestRemoteIP     = "httpRequest.remoteIp"
	fieldStackdriverHTTPRequestServerIP     = "httpRequest.serverIp"
	fieldStackdriverHTTPRequestReferrer     = "httpRequest.referer"
	fieldStackdriverHTTPRequestLatency      = "httpRequest.latency"
	fieldStackdriverHTTPRequestProtocol     = "httpRequest.protocol"
)

// appendStackdriver takes a map in ECS dotted notation and for each field
// with a special value in Stackdriver (e.g. 'severity'), it will apply the
// necessary transforms and inject these values into a new map.
func appendStackdriver(entry map[string]interface{}) map[string]interface{} {
	newEntry := make(map[string]interface{})
	for key, value := range entry {
		// copy existing fields over unmodified
		newEntry[key] = value

		switch key {
		case FieldTimestamp:
			newEntry[fieldStackdriverTimestamp] = value
		case FieldLogLevel:
			if level, ok := value.(string); ok {
				newEntry[fieldStackdriverSeverity] = stackdriverSeverity(level)
			}
		case FieldHTTPRequestMethod:
			newEntry[fieldStackdriverHTTPRequestMethod] = value
		case FieldURLFull:
			newEntry[fieldStackdriverHTTPRequestURL] = value
		case FieldHTTPRequestBytes:
			if bytes, ok := value.(int64); ok {
				newEntry[fieldStackdriverHTTPRequestSize] = string(bytes)
			}
		case FieldHTTPResponseStatusCode:
			newEntry[fieldStackdriverHTTPRequestStatus] = value
		case FieldUserAgentOriginal:
			newEntry[fieldStackdriverHTTPRequestUserAgent] = value
		case FieldHTTPRequestReferrer:
			newEntry[fieldStackdriverHTTPRequestReferrer] = value
		case FieldClientIP:
			newEntry[fieldStackdriverHTTPRequestRemoteIP] = value
		case FieldServerIP:
			newEntry[fieldStackdriverHTTPRequestServerIP] = value
		case FieldHTTPVersion:
			if version, ok := value.(string); ok {
				newEntry[fieldStackdriverHTTPRequestProtocol] = stackdriverHTTPProtocol(version)
			}
		}
	}
	return newEntry
}

func stackdriverSeverity(level string) string {
	// ECS doesn't specify accepted values for log.level
	switch strings.ToLower(level) {
	case "t":
		fallthrough
	case "trc":
		fallthrough
	case "trace":
		fallthrough
	case "d":
		fallthrough
	case "dbg":
		fallthrough
	case "debug":
		return "DEBUG"
	case "i":
		fallthrough
	case "inf":
		fallthrough
	case "informational":
		fallthrough
	case "info":
		return "INFO"
	case "n":
		fallthrough
	case "not":
		fallthrough
	case "ntc":
		fallthrough
	case "notice":
		return "NOTICE"
	case "w":
		fallthrough
	case "wrn":
		fallthrough
	case "warn":
		fallthrough
	case "warning":
		return "WARNING"
	case "e":
		fallthrough
	case "err":
		fallthrough
	case "error":
		return "ERROR"
	case "c":
		fallthrough
	case "crt":
		fallthrough
	case "crit":
		fallthrough
	case "critical":
		return "CRITICAL"
	case "a":
		fallthrough
	case "alr":
		fallthrough
	case "alrt":
		fallthrough
	case "alrm":
		fallthrough
	case "alarm":
		fallthrough
	case "alert":
		return "ALERT"
	case "f":
		fallthrough
	case "ftl":
		fallthrough
	case "fat":
		fallthrough
	case "fatal":
		fallthrough
	case "emg":
		fallthrough
	case "emrg":
		fallthrough
	case "emergency":
		return "EMERGENCY"
	default:
		return "DEFAULT"
	}
}

func stackdriverHTTPProtocol(version string) string {
	// ECS uses just a version number
	switch strings.ToLower(version) {
	case "1.0":
		return "HTTP/1.0"
	case "1.1":
		return "HTTP/1.1"
	case "2":
		return "HTTP/2"
	default:
		return version
	}
}
