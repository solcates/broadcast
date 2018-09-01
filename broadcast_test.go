package broadcast

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func TestNewBroadcaster(t *testing.T) {
	type args struct {
		port    int
		payload string
		timeout int
	}
	tests := []struct {
		name string
		args args
		want *UDPBroadcaster
	}{
		{
			name: "Default",
			args: args{
				port:    8888,
				payload: "hello",
			},
			want: &UDPBroadcaster{
				port:    8888,
				payload: "hello",
				timeout: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUDPBroadcaster(tt.args.port, tt.args.payload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUDPBroadcaster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_broadcaster_Discover(t *testing.T) {

	tests := []struct {
		name             string
		bc               *UDPBroadcaster
		wantSelf         bool
		wantErr          bool
		wantBadbroadcast bool
	}{
		{
			name: "Default",
			bc: &UDPBroadcaster{
				port:     30303,
				payload:  "Discovery: Who is out there?",
				findself: false,
				timeout:  5,
			},
			wantSelf: false,
			wantErr:  false,
		},
		{
			name: "Findself",
			bc: &UDPBroadcaster{
				port:     30303,
				payload:  "Discovery: Who is out there?",
				findself: true,
				timeout:  5,
			},
			wantSelf: true,
			wantErr:  false,
		}, {
			name: "Error - Bad Port",
			bc: &UDPBroadcaster{
				port:     0,
				payload:  "Discovery: Who is out there?",
				findself: false,
				timeout:  5,
			},
			wantSelf: false,
			wantErr:  true,
		}, {
			name: "Error - Bad Payload",
			bc: &UDPBroadcaster{
				port:     30303,
				payload:  "",
				findself: false,
				timeout:  5,
			},
			wantSelf: false,
			wantErr:  true,
		}, {
			name: "Error - Server ResolveUDPAddr",
			bc: &UDPBroadcaster{
				port:     30303,
				payload:  "Discovery: Who is out there?",
				findself: false,
				timeout:  5,
			},
			wantSelf:         false,
			wantErr:          true,
			wantBadbroadcast: true,
		}, {
			name: "Error - Local ResolveUDPAddr",
			bc: &UDPBroadcaster{
				port:     0,
				payload:  "Discovery: Who is out there?",
				findself: false,
				timeout:  5,
			},
			wantSelf:         false,
			wantErr:          true,
			wantBadbroadcast: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantBadbroadcast {
				broadcastAddress = ""
				network = "te"
			}
			gotFound, err := tt.bc.Discover()
			if (err != nil) != tt.wantErr {
				t.Errorf("UDPBroadcaster.Discover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantSelf {
				assert.NotNil(t, gotFound)
			}

		})
	}
}

func ExampleUDPBroadcaster_Discover() {
	bc := NewUDPBroadcaster(8080, "Hello?")
	bc.SetFindself(true)
	ips, err := bc.Discover()
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		fmt.Println(ip)
	}
}
