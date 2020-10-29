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

// Client is a struct that represnts the current Conch client.
type Client struct {
	Sling *sling.Sling
	logger.Logger
}

// Config is an interface for an acceptable struct to configure the conch
// client.
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

// New takes a Config struct and returns a new instance of Client
func New(c Config) *Client {
	c.GetLogger().Debug(c)
	s := sling.New().
		Client(&http.Client{Transport: defaultTransport}).
		Set("User-Agent", defaultUserAgent()).
		Base(c.GetURL()).
		Set("Authorization", "Bearer "+c.GetToken())

	return &Client{s, c.GetLogger()}
}

// New performs a shallow clone of the current client and returns the
// new instance
func (c *Client) New() *Client {
	s := c.Sling.New()
	l := c.Logger
	return &Client{s, l}
}

// UserAgent sets the client's User-Agent header to the given string
func (c *Client) UserAgent(ua string) *Client {
	c = c.New()
	c.Sling.Set("User-Agent", ua)
	return c
}

// Path sets the sling path to the given path, optionally appending 0 or
// more IDs to the end of the path
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

// Token sets the last element in the pasth to token and optionally to the given
// token identifier
func (c *Client) Token(id ...string) *Client {
	return c.Path("token", id...)
}

// Relay sets the last element in the path to /relay and optionally appends the
// given relay identifier
func (c *Client) Relay(id ...string) *Client {
	return c.Path("relay", id...)
}

// Register sets the last element in the path to /register
func (c *Client) Register() *Client {
	return c.Path("register")
}

// HardwareVendor sets the last element in the path to /hardware_vendor and
// optionally appends teh given Hardare Vendor identifiers
func (c *Client) HardwareVendor(id ...string) *Client {
	return c.Path("hardware_vendor", id...)
}

// HardwareProduct sets the last element in the path to /hardware_product and
// optionally appends the given Hardare Product identifiers
func (c *Client) HardwareProduct(id ...string) *Client {
	return c.Path("hardware_product", id...)
}

// Specification sets the last element in the path to /specification and
// optionally sets the path query string
func (c *Client) Specification(path ...string) *Client {
	return c.Path("specification").
		stripTrailingSlash().
		Path(fmt.Sprintf("?path=%s", path))
}

// DC sets the last element in the path to /dc and
// appends the given identifiers
func (c *Client) DC(id ...string) *Client {
	return c.Path("dc", id...)
}

// Rooms sets the last element in the path to /rooms and
// appends the given identifiers
func (c *Client) Rooms() *Client {
	return c.Path("rooms")
}

// User sets the last element in the path to /user and
// appends the given identifiers
func (c *Client) User(id ...string) *Client {
	return c.Path("user", id...)
}

// Device sets the last element in the path to /device and
// appends the given identifiers
func (c *Client) Device(id ...string) *Client {
	return c.Path("device", id...)
}

// Settings sets the last element in the path to /settings and
// appends the given identifiers
func (c *Client) Settings(id ...string) *Client {
	return c.Path("settings", id...)
}

// Build sets the last element in the path to /build and
// appends the given identifiers
func (c *Client) Build(id ...string) *Client {
	return c.Path("build", id...)
}

// Organization sets the last element in the path to /organization and
// appends the given identifiers
func (c *Client) Organization(id ...string) *Client {
	return c.Path("organization", id...)
}

// Rack sets the last element in the path to /rack and
// appends the given identifiers
func (c *Client) Rack(id ...string) *Client {
	return c.Path("rack", id...)
}

// RackRole sets the last element in the path to /rack_role and
// appends the given identifiers
func (c *Client) RackRole(id ...string) *Client {
	return c.Path("rack_role", id...)
}

// Room sets the last element in the path to /room and
// appends the given identifiers
func (c *Client) Room(id ...string) *Client {
	return c.Path("room", id...)
}

// Layout sets the last element in the path to /layout and
// appends the given identifiers
func (c *Client) Layout(id ...string) *Client {
	return c.Path("layout", id...)
}

// Validation sets the last element in the path to /validation and
// appends the given identifiers
func (c *Client) Validation(id ...string) *Client {
	return c.Path("validation", id...)
}

// ValidationPlan sets the last element in the path to /validation_plan and
// appends the given identifiers
func (c *Client) ValidationPlan(id ...string) *Client {
	return c.Path("validation_plan", id...)
}

// ValidationState sets the last element in the path to /validation_state and
// appends the given identifiers
func (c *Client) ValidationState(id ...string) *Client {
	return c.Path("validation_state", id...)
}

// Interface sets the last element in the path to /interface and
// appends the given identifiers
func (c *Client) Interface(id ...string) *Client {
	return c.Path("interface", id...)
}

// Schema sets the last element in the path to /schema and
// appends the given identifiers
func (c *Client) Schema(name ...string) *Client {
	return c.Path("schema", name...)
}

// Field appends the given identifiers
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

// WithParams sets the query arguments to the given  key values in the params map
func (c *Client) WithParams(params map[string]string) *Client {
	c = c.stripTrailingSlash()
	segments := []string{}
	for k, v := range params {
		segments = append(segments, fmt.Sprintf("%s=%s&", url.PathEscape(k), url.PathEscape(v)))
	}
	qs := strings.Join(segments, "&")
	c.Sling.Path(fmt.Sprintf("?%s", qs))
	return c
}

// ValidationStates sets the last element in the path to /validation_state
// BUG(perigrin) Doesn't support query parameters yet
func (c *Client) ValidationStates(states ...string) *Client {
	c = c.New()
	c.Sling.Path("validation_state") // TODO add query params from states
	return c
}

// SKU sets the last element in the path to /sku
func (c *Client) SKU() *Client {
	return c.Path("sku")
}

// Assignment sets the last element in the path to /assignment
func (c *Client) Assignment() *Client {
	return c.Path("assignment")
}

// Links sets the last element in the path to /links
func (c *Client) Links() *Client {
	return c.Path("links")
}

// PXE sets the last element in the path to /pxe
func (c *Client) PXE() *Client {
	return c.Path("pxe")
}

// Location sets the last element in the path to /location
func (c *Client) Location() *Client {
	return c.Path("location")
}

// Revoke sets the last element in the path to /revoke
func (c *Client) Revoke() *Client {
	return c.Path("revoke")
}

// Password sets the last element in the path to /password
func (c *Client) Password() *Client {
	return c.Path("password")
}

// AssetTag sets the last element in the path to /asset_tag
func (c *Client) AssetTag() *Client {
	return c.Path("asset_tag")
}

// Validated sets the last element in the path to /validated
func (c *Client) Validated() *Client {
	return c.Path("validated")
}

// Phase sets the last element in the path to /phase
func (c *Client) Phase() *Client {
	return c.Path("phase")
}

// DeviceReport sets the last element in the path to /device_report and appends
// the given identifier
func (c *Client) DeviceReport(id ...string) *Client {
	return c.Path("device_report", id...)
}

// Post sets the HTTP method to POST and sets the JSON body to the given data
func (c *Client) Post(data interface{}) *Client {
	c = c.New()
	c.Sling.Post("").BodyJSON(data)
	return c
}

// Put sets the HTTP method to PUT and sets the JSON body to the given data
func (c *Client) Put(data interface{}) *Client {
	c = c.New()
	c.Sling.Put("").BodyJSON(data)
	return c
}

// Delete sets the HTTP method to DELETE and optionally sets the JSON body to the given data
func (c *Client) Delete(data ...interface{}) *Client {
	c = c.New()
	c.Sling.Delete("")
	for _, d := range data {
		c.Sling.BodyJSON(d)
	}
	return c
}

// Send sends a HTTP request to the API server  without expecting a return data
// structure. It returns the *http.Response and/or error from the request.
func (c *Client) Send() (*http.Response, error) {
	c.Debug("Send")
	req, err := c.Sling.Request()
	c.Info("URL: ", req.URL)
	c.Debug(req, err)

	res, err := c.Sling.Do(req, nil, nil)
	c.Debug(res, err)

	return res, err
}

// Receive sends a HTTP request to the API server and decodes the results into
// the provided structure structure. It returns the *http.Response and/or error
// from the request.
func (c *Client) Receive(data interface{}) (*http.Response, error) {
	c.Debug("Receive")
	req, err := c.Sling.Request()
	c.Info("URL: ", req.URL)
	c.Debug(req, err)

	res, err := c.Sling.ReceiveSuccess(data)
	c.Debug(res, err)

	return res, err
}
