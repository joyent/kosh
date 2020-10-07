package conch

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dghubble/sling"
)

type Logger interface {
	Debug(...interface{})
}

type Client struct {
	Sling *sling.Sling
	Logger
}

func New(api, token string, l Logger) *Client {
	s := sling.New().Base(api).Set("Authorization", "Bearer "+token)
	return &Client{s, l}
}

func (c *Client) New() *Client {
	s := c.Sling.New()
	l := c.Logger
	return &Client{s, l}
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

func (c *Client) Delete() *Client {
	c = c.New()
	c.Sling.Delete("")
	return c
}

// when you don't expect a response
func (c *Client) Send() (*http.Response, error) {
	req, err := c.Sling.Request()
	if err != nil {
		return nil, err
	}

	c.Debug("Sending request: ", req.URL)
	res, err := c.Sling.Do(req, nil, nil)
	if err != nil {
		c.Debug("Error: %+v", err)
		return nil, err
	}
	return res, nil
}

// when you do expect a response
func (c *Client) Receive(data interface{}) (interface{}, error) {
	req, err := c.Sling.Request()
	if err != nil {
		return nil, err
	}

	c.Debug("Sending request: ", req.URL)
	_, err = c.Sling.ReceiveSuccess(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
