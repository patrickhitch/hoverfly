package functional_tests

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"io"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/dghubble/sling"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

func DoRequest(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)
	Expect(err).To(BeNil())
	return response
}

type Hoverfly struct {
	adminPort int
	adminUrl  string
	proxyPort int
	proxyUrl  string
	process   *exec.Cmd
	commands  []string
}

func NewHoverfly() *Hoverfly {
	return &Hoverfly{
		adminPort: freeport.GetPort(),
		proxyPort: freeport.GetPort(),
	}
}

func (this *Hoverfly) Start(commands ...string) {
	this.process = this.startHoverflyInternal(this.adminPort, this.proxyPort, commands...)
	this.adminUrl = fmt.Sprintf("http://localhost:%v", this.adminPort)
	this.proxyUrl = fmt.Sprintf("http://localhost:%v", this.proxyPort)
}

func (this Hoverfly) Stop() {
	this.process.Process.Kill()
}

func (this Hoverfly) DeleteBoltDb() {
	workingDirectory, _ := os.Getwd()
	os.Remove(workingDirectory + "requests.db")
}

func (this Hoverfly) GetMode() string {
	currentState := &v2.ModeView{}
	resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/mode", this.adminPort)))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	json.Unmarshal(body, currentState)

	return currentState.Mode
}

func (this Hoverfly) SetMode(mode string) {
	this.SetModeWithArgs(mode, v2.ModeArgumentsView{})
}

func (this Hoverfly) SetModeWithArgs(mode string, arguments v2.ModeArgumentsView) {
	newMode := &v2.ModeView{
		Mode:      mode,
		Arguments: arguments,
	}

	DoRequest(sling.New().Put(this.adminUrl + "/api/v2/hoverfly/mode").BodyJSON(newMode))
}

func (this Hoverfly) SetMiddleware(binary, script string) {
	newMiddleware := v2.MiddlewareView{
		Binary: binary,
		Script: script,
	}

	DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", this.adminPort)).BodyJSON(newMiddleware))
}

func (this Hoverfly) GetSimulation() io.Reader {
	res := sling.New().Get(this.adminUrl + "/api/v2/simulation")
	req := DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func (this Hoverfly) ImportSimulation(simulation string) {
	req := sling.New().Put(this.adminUrl + "/api/v2/simulation").Body(bytes.NewBufferString(simulation))
	response := DoRequest(req)
	Expect(response.StatusCode).To(Equal(http.StatusOK))
}

func (this Hoverfly) ExportSimulation() v2.SimulationViewV2 {
	reader := this.GetSimulation()
	simulationBytes, err := ioutil.ReadAll(reader)
	Expect(err).To(BeNil())

	var simulation v2.SimulationViewV2

	err = json.Unmarshal(simulationBytes, &simulation)
	Expect(err).To(BeNil())

	return simulation
}

func (this Hoverfly) GetCache() v2.CacheView {
	req := sling.New().Get(this.adminUrl + "/api/v2/cache")
	response := DoRequest(req)
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	cacheBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	var cache v2.CacheView

	err = json.Unmarshal(cacheBytes, &cache)
	Expect(err).To(BeNil())

	return cache
}

func (this Hoverfly) FlushCache() v2.CacheView {
	req := sling.New().Delete(this.adminUrl + "/api/v2/cache")
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))

	cacheBytes, err := ioutil.ReadAll(res.Body)
	Expect(err).To(BeNil())

	var cache v2.CacheView

	err = json.Unmarshal(cacheBytes, &cache)
	Expect(err).To(BeNil())

	return cache
}

func (this Hoverfly) Proxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse(this.proxyUrl)
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func (this Hoverfly) ProxyWithAuth(r *sling.Sling, user, password string) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse(fmt.Sprintf("http://%s:%s@localhost:%v", user, password, this.proxyPort))
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func (this Hoverfly) GetAdminPort() string {
	return strconv.Itoa(this.adminPort)
}

func (this Hoverfly) GetProxyPort() string {
	return strconv.Itoa(this.proxyPort)
}

func (this Hoverfly) GetPid() int {
	return this.process.Process.Pid
}

func (this Hoverfly) startHoverflyInternal(adminPort, proxyPort int, additionalCommands ...string) *exec.Cmd {
	hoverflyBinaryUri := BuildBinaryPath()

	commands := []string{
		"-ap",
		strconv.Itoa(adminPort),
		"-pp",
		strconv.Itoa(proxyPort),
	}

	commands = append(commands, additionalCommands...)
	this.commands = commands
	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)
	err := hoverflyCmd.Start()

	BinaryErrorCheck(err, hoverflyBinaryUri)

	for _, command := range commands {
		if command == "-add" {
			time.Sleep(time.Second * 3)
			return hoverflyCmd
		}
	}

	this.healthcheck()

	return hoverflyCmd
}

func BuildBinaryPath() string {
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, "bin/hoverfly")
}

func BinaryErrorCheck(err error, binaryPath string) {
	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(binaryPath)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}
}

func (this Hoverfly) healthcheck() {
	Eventually(func() int {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/health", this.adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second*5).Should(BeNumerically("==", http.StatusOK), "Hoverfly not running on %d", this.adminPort, this.commands)
}

func Healthcheck(adminPort int) {
	var err error
	var resp *http.Response

	Eventually(func() int {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/api/health", adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second*5).Should(BeNumerically("==", http.StatusOK), "Hoverfly not running on %d but have no extra information", adminPort)
}

func Run(binary string, commands ...string) string {
	cmd := exec.Command(binary, commands...)
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			return string(exitError.Stderr)
		}
	}

	return strings.TrimSpace(string(out))
}

func GenerateFileName() string {

	rb := make([]byte, 6)
	rand.Read(rb)

	rs := base64.URLEncoding.EncodeToString(rb)

	return "testdata-gen/" + rs + ".json"
}

func TableToSliceMapStringString(table string) map[string]map[string]string {
	results := map[string]map[string]string{}

	tableRows := strings.Split(table, "\n")
	headings := []string{}

	for _, heading := range strings.Split(tableRows[1], "|") {
		headings = append(headings, strings.TrimSpace(heading))
	}

	for _, row := range tableRows[2:] {
		if !strings.Contains(row, "-+-") {
			rowValues := strings.Split(row, "|")

			result := map[string]string{}
			for i, value := range rowValues {
				result[headings[i]] = strings.TrimSpace(value)
			}

			results[result["TARGET NAME"]] = result
		}
	}

	return results
}
