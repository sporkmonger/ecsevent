package ecsevent

import (
	"strings"
)

var (
	fieldStackDriverTimestamp               = "timestamp"
	fieldStackDriverSeverity                = "severity"
	fieldStackDriverHTTPRequestMethod       = "httpRequest.requestMethod"
	fieldStackDriverHTTPRequestURL          = "httpRequest.requestUrl"
	fieldStackDriverHTTPRequestSize         = "httpRequest.requestSize"
	fieldStackDriverHTTPRequestStatus       = "httpRequest.status"       // response
	fieldStackDriverHTTPRequestResponseSize = "httpRequest.responseSize" // response
	fieldStackDriverHTTPRequestUserAgent    = "httpRequest.userAgent"
	fieldStackDriverHTTPRequestRemoteIP     = "httpRequest.remoteIp"
	fieldStackDriverHTTPRequestServerIP     = "httpRequest.serverIp"
	fieldStackDriverHTTPRequestReferrer     = "httpRequest.referer"
	fieldStackDriverHTTPRequestLatency      = "httpRequest.latency"
	fieldStackDriverHTTPRequestProtocol     = "httpRequest.protocol"
)

// appendStackDriver takes a map in ECS dotted notation and for each field
// with a special value in StackDriver (e.g. 'severity'), it will apply the
// necessary transforms and inject these values into a new map.
func appendStackDriver(entry map[string]interface{}) map[string]interface{} {
	newEntry := make(map[string]interface{})
	for key, value := range entry {
		// copy existing fields over unmodified
		newEntry[key] = value

		switch key {
		case FieldTimestamp:
			newEntry[fieldStackDriverTimestamp] = value
		case FieldLogLevel:
			if level, ok := value.(string); ok {
				newEntry[fieldStackDriverSeverity] = stackDriverSeverity(level)
			}
		case FieldHTTPRequestMethod:
			newEntry[fieldStackDriverHTTPRequestMethod] = value
		case FieldURLFull:
			newEntry[fieldStackDriverHTTPRequestURL] = value
		case FieldHTTPRequestBytes:
			if bytes, ok := value.(int64); ok {
				newEntry[fieldStackDriverHTTPRequestSize] = string(bytes)
			}
		case FieldHTTPResponseStatusCode:
			newEntry[fieldStackDriverHTTPRequestStatus] = value
		case FieldUserAgentOriginal:
			newEntry[fieldStackDriverHTTPRequestUserAgent] = value
		case FieldHTTPRequestReferrer:
			newEntry[fieldStackDriverHTTPRequestReferrer] = value
		case FieldClientIP:
			newEntry[fieldStackDriverHTTPRequestRemoteIP] = value
		case FieldServerIP:
			newEntry[fieldStackDriverHTTPRequestServerIP] = value
		case FieldHTTPVersion:
			if version, ok := value.(string); ok {
				newEntry[fieldStackDriverHTTPRequestProtocol] = stackDriverHTTPProtocol(version)
			}
		}
	}
	return newEntry
}

func stackDriverSeverity(level string) string {
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

func stackDriverHTTPProtocol(version string) string {
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
