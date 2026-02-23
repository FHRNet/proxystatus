package proxystatus

type proxyError struct {
	errStr  string
	errCode int
}

// Returns a string associated with the error key
func (pe proxyError) ErrorString() string {
	return pe.errStr
}

// Returns a status code associated with the error key
func (pe proxyError) ErrorCode() int {
	return pe.errCode
}

// https://www.iana.org/assignments/http-proxy-status/http-proxy-status.xhtml
var (
	DNSTimeout = proxyError{"dns_timeout", 504}
	DNSError   = proxyError{"dns_error", 502}

	DestinationNotFound     = proxyError{"destination_not_found", 500}
	DestinationUnavailable  = proxyError{"destination_unavailable", 503}
	DestinationIPProhibited = proxyError{"destination_ip_prohibited", 502}
	DestinationIPUnroutable = proxyError{"destination_ip_unroutable", 502}

	ConnectionRefused        = proxyError{"connection_refused", 502}
	ConnectionTerminated     = proxyError{"connection_terminated", 502}
	ConnectionTimeout        = proxyError{"connection_timeout", 504}
	ConnectionReadTimeout    = proxyError{"connection_read_timeout", 504}
	ConnectionWriteTimeout   = proxyError{"connection_write_timeout", 504}
	ConnectionLimitedReached = proxyError{"connection_limit_reached", 503}

	TLSProtocolError    = proxyError{"tls_protocol_error", 502}
	TLSCertificateError = proxyError{"tls_certificate_error", 502}
	TLSAlertReceived    = proxyError{"tls_alert_received", 502}

	HTTPRequestError                     = proxyError{"http_request_error", 0}
	HTTPRequestDenied                    = proxyError{"http_request_denied", 403}
	HTTPIncompleteResponse               = proxyError{"http_response_incomplete", 502}
	HTTPResponseHeaderSectionTooLarge    = proxyError{"http_response_header_section_size", 502}
	HTTPResponseHeaderFieldLineTooLarge  = proxyError{"http_response_header_size", 502}
	HTTPResponseBodyTooLarge             = proxyError{"http_response_body_size", 502}
	HTTPResponseTrailerSectionTooLarge   = proxyError{"http_response_trailer_section_size", 502}
	HTTPResponseTrailerFieldLineTooLarge = proxyError{"http_response_trailer_size", 502}
	HTTPResponseTransferCodingError      = proxyError{"http_response_transfer_coding", 502}
	HTTPResponseContentCodingError       = proxyError{"http_response_content_coding", 502}
	HTTPResponseTimeout                  = proxyError{"http_response_timeout", 504}
	HTTPUpgradeFailed                    = proxyError{"http_upgrade_failed", 502}
	HTTPProtocolError                    = proxyError{"http_protocol_error", 502}

	ProxyInternalResponse   = proxyError{"proxy_internal_response", 0}
	ProxyInternalError      = proxyError{"proxy_internal_error", 500}
	ProxyConfigurationError = proxyError{"proxy_configuration_error", 500}
	ProxyLoopDetected       = proxyError{"proxy_loop_detected", 502}
)
