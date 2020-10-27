package conch

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/joyent/kosh/logger"
)

type Client struct {
	Sling *sling.Sling
	logger.Logger
}

type Config interface {
	GetURL() string
	GetToken() string
	GetLogger() logger.Logger
}

var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func defaultUserAgent() string {
	f, _ := os.Executable()
	return fmt.Sprintf("go-conch %s", filepath.Base(f))
}

func New(c Config) *Client {
	c.GetLogger().Debug(c)
	s := sling.New().
		Client(&http.Client{Transport: defaultTransport}).
		Set("User-Agent", defaultUserAgent()).
		Base(c.GetURL()).
		Set("Authorization", "Bearer "+c.GetToken())

	return &Client{s, c.GetLogger()}
}

func (c *Client) New() *Client {
	s := c.Sling.New()
	l := c.Logger
	return &Client{s, l}
}

func (c *Client) UserAgent(ua string) *Client {
	c = c.New()
	c.Sling.Set("User-Agent", ua)
	return c
}

func (c *Client) Path(path string, ids ...string) *Client {
	c = c.New()
	c.Sling.Path(fmt.Sprintf("%s/", path))
	for _, id := range ids {
		if id != "" {
			c.Sling.Path(fmt.Sprintf("%s/", id))
		}
	}
	return c
}

func (c *Client) Token(id ...string) *Client {
	return c.Path("token", id...)
}

func (c *Client) Relay(id ...string) *Client {
	return c.Path("relay", id...)
}

func (c *Client) Register() *Client {
	return c.Path("register")
}

func (c *Client) HardwareVendor(id ...string) *Client {
	return c.Path("hardware_vendor", id...)
}

func (c *Client) HardwareProduct(id ...string) *Client {
	return c.Path("hardware_product", id...)
}

func (c *Client) Specification(path ...string) *Client {
	return c.Path("specification").
		stripTrailingSlash().
		Path(fmt.Sprintf("?path=%s", path))
}

func (c *Client) DC(id ...string) *Client {
	return c.Path("dc", id...)
}

func (c *Client) Rooms() *Client {
	return c.Path("rooms")
}

func (c *Client) User(id ...string) *Client {
	return c.Path("user", id...)
}

func (c *Client) Device(id ...string) *Client {
	return c.Path("device", id...)
}

func (c *Client) Settings(id ...string) *Client {
	return c.Path("settings", id...)
}

func (c *Client) Build(id ...string) *Client {
	return c.Path("build", id...)
}

func (c *Client) Organization(id ...string) *Client {
	return c.Path("organization", id...)
}

func (c *Client) Rack(id ...string) *Client {
	return c.Path("rack", id...)
}

func (c *Client) RackRole(id ...string) *Client {
	return c.Path("rack_role", id...)
}

func (c *Client) Room(id ...string) *Client {
	return c.Path("room", id...)
}

func (c *Client) Layout(id ...string) *Client {
	return c.Path("layout", id...)
}

func (c *Client) Validation(id ...string) *Client {
	return c.Path("validation", id...)
}

func (c *Client) ValidationPlan(id ...string) *Client {
	return c.Path("validation_plan", id...)
}

func (c *Client) ValidationState(id ...string) *Client {
	return c.Path("validation_state", id...)
}

func (c *Client) Interface(id ...string) *Client {
	return c.Path("interface", id...)
}

func (c *Client) Schema(name ...string) *Client {
	return c.Path("schema", name...)
}

func (c *Client) Field(ids ...string) *Client {
	c = c.New()
	for _, id := range ids {
		c.Sling.Path(fmt.Sprintf("%s/", id))
	}
	return c
}

func (c *Client) stripTrailingSlash() *Client {
	c = c.New()
	r, _ := c.Sling.Request()
	base := strings.TrimRight(r.URL.String(), "/")
	c.Sling.Base(base)
	return c
}

func (c *Client) WithParams(settings map[string]string) *Client {
	c = c.stripTrailingSlash()
	segments := []string{}
	for k, v := range settings {
		segments = append(segments, fmt.Sprintf("%s=%s&", url.PathEscape(k), url.PathEscape(v)))
	}
	qs := strings.Join(segments, "&")
	c.Sling.Path(fmt.Sprintf("?%s", qs))
	return c
}

func (c *Client) ValidationStates(states ...string) *Client {
	c = c.New()
	c.Sling.Path("validation_state") // TODO add query params from states
	return c
}

func (c *Client) SKU() *Client {
	return c.Path("sku")
}

func (c *Client) Assignment() *Client {
	return c.Path("assignment")
}

func (c *Client) Links() *Client {
	return c.Path("links")
}

func (c *Client) PXE() *Client {
	return c.Path("pxe")
}

func (c *Client) Location() *Client {
	return c.Path("location")
}

func (c *Client) Revoke() *Client {
	return c.Path("revoke")
}

func (c *Client) Password() *Client {
	return c.Path("password")
}

func (c *Client) AssetTag() *Client {
	return c.Path("asset_tag")
}

func (c *Client) Validated() *Client {
	return c.Path("validated")
}

func (c *Client) Phase() *Client {
	return c.Path("phase")
}

func (c *Client) DeviceReport(id ...string) *Client {
	return c.Path("device_report", id...)
}

func (c *Client) Post(data interface{}) *Client {
	c = c.New()
	c.Sling.Post("").BodyJSON(data)
	return c
}

func (c *Client) Put(data interface{}) *Client {
	c = c.New()
	c.Sling.Put("").BodyJSON(data)
	return c
}

func (c *Client) Delete(data ...interface{}) *Client {
	c = c.New()
	c.Sling.Delete("")
	for _, d := range data {
		c.Sling.BodyJSON(d)
	}
	return c
}

// when you don't expect a response
func (c *Client) Send() (*http.Response, error) {
	c.Debug("Send")
	req, err := c.Sling.Request()
	c.Info("URL: ", req.URL)
	c.Debug(req, err)

	res, err := c.Sling.Do(req, nil, nil)
	c.Debug(res, err)

	return res, err
}

// when you do expect a response
func (c *Client) Receive(data interface{}) (*http.Response, error) {
	c.Debug("Receive")
	req, err := c.Sling.Request()
	c.Info("URL: ", req.URL)
	c.Debug(req, err)

	res, err := c.Sling.ReceiveSuccess(data)
	c.Debug(res, err)

	return res, err
}
