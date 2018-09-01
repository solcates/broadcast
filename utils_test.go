package broadcast

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func Test_getLocalIPs(t *testing.T) {
	tests := []struct {
		name    string
		wantIps bool
	}{
		{
			name:    "Default",
			wantIps: true,
		}, {
			name:    "Error",
			wantIps: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantIps {
				// override
				netInterface = func() (out []net.Addr, err error) {
					return out, errors.New("forced error")
				}
			}
			gotIps := getLocalIPs()
			// make sure we get something back
			assert.NotNil(t, gotIps)

		})
	}
}
