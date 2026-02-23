package proxystatus

import (
	"net/http"
	"strconv"
	"strings"
)

type ProxyStatusEntry struct {
	ProxyName string

	Error          *proxyError
	NextHop        string
	NextProtocol   string
	ReceivedStatus int
	Details        string
}

// Returns a Proxy-Status header value and a corresponding recommended status code
func GetProxyStatus(pse ProxyStatusEntry) (proxyStatus string, statusCode int) {
	var b strings.Builder

	// reduces reallocations, 64 should be a decent tradeoff between memory usage and reallocs
	b.Grow(len(pse.ProxyName) + 64)

	b.WriteString(pse.ProxyName)

	if e := pse.Error; e != nil {
		b.WriteString("; error=")
		b.WriteString(e.ErrorString())
		statusCode = e.ErrorCode()
	}
	if pse.NextHop != "" {
		b.WriteString("; next-hop=")
		b.WriteString(pse.NextHop)
	}
	if pse.NextProtocol != "" {
		b.WriteString("; next-protocol=")
		b.WriteString(pse.NextProtocol)
	}
	if pse.ReceivedStatus != 0 {
		b.WriteString("; received_status=")
		b.WriteString(strconv.Itoa(pse.ReceivedStatus))
	}

	if pse.Details != "" {
		b.WriteString("; details=\"")
		b.WriteString(pse.Details)
		b.WriteByte('"')
	}

	return b.String(), statusCode
}

// Add a proxy status header and return an associated recommended status code
func AddProxyStatusHeader(h *http.Header, pse ProxyStatusEntry) (statusCode int) {
	pmsg, statusCode := GetProxyStatus(pse)

	pStatusHdr := h.Get("Proxy-Status")
	if pStatusHdr == "" {
		h.Set("Proxy-Status", pmsg)
		return
	}

	h.Set("Proxy-Status", pStatusHdr+", "+pmsg)
	return
}
