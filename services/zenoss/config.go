package zenoss

import (
	"net/url"

	"github.com/influxdata/kapacitor/alert"
	"github.com/pkg/errors"
)

type SeverityMap struct {
	// OK level mapping to severity (number or string)
	OK interface{} `toml:"OK" json:"ok"`
	// Info level mapping to severity (number or string)
	Info interface{} `toml:"Info" json:"info"`
	// Warning level mapping to severity (number or string)
	Warning interface{} `toml:"Warning" json:"warning"`
	// Critical level mapping to severity (number or string)
	Critical interface{} `toml:"Critical" json:"critical"`
}

func (s *SeverityMap) ValueFor(level alert.Level) interface{} {
	var severity interface{}
	switch level {
	case alert.OK:
		severity = s.OK
	case alert.Info:
		severity = s.Info
	case alert.Warning:
		severity = s.Warning
	case alert.Critical:
		severity = s.Critical
	}
	return severity
}

type Config struct {
	// Whether Zenoss integration is enabled.
	Enabled bool `toml:"enabled" override:"enabled"`
	// Zenoss events API URL.
	URL string `toml:"url" override:"url"`
	// ServiceNow authentication username.
	Username string `toml:"username" override:"username"`
	// ServiceNow authentication password.
	Password string `toml:"password" override:"password,redact"`
	// Action (router name).
	Action string `toml:"action" override:"action"`
	// Router method.
	Method string `toml:"method" override:"method"`
	// Event type.
	Type string `toml:"type" override:"type"`
	// Event TID.
	TID int64 `toml:"tid" override:"tid"`
	// Level to severity map.
	SeverityMap SeverityMap `toml:"severity-map" override:"severity-map"`
	// Whether all alerts should automatically post to ServiceNow.
	Global bool `toml:"global" override:"global"`
	// Whether all alerts should automatically use stateChangesOnly mode.
	// Only applies if global is also set.
	StateChangesOnly bool `toml:"state-changes-only" override:"state-changes-only"`
}

func NewConfig() Config {
	return Config{
		URL:         "https://tenant.zenoss.io:8080/zport/dmd/evconsole_router",
		Action:      "EventsRouter",
		Method:      "add_event",
		Type:        "rpc",
		TID:         1,
		SeverityMap: SeverityMap{OK: "Clear", Info: "Info", Warning: "Warning", Critical: "Critical"},
	}
}

func (c Config) Validate() error {
	if c.Enabled && c.URL == "" {
		return errors.New("must specify events URL")
	}
	if _, err := url.Parse(c.URL); err != nil {
		return errors.Wrapf(err, "invalid url %q", c.URL)
	}
	if c.SeverityMap.OK == nil || c.SeverityMap.Info == nil || c.SeverityMap.Warning == nil || c.SeverityMap.Critical == nil {
		return errors.New("invalid severity mapping")
	}
	return nil
}
