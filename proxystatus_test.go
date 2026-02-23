package proxystatus_test

import (
	"net/http"
	"testing"

	"github.com/FHRNet/proxystatus"
)

const (
	testProxyName = "ExampleCDN"
)

func Test_GetProxyStatus(t *testing.T) {
	tests := []struct {
		name           string
		pse            proxystatus.ProxyStatusEntry
		expectedString string
		expectedCode   int
	}{
		{
			name:           "No args returns an empty string",
			pse:            proxystatus.ProxyStatusEntry{},
			expectedString: "",
		},
		{
			name: "Supplying a ProxyName should return only that",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
			},
			expectedString: testProxyName,
		},
		{
			name: "Test DestinationIPUnroutable",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
				Error:     &proxystatus.DestinationIPUnroutable,
			},
			expectedString: testProxyName + "; error=destination_ip_unroutable",
			expectedCode:   proxystatus.DestinationIPUnroutable.ErrorCode(),
		},
		{
			name: "Combine Error and NextHop",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
				Error:     &proxystatus.ConnectionTimeout,
				NextHop:   "backend.example.org:8001",
			},
			expectedString: testProxyName + "; error=connection_timeout; next-hop=backend.example.org:8001",
			expectedCode:   proxystatus.ConnectionTimeout.ErrorCode(),
		},
		{
			name: "Combine NextHop and NextProtocol",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName:    testProxyName,
				NextHop:      "backend.example.org:8001",
				NextProtocol: "h2",
			},
			expectedString: testProxyName + "; next-hop=backend.example.org:8001; next-protocol=h2",
			expectedCode:   0,
		},
		{
			name: "Combine Error and ReceivedStatus and Details",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName:      testProxyName,
				Error:          &proxystatus.HTTPRequestError,
				ReceivedStatus: 403,
				Details:        "Origin denied the request",
			},
			expectedString: testProxyName + "; error=http_request_error; received_status=403; details=\"Origin denied the request\"",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rcvdString, rcvdCode := proxystatus.GetProxyStatus(tt.pse)
			if rcvdString != tt.expectedString {
				t.Errorf("ProxyStatus.GetProxyStatus() = %v, want %v", rcvdString, tt.expectedString)
			}
			if rcvdCode != tt.expectedCode {
				t.Errorf("ProxyStatus.GetProxyStatus() = %v, want %v", rcvdCode, tt.expectedCode)
			}
		})
	}
}

func Test_ProxyStatusHeader(t *testing.T) {
	tests := []struct {
		name           string
		hdr            http.Header
		pse            proxystatus.ProxyStatusEntry
		expectedHeader string
		expectedCode   int
	}{
		{
			name:           "No args does not append the header",
			pse:            proxystatus.ProxyStatusEntry{},
			expectedHeader: "",
		},
		{
			name: "Proxy name shows up bare",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
			},
			expectedHeader: testProxyName,
		},
		{
			name: "Header is set correctly with args",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
				Error:     &proxystatus.DestinationIPUnroutable,
				NextHop:   "backend.example.org:8001",
			},
			expectedHeader: testProxyName + "; error=destination_ip_unroutable; next-hop=backend.example.org:8001",
		},
		{
			name: "Header is correctly appended",
			pse: proxystatus.ProxyStatusEntry{
				ProxyName: testProxyName,
				Error:     &proxystatus.DestinationIPUnroutable,
				NextHop:   "backend.example.org:8001",
			},
			hdr: func() http.Header {
				h := http.Header{}
				h.Add("Proxy-Status", "SomeOtherProxy")
				return h
			}(),
			expectedHeader: "SomeOtherProxy, " + testProxyName + "; error=destination_ip_unroutable; next-hop=backend.example.org:8001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hdr
			if h == nil {
				h = http.Header{}
			}

			proxystatus.AddProxyStatusHeader(&h, tt.pse)
			rcvdHeader := h.Get("Proxy-Status")
			if rcvdHeader != tt.expectedHeader {
				t.Errorf("ProxyStatus.AddProxyStatusHeader() = %v, want %v", rcvdHeader, tt.expectedHeader)
			}
		})
	}
}

func Benchmark_GetProxyStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		proxystatus.GetProxyStatus(proxystatus.ProxyStatusEntry{
			ProxyName:      testProxyName,
			Error:          &proxystatus.HTTPRequestError,
			ReceivedStatus: 403,
			Details:        "Origin denied the request",
			NextHop:        "backend.example.org:8001",
			NextProtocol:   "h2",
		})
	}
}
