package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestBloqueioDeIP(t *testing.T) {
	IpList := CreateIpList()

	testes := []struct {
		ip         string
		tentativas int
		deveErrar  bool
	}{
		{"192.168.0.1", 1, false},
		{"192.168.0.2", 2, false},
		{"192.168.0.3", MAXIMUM_REQUESTS, true},
		{"192.168.0.4", MAXIMUM_REQUESTS + 1, true},
	}

	for _, teste := range testes {
		t.Run(teste.ip, func(t *testing.T) {
			var err error

			for i := 0; i < teste.tentativas; i++ {
				err = IpList.AddAndCheckIP(teste.ip)
			}

			if teste.deveErrar {
				if err == nil {
					t.Errorf("Esperava erro para IP: %s", teste.ip)
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro para IP: %s, mas veio %v", teste.ip, err)
				}
			}
		})
	}
}

func TestIpBloquado(t *testing.T) {
	const IP = "192.168.1.0"
	ipList := CreateIpList()
	ipAddr := &IpAddr{
		Count:   4,
		Blocked: true,
		Expire:  time.Now().Add(12 * time.Hour),
	}
	ipList.IpAddr[IP] = ipAddr

	var ERROR_MESSAGE = fmt.Sprintf("Endereço bloqueado por rate limits, tente novamente em %d", ipAddr.Expire.Hour())
	err := ipList.AddAndCheckIP(IP)
	if err != nil {
		if err.Error() != ERROR_MESSAGE {
			t.Errorf("Mensagem de erro incorreta:\n    %s", err.Error())
		}
	} else {
		t.Errorf("Esperava erro mas não veio")
	}

}
