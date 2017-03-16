package wrapper

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_Hoverfly_isLocal_WhenLocalhost(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "localhost"}

	result := hoverfly.isLocal()

	Expect(result).To(BeTrue())
}

func Test_Hoverfly_isLocal_WhenLocalhostIP(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "127.0.0.1"}

	result := hoverfly.isLocal()

	Expect(result).To(BeTrue())
}

func Test_Hoverfly_isLocal_WhenAnotherDNS(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "specto.io"}

	result := hoverfly.isLocal()

	Expect(result).To(BeFalse())
}

func Test_Something(t *testing.T) {
	cfg := hoverfly.InitSettings()
	cfg.AdminPort = "8888"
	cfg.ProxyPort = "8500"
	cfg.Webserver = true
	cfg.Mode = "simulate"

	var b bytes.Buffer
	devNull := bufio.NewWriter(&b)
	logrus.SetOutput(devNull)

	hf := hoverfly.NewHoverflyWithConfiguration(cfg)
	hf.StartProxy()
	// adminApi := hoverfly.AdminApi{}
	// adminApi.StartAdminInterface(hf)

	err := hf.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("DELETE"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "hello",
						EncodedBody: false,
						Headers:     map[string][]string{},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	unit := Hoverfly{
		Host:      "localhost",
		AdminPort: "8500",
	}

	fmt.Println("error incoming")
	fmt.Println(err)

	Expect(unit.DeleteSimulations()).Should(Succeed())

	hf.StopProxy()
}
