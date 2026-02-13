package radiko

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/yyoshiki41/go-radiko/internal/m3u8"
	"github.com/yyoshiki41/go-radiko/internal/util"
)

const (
	timeshiftEndpoint = "https://tf-f-rpaa-radiko.smartstream.ne.jp/tf/playlist.m3u8"

	radikoAreaIDHeader = "X-Radiko-AreaId"
)

// TimeshiftPlaylistM3U8 returns uri.
func (c *Client) TimeshiftPlaylistM3U8(ctx context.Context, stationID string, start time.Time) (string, error) {
	prog, err := c.GetProgramByStartTime(ctx, stationID, start)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(timeshiftEndpoint)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("station_id", stationID)
	q.Set("start_at", prog.Ft)
	q.Set("ft", prog.Ft)
	q.Set("end_at", prog.To)
	q.Set("to", prog.To)
	q.Set("l", "15")
	q.Set("lsid", generateLsid())
	q.Set("type", "b")
	q.Set("preroll", "0")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(radikoAuthTokenHeader, c.AuthToken())
	req.Header.Set(radikoAreaIDHeader, c.AreaID())

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return m3u8.GetURI(resp.Body)
}

// GetTimeshiftURL returns a timeshift url for web browser.
func GetTimeshiftURL(stationID string, start time.Time) string {
	endpoint := path.Join("#!/ts", stationID, util.Datetime(start))
	return defaultEndpoint + "/" + endpoint
}

// generateLsid generates a random 32-character hex string for the listener session ID.
func generateLsid() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
