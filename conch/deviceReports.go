package conch

import (
	"encoding/json"
	"io"

	"github.com/joyent/kosh/v3/conch/types"
)

// SendDeviceReport (POST /device_report) reads a new device report from an
// io.Reader and sends it to the API, it does not return the results
func (c *Client) SendDeviceReport(r io.Reader) error {
	report := &types.DeviceReport{}
	json.NewDecoder(r).Decode(report)
	_, e := c.DeviceReport().Post(report).Send()
	return e
}

// ValidateDeviceReport (POST /device_report?no_update_db=1) reads a new device
// report from an io.Reader and sends it to the API returning the validation
// results
func (c *Client) ValidateDeviceReport(r io.Reader) (results types.ReportValidationResults) {
	report := &types.DeviceReport{}
	json.NewDecoder(r).Decode(report)
	c.DeviceReport().Post(report).Receive(&results)
	return
}

// GetDeviceReport (GET /device_report/:device_report_id) returns the
// previously sent report for the given id string
func (c *Client) GetDeviceReport(id string) (report types.DeviceReportRow) {
	c.DeviceReport(id).Receive(&report)
	return
}
