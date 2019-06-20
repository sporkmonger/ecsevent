package ecsevent

import (
	"fmt"
	"reflect"
)

// Field name constants for the Elastic Common Schema.
// See: https://www.elastic.co/guide/en/ecs/current/ecs-field-reference.html
var (
	FieldTimestamp                    = "@timestamp"
	FieldLabels                       = "labels"
	FieldTags                         = "tags"
	FieldMessage                      = "message"
	FieldAgentEphemeralID             = "agent.ephemeral_id"
	FieldAgentID                      = "agent.id"
	FieldAgentName                    = "agent.name"
	FieldAgentType                    = "agent.type"
	FieldAgentVersion                 = "agent.version"
	FieldClientAddress                = "client.address"
	FieldClientBytes                  = "client.bytes"
	FieldClientDomain                 = "client.domain"
	FieldClientIP                     = "client.ip"
	FieldClientMAC                    = "client.mac"
	FieldClientPackets                = "client.packets"
	FieldClientPort                   = "client.port"
	FieldClientGeoCityName            = "client.geo.city_name"
	FieldClientGeoContinentName       = "client.geo.continent_name"
	FieldClientGeoCountryISOCode      = "client.geo.country_iso_code"
	FieldClientGeoCountryName         = "client.geo.country_name"
	FieldClientGeoLocation            = "client.geo.location"
	FieldClientGeoName                = "client.geo.name"
	FieldClientGeoRegionISOCode       = "client.geo.region_iso_code"
	FieldClientGeoRegionName          = "client.geo.region_name"
	FieldCloudAccountID               = "cloud.account.id"
	FieldCloudAvailabilityZone        = "cloud.availability_zone"
	FieldCloudInstanceID              = "cloud.instance.id"
	FieldCloudInstanceName            = "cloud.instance.name"
	FieldCloudMachineType             = "cloud.machine.type"
	FieldCloudProvider                = "cloud.provider"
	FieldCloudRegion                  = "cloud.region"
	FieldContainerID                  = "container.id"
	FieldContainerImageName           = "container.image.name"
	FieldContainerImageTag            = "container.image.tag"
	FieldContainerLabels              = "container.labels"
	FieldContainerName                = "container.name"
	FieldContainerRuntime             = "container.runtime"
	FieldDestinationAddress           = "destination.address"
	FieldDestinationBytes             = "destination.bytes"
	FieldDestinationDomain            = "destination.domain"
	FieldDestinationIP                = "destination.ip"
	FieldDestinationMAC               = "destination.mac"
	FieldDestinationPackets           = "destination.packets"
	FieldDestinationPort              = "destination.port"
	FieldDestinationGeoCityName       = "destination.geo.city_name"
	FieldDestinationGeoContinentName  = "destination.geo.continent_name"
	FieldDestinationGeoCountryISOCode = "destination.geo.country_iso_code"
	FieldDestinationGeoCountryName    = "destination.geo.country_name"
	FieldDestinationGeoLocation       = "destination.geo.location"
	FieldDestinationGeoName           = "destination.geo.name"
	FieldDestinationGeoRegionISOCode  = "destination.geo.region_iso_code"
	FieldDestinationGeoRegionName     = "destination.geo.region_name"
	FieldDestinationUserEmail         = "destination.user.email"
	FieldDestinationUserFullName      = "destination.user.full_name"
	FieldDestinationUserGroupID       = "destination.user.group.id"
	FieldDestinationUserGroupName     = "destination.user.group.name"
	FieldDestinationUserHash          = "destination.user.hash"
	FieldDestinationUserID            = "destination.user.id"
	FieldDestinationUserName          = "destination.user.name"
	FieldECSVersion                   = "ecs.version"
	FieldErrorCode                    = "error.code"
	FieldErrorID                      = "error.id"
	FieldErrorMessage                 = "error.message"
	FieldErrorStackTrace              = "error.stack_trace"
	FieldEventAction                  = "event.action"
	FieldEventCategory                = "event.category"
	FieldEventCreated                 = "event.created"
	FieldEventDataset                 = "event.dataset"
	FieldEventDuration                = "event.duration"
	FieldEventEnd                     = "event.end"
	FieldEventHash                    = "event.hash"
	FieldEventKind                    = "event.kind"
	FieldEventModule                  = "event.module"
	FieldEventOriginal                = "event.original"
	FieldEventOutcome                 = "event.outcome"
	FieldEventRiskScore               = "event.risk_score"
	FieldEventRiskScoreNorm           = "event.risk_score_norm"
	FieldEventSeverity                = "event.severity"
	FieldEventStart                   = "event.start"
	FieldEventSubevents               = "event.subevents"
	FieldEventTimezone                = "event.timezone"
	FieldEventType                    = "event.type"
	FieldFileCTime                    = "file.ctime"
	FieldFileDevice                   = "file.device"
	FieldFileExtension                = "file.extension"
	FieldFileGID                      = "file.gid"
	FieldFileGroup                    = "file.group"
	FieldFileINode                    = "file.inode"
	FieldFileMode                     = "file.mode"
	FieldFileMTime                    = "file.mtime"
	FieldFileOwner                    = "file.owner"
	FieldFilePath                     = "file.path"
	FieldFileSize                     = "file.size"
	FieldFileTargetPath               = "file.target_path"
	FieldFileType                     = "file.type"
	FieldFileUID                      = "file.uid"
	FieldGroupID                      = "group.id"
	FieldGroupName                    = "group.name"
	FieldHostArchitecture             = "host.architecture"
	FieldHostHostname                 = "host.hostname"
	FieldHostID                       = "host.id"
	FieldHostIP                       = "host.ip"
	FieldHostMAC                      = "host.mac"
	FieldHostName                     = "host.name"
	FieldHostType                     = "host.type"
	FieldHostGeoCityName              = "host.geo.city_name"
	FieldHostGeoContinentName         = "host.geo.continent_name"
	FieldHostGeoCountryISOCode        = "host.geo.country_iso_code"
	FieldHostGeoCountryName           = "host.geo.country_name"
	FieldHostGeoLocation              = "host.geo.location"
	FieldHostGeoName                  = "host.geo.name"
	FieldHostGeoRegionISOCode         = "host.geo.region_iso_code"
	FieldHostGeoRegionName            = "host.geo.region_name"
	FieldHostOSFamily                 = "host.os.family"
	FieldHostOSFull                   = "host.os.full"
	FieldHostOSKernel                 = "host.os.kernel"
	FieldHostOSName                   = "host.os.name"
	FieldHostOSPlatform               = "host.os.platform"
	FieldHostOSVersion                = "host.os.version"
	FieldHostUserEmail                = "host.user.email"
	FieldHostUserFullName             = "host.user.full_name"
	FieldHostUserGroupID              = "host.user.group.id"
	FieldHostUserGroupName            = "host.user.group.name"
	FieldHostUserHash                 = "host.user.hash"
	FieldHostUserID                   = "host.user.id"
	FieldHostUserName                 = "host.user.name"
	FieldHTTPRequestBodyBytes         = "http.request.body.bytes"
	FieldHTTPRequestBodyContent       = "http.request.body.content"
	FieldHTTPRequestBytes             = "http.request.bytes"
	FieldHTTPRequestMethod            = "http.request.method"
	FieldHTTPRequestReferrer          = "http.request.referrer"
	FieldHTTPResponseBodyBytes        = "http.response.body.bytes"
	FieldHTTPResponseBodyContent      = "http.response.body.content"
	FieldHTTPResponseBytes            = "http.response.bytes"
	FieldHTTPResponseStatusCode       = "http.response.status_code"
	FieldHTTPVersion                  = "http.version"
	FieldLogLevel                     = "log.level"
	FieldLogOriginal                  = "log.original"
	FieldNetworkApplication           = "network.application"
	FieldNetworkBytes                 = "network.bytes"
	FieldNetworkCommunityID           = "network.community_id"
	FieldNetworkDirection             = "network.direction"
	FieldNetworkForwardedIP           = "network.forwarded_ip"
	FieldNetworkIANANumber            = "network.iana_number"
	FieldNetworkName                  = "network.name"
	FieldNetworkPackets               = "network.packets"
	FieldNetworkProtocol              = "network.protocol"
	FieldNetworkTransport             = "network.transport"
	FieldNetworkType                  = "network.type"
	FieldObserverHostname             = "observer.hostname"
	FieldObserverIP                   = "observer.ip"
	FieldObserverMAC                  = "observer.mac"
	FieldObserverSerialNumber         = "observer.serial_number"
	FieldObserverType                 = "observer.type"
	FieldObserverVendor               = "observer.vendor"
	FieldObserverVersion              = "observer.version"
	FieldObserverOSFamily             = "observer.os.family"
	FieldObserverOSFull               = "observer.os.full"
	FieldObserverOSKernel             = "observer.os.kernel"
	FieldObserverOSName               = "observer.os.name"
	FieldObserverOSPlatform           = "observer.os.platform"
	FieldObserverOSVersion            = "observer.os.version"
	FieldOrganizationID               = "organization.id"
	FieldOrganizationName             = "organization.name"
	FieldProcessArgs                  = "process.args"
	FieldProcessExecutable            = "process.executable"
	FieldProcessName                  = "process.name"
	FieldProcessPID                   = "process.pid"
	FieldProcessPPID                  = "process.ppid"
	FieldProcessStart                 = "process.start"
	FieldProcessThreadID              = "process.thread.id"
	FieldProcessTitle                 = "process.title"
	FieldProcessWorkingDirectory      = "process.working_directory"
	FieldRelatedIP                    = "related.ip"
	FieldServerAddress                = "server.address"
	FieldServerBytes                  = "server.bytes"
	FieldServerDomain                 = "server.domain"
	FieldServerIP                     = "server.ip"
	FieldServerMAC                    = "server.mac"
	FieldServerPackets                = "server.packets"
	FieldServerPort                   = "server.port"
	FieldServerGeoCityName            = "server.geo.city_name"
	FieldServerGeoContinentName       = "server.geo.continent_name"
	FieldServerGeoCountryISOCode      = "server.geo.country_iso_code"
	FieldServerGeoCountryName         = "server.geo.country_name"
	FieldServerGeoLocation            = "server.geo.location"
	FieldServerGeoName                = "server.geo.name"
	FieldServerGeoRegionISOCode       = "server.geo.region_iso_code"
	FieldServerGeoRegionName          = "server.geo.region_name"
	FieldServerUserEmail              = "server.user.email"
	FieldServerUserFullName           = "server.user.full_name"
	FieldServerUserGroupID            = "server.user.group.id"
	FieldServerUserGroupName          = "server.user.group.name"
	FieldServerUserHash               = "server.user.hash"
	FieldServerUserID                 = "server.user.id"
	FieldServerUserName               = "server.user.name"
	FieldServiceEphemeralID           = "service.ephemeral_id"
	FieldServiceID                    = "service.id"
	FieldServiceName                  = "service.name"
	FieldServiceState                 = "service.state"
	FieldServiceType                  = "service.type"
	FieldServiceVersion               = "service.version"
	FieldSourceAddress                = "source.address"
	FieldSourceBytes                  = "source.bytes"
	FieldSourceDomain                 = "source.domain"
	FieldSourceIP                     = "source.ip"
	FieldSourceMAC                    = "source.mac"
	FieldSourcePackets                = "source.packets"
	FieldSourcePort                   = "source.port"
	FieldSourceGeoCityName            = "source.geo.city_name"
	FieldSourceGeoContinentName       = "source.geo.continent_name"
	FieldSourceGeoCountryISOCode      = "source.geo.country_iso_code"
	FieldSourceGeoCountryName         = "source.geo.country_name"
	FieldSourceGeoLocation            = "source.geo.location"
	FieldSourceGeoName                = "source.geo.name"
	FieldSourceGeoRegionISOCode       = "source.geo.region_iso_code"
	FieldSourceGeoRegionName          = "source.geo.region_name"
	FieldSourceUserEmail              = "source.user.email"
	FieldSourceUserFullName           = "source.user.full_name"
	FieldSourceUserGroupID            = "source.user.group.id"
	FieldSourceUserGroupName          = "source.user.group.name"
	FieldSourceUserHash               = "source.user.hash"
	FieldSourceUserID                 = "source.user.id"
	FieldSourceUserName               = "source.user.name"
	FieldURLDomain                    = "url.domain"
	FieldURLFragment                  = "url.fragment"
	FieldURLFull                      = "url.full"
	FieldURLOriginal                  = "url.original"
	FieldURLPassword                  = "url.password"
	FieldURLPath                      = "url.path"
	FieldURLPort                      = "url.port"
	FieldURLQuery                     = "url.query"
	FieldURLScheme                    = "url.scheme"
	FieldURLUsername                  = "url.username"
	FieldUserEmail                    = "user.email"
	FieldUserFullName                 = "user.full_name"
	FieldUserGroupID                  = "user.group.id"
	FieldUserGroupName                = "user.group.name"
	FieldUserHash                     = "user.hash"
	FieldUserID                       = "user.id"
	FieldUserName                     = "user.name"
	FieldUserAgentDeviceName          = "user_agent.device.name"
	FieldUserAgentName                = "user_agent.name"
	FieldUserAgentOriginal            = "user_agent.original"
	FieldUserAgentVersion             = "user_agent.version"
)

// Sanity check, since we're working with lots of map[string]interface{}
var fieldKinds = map[string]reflect.Kind{
	FieldTimestamp:                    reflect.Struct, // time.Time
	FieldLabels:                       reflect.Map,    // map[string]string
	FieldTags:                         reflect.Slice,  // []string
	FieldMessage:                      reflect.String,
	FieldAgentEphemeralID:             reflect.String,
	FieldAgentID:                      reflect.String,
	FieldAgentName:                    reflect.String,
	FieldAgentType:                    reflect.String,
	FieldAgentVersion:                 reflect.String,
	FieldClientAddress:                reflect.String,
	FieldClientBytes:                  reflect.Int,
	FieldClientDomain:                 reflect.String,
	FieldClientIP:                     reflect.String,
	FieldClientMAC:                    reflect.String,
	FieldClientPackets:                reflect.Int,
	FieldClientPort:                   reflect.Int,
	FieldClientGeoCityName:            reflect.String,
	FieldClientGeoContinentName:       reflect.String,
	FieldClientGeoCountryISOCode:      reflect.String,
	FieldClientGeoCountryName:         reflect.String,
	FieldClientGeoLocation:            reflect.Map, // map[string]float
	FieldClientGeoName:                reflect.String,
	FieldClientGeoRegionISOCode:       reflect.String,
	FieldClientGeoRegionName:          reflect.String,
	FieldCloudAccountID:               reflect.String,
	FieldCloudAvailabilityZone:        reflect.String,
	FieldCloudInstanceID:              reflect.String,
	FieldCloudInstanceName:            reflect.String,
	FieldCloudMachineType:             reflect.String,
	FieldCloudProvider:                reflect.String,
	FieldCloudRegion:                  reflect.String,
	FieldContainerID:                  reflect.String,
	FieldContainerImageName:           reflect.String,
	FieldContainerImageTag:            reflect.String,
	FieldContainerLabels:              reflect.Map, // map[string]string
	FieldContainerName:                reflect.String,
	FieldContainerRuntime:             reflect.String,
	FieldDestinationAddress:           reflect.String,
	FieldDestinationBytes:             reflect.Int,
	FieldDestinationDomain:            reflect.String,
	FieldDestinationIP:                reflect.String,
	FieldDestinationMAC:               reflect.String,
	FieldDestinationPackets:           reflect.Int,
	FieldDestinationPort:              reflect.Int,
	FieldDestinationGeoCityName:       reflect.String,
	FieldDestinationGeoContinentName:  reflect.String,
	FieldDestinationGeoCountryISOCode: reflect.String,
	FieldDestinationGeoCountryName:    reflect.String,
	FieldDestinationGeoLocation:       reflect.Map, // map[string]float
	FieldDestinationGeoName:           reflect.String,
	FieldDestinationGeoRegionISOCode:  reflect.String,
	FieldDestinationGeoRegionName:     reflect.String,
	FieldDestinationUserEmail:         reflect.String,
	FieldDestinationUserFullName:      reflect.String,
	FieldDestinationUserGroupID:       reflect.String,
	FieldDestinationUserGroupName:     reflect.String,
	FieldDestinationUserHash:          reflect.String,
	FieldDestinationUserID:            reflect.String,
	FieldDestinationUserName:          reflect.String,
	FieldECSVersion:                   reflect.String,
	FieldErrorCode:                    reflect.String,
	FieldErrorID:                      reflect.String,
	FieldErrorMessage:                 reflect.String,
	FieldErrorStackTrace:              reflect.String,
	FieldEventAction:                  reflect.String,
	FieldEventCategory:                reflect.String,
	FieldEventCreated:                 reflect.Struct, // time.Time
	FieldEventDataset:                 reflect.String,
	FieldEventDuration:                reflect.Int,
	FieldEventEnd:                     reflect.Struct, // time.Time
	FieldEventHash:                    reflect.String,
	FieldEventKind:                    reflect.String,
	FieldEventModule:                  reflect.String,
	FieldEventOriginal:                reflect.String,
	FieldEventOutcome:                 reflect.String,
	FieldEventRiskScore:               reflect.Float64,
	FieldEventRiskScoreNorm:           reflect.Float64,
	FieldEventSeverity:                reflect.Int,
	FieldEventStart:                   reflect.Struct, // time.Time
	FieldEventTimezone:                reflect.String,
	FieldEventType:                    reflect.String,
	FieldFileCTime:                    reflect.Struct, // time.Time
	FieldFileDevice:                   reflect.String,
	FieldFileExtension:                reflect.String,
	FieldFileGID:                      reflect.String,
	FieldFileGroup:                    reflect.String,
	FieldFileINode:                    reflect.String,
	FieldFileMode:                     reflect.String,
	FieldFileMTime:                    reflect.Struct, // time.Time
	FieldFileOwner:                    reflect.String,
	FieldFileSize:                     reflect.Int,
	FieldFileTargetPath:               reflect.String,
	FieldFileType:                     reflect.String,
	FieldFileUID:                      reflect.String,
	FieldGroupID:                      reflect.String,
	FieldGroupName:                    reflect.String,
	FieldHostArchitecture:             reflect.String,
	FieldHostHostname:                 reflect.String,
	FieldHostID:                       reflect.String,
	FieldHostIP:                       reflect.String,
	FieldHostMAC:                      reflect.String,
	FieldHostName:                     reflect.String,
	FieldHostType:                     reflect.String,
	FieldHostGeoCityName:              reflect.String,
	FieldHostGeoContinentName:         reflect.String,
	FieldHostGeoCountryISOCode:        reflect.String,
	FieldHostGeoCountryName:           reflect.String,
	FieldHostGeoLocation:              reflect.Map, // map[string]float
	FieldHostGeoName:                  reflect.String,
	FieldHostGeoRegionISOCode:         reflect.String,
	FieldHostGeoRegionName:            reflect.String,
	FieldHostOSFamily:                 reflect.String,
	FieldHostOSFull:                   reflect.String,
	FieldHostOSKernel:                 reflect.String,
	FieldHostOSName:                   reflect.String,
	FieldHostOSPlatform:               reflect.String,
	FieldHostOSVersion:                reflect.String,
	FieldHostUserEmail:                reflect.String,
	FieldHostUserFullName:             reflect.String,
	FieldHostUserGroupID:              reflect.String,
	FieldHostUserGroupName:            reflect.String,
	FieldHostUserHash:                 reflect.String,
	FieldHostUserID:                   reflect.String,
	FieldHostUserName:                 reflect.String,
	FieldHTTPRequestBodyBytes:         reflect.Int,
	FieldHTTPRequestBodyContent:       reflect.String,
	FieldHTTPRequestBytes:             reflect.Int,
	FieldHTTPRequestMethod:            reflect.String,
	FieldHTTPRequestReferrer:          reflect.String,
	FieldHTTPResponseBodyBytes:        reflect.Int,
	FieldHTTPResponseBodyContent:      reflect.String,
	FieldHTTPResponseBytes:            reflect.Int,
	FieldHTTPResponseStatusCode:       reflect.Int,
	FieldHTTPVersion:                  reflect.String,
	FieldLogLevel:                     reflect.String,
	FieldLogOriginal:                  reflect.String,
	FieldNetworkApplication:           reflect.String,
	FieldNetworkBytes:                 reflect.Int,
	FieldNetworkCommunityID:           reflect.String,
	FieldNetworkDirection:             reflect.String,
	FieldNetworkForwardedIP:           reflect.String,
	FieldNetworkIANANumber:            reflect.String,
	FieldNetworkName:                  reflect.String,
	FieldNetworkPackets:               reflect.Int,
	FieldNetworkProtocol:              reflect.String,
	FieldNetworkTransport:             reflect.String,
	FieldNetworkType:                  reflect.String,
	FieldObserverHostname:             reflect.String,
	FieldObserverIP:                   reflect.String,
	FieldObserverMAC:                  reflect.String,
	FieldObserverSerialNumber:         reflect.String,
	FieldObserverType:                 reflect.String,
	FieldObserverVendor:               reflect.String,
	FieldObserverVersion:              reflect.String,
	FieldObserverOSFamily:             reflect.String,
	FieldObserverOSFull:               reflect.String,
	FieldObserverOSKernel:             reflect.String,
	FieldObserverOSName:               reflect.String,
	FieldObserverOSPlatform:           reflect.String,
	FieldObserverOSVersion:            reflect.String,
	FieldOrganizationID:               reflect.String,
	FieldOrganizationName:             reflect.String,
	FieldProcessArgs:                  reflect.Slice, // []string
	FieldProcessExecutable:            reflect.String,
	FieldProcessName:                  reflect.String,
	FieldProcessPID:                   reflect.Int,
	FieldProcessPPID:                  reflect.Int,
	FieldProcessStart:                 reflect.Struct, // time.Time
	FieldProcessThreadID:              reflect.Int,
	FieldProcessTitle:                 reflect.String,
	FieldProcessWorkingDirectory:      reflect.String,
	FieldRelatedIP:                    reflect.String,
	FieldServerAddress:                reflect.String,
	FieldServerBytes:                  reflect.Int,
	FieldServerDomain:                 reflect.String,
	FieldServerIP:                     reflect.String,
	FieldServerMAC:                    reflect.String,
	FieldServerPackets:                reflect.Int,
	FieldServerPort:                   reflect.String,
	FieldServerGeoCityName:            reflect.String,
	FieldServerGeoContinentName:       reflect.String,
	FieldServerGeoCountryISOCode:      reflect.String,
	FieldServerGeoCountryName:         reflect.String,
	FieldServerGeoLocation:            reflect.Map, // map[string]float
	FieldServerGeoName:                reflect.String,
	FieldServerGeoRegionISOCode:       reflect.String,
	FieldServerGeoRegionName:          reflect.String,
	FieldServerUserEmail:              reflect.String,
	FieldServerUserFullName:           reflect.String,
	FieldServerUserGroupID:            reflect.String,
	FieldServerUserGroupName:          reflect.String,
	FieldServerUserHash:               reflect.String,
	FieldServerUserID:                 reflect.String,
	FieldServerUserName:               reflect.String,
	FieldServiceEphemeralID:           reflect.String,
	FieldServiceID:                    reflect.String,
	FieldServiceName:                  reflect.String,
	FieldServiceState:                 reflect.String,
	FieldServiceType:                  reflect.String,
	FieldServiceVersion:               reflect.String,
	FieldSourceAddress:                reflect.String,
	FieldSourceBytes:                  reflect.Int,
	FieldSourceDomain:                 reflect.String,
	FieldSourceIP:                     reflect.String,
	FieldSourceMAC:                    reflect.String,
	FieldSourcePackets:                reflect.Int,
	FieldSourcePort:                   reflect.String,
	FieldSourceGeoCityName:            reflect.String,
	FieldSourceGeoContinentName:       reflect.String,
	FieldSourceGeoCountryISOCode:      reflect.String,
	FieldSourceGeoCountryName:         reflect.String,
	FieldSourceGeoLocation:            reflect.Map, // map[string]float
	FieldSourceGeoName:                reflect.String,
	FieldSourceGeoRegionISOCode:       reflect.String,
	FieldSourceGeoRegionName:          reflect.String,
	FieldSourceUserEmail:              reflect.String,
	FieldSourceUserFullName:           reflect.String,
	FieldSourceUserGroupID:            reflect.String,
	FieldSourceUserGroupName:          reflect.String,
	FieldSourceUserHash:               reflect.String,
	FieldSourceUserID:                 reflect.String,
	FieldSourceUserName:               reflect.String,
	FieldURLDomain:                    reflect.String,
	FieldURLFragment:                  reflect.String,
	FieldURLFull:                      reflect.String,
	FieldURLOriginal:                  reflect.String,
	FieldURLPassword:                  reflect.String,
	FieldURLPath:                      reflect.String,
	FieldURLPort:                      reflect.Int,
	FieldURLQuery:                     reflect.String,
	FieldURLScheme:                    reflect.String,
	FieldURLUsername:                  reflect.String,
	FieldUserEmail:                    reflect.String,
	FieldUserFullName:                 reflect.String,
	FieldUserGroupID:                  reflect.String,
	FieldUserGroupName:                reflect.String,
	FieldUserHash:                     reflect.String,
	FieldUserID:                       reflect.String,
	FieldUserName:                     reflect.String,
	FieldUserAgentDeviceName:          reflect.String,
	FieldUserAgentName:                reflect.String,
	FieldUserAgentOriginal:            reflect.String,
	FieldUserAgentVersion:             reflect.String,
}

func typeCheck(fieldName string, value interface{}) error {
	kind, ok := fieldKinds[fieldName]
	if ok {
		valueType := reflect.TypeOf(value)
		if kind != valueType.Kind() {
			return fmt.Errorf("unexpected value for field '%s', refer to ECS specification", fieldName)
		}
	}
	return nil
}
