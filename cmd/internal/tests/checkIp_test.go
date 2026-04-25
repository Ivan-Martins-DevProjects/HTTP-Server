package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/Ivan-Martins-DevProjects/HTTP-Server/cmd/internal"
)

func TestBloqueioDeIP(t *testing.T) {
	IpList := internal.CreateIpList()

	testes := []struct {
		ip         string
		tentativas int
		deveErrar  bool
	}{
		{"192.168.0.1", 1, false},
		{"192.168.0.2", 2, false},
		{"192.168.0.3", internal.MAXIMUM_REQUESTS, true},
		{"192.168.0.4", internal.MAXIMUM_REQUESTS + 1, true},
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

func TestIpBloqueado(t *testing.T) {
	const IP = "192.168.1.0"
	ipList := internal.CreateIpList()
	ipAddr := &internal.IpAddr{
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
