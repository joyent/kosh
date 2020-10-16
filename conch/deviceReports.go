package conch

import (
	"encoding/json"
	"io"

	"github.com/joyent/kosh/conch/types"
)

// POST /device_report
func (c *Client) SendDeviceReport(r io.Reader) error {
	report := &types.DeviceReport{}
	json.NewDecoder(r).Decode(report)
	_, e := c.DeviceReport().Post(report).Send()
	return e
}

// POST /device_report?no_update_db=1
func (c *Client) ValidateDeviceReport(r io.Reader) (results types.ReportValidationResults) {
	report := &types.DeviceReport{}
	json.NewDecoder(r).Decode(report)
	c.DeviceReport().Post(report).Receive(&results)
	return
}

// GET /device_report/:device_report_id
func (c *Client) GetDeviceReport(id string) (report types.DeviceReportRow) {
	c.DeviceReport(id).Receive(&report)
	return
}
