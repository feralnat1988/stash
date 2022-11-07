// Package network provides the Chrome DevTools Protocol
// commands, types, and events for the Network domain.
//
// Network domain allows tracking network activities of the page. It exposes
// information about http, file, data and other requests and responses, their
// headers, bodies, timing, etc.
//
// Generated by the cdproto-gen command.
package network

// Code generated by cdproto-gen. DO NOT EDIT.

import (
	"context"
	"encoding/base64"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/debugger"
	"github.com/chromedp/cdproto/io"
)

// SetAcceptedEncodingsParams sets a list of content encodings that will be
// accepted. Empty list means no encoding is accepted.
type SetAcceptedEncodingsParams struct {
	Encodings []ContentEncoding `json:"encodings"` // List of accepted content encodings.
}

// SetAcceptedEncodings sets a list of content encodings that will be
// accepted. Empty list means no encoding is accepted.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setAcceptedEncodings
//
// parameters:
//   encodings - List of accepted content encodings.
func SetAcceptedEncodings(encodings []ContentEncoding) *SetAcceptedEncodingsParams {
	return &SetAcceptedEncodingsParams{
		Encodings: encodings,
	}
}

// Do executes Network.setAcceptedEncodings against the provided context.
func (p *SetAcceptedEncodingsParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetAcceptedEncodings, p, nil)
}

// ClearAcceptedEncodingsOverrideParams clears accepted encodings set by
// setAcceptedEncodings.
type ClearAcceptedEncodingsOverrideParams struct{}

// ClearAcceptedEncodingsOverride clears accepted encodings set by
// setAcceptedEncodings.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-clearAcceptedEncodingsOverride
func ClearAcceptedEncodingsOverride() *ClearAcceptedEncodingsOverrideParams {
	return &ClearAcceptedEncodingsOverrideParams{}
}

// Do executes Network.clearAcceptedEncodingsOverride against the provided context.
func (p *ClearAcceptedEncodingsOverrideParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandClearAcceptedEncodingsOverride, nil, nil)
}

// ClearBrowserCacheParams clears browser cache.
type ClearBrowserCacheParams struct{}

// ClearBrowserCache clears browser cache.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-clearBrowserCache
func ClearBrowserCache() *ClearBrowserCacheParams {
	return &ClearBrowserCacheParams{}
}

// Do executes Network.clearBrowserCache against the provided context.
func (p *ClearBrowserCacheParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandClearBrowserCache, nil, nil)
}

// ClearBrowserCookiesParams clears browser cookies.
type ClearBrowserCookiesParams struct{}

// ClearBrowserCookies clears browser cookies.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-clearBrowserCookies
func ClearBrowserCookies() *ClearBrowserCookiesParams {
	return &ClearBrowserCookiesParams{}
}

// Do executes Network.clearBrowserCookies against the provided context.
func (p *ClearBrowserCookiesParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandClearBrowserCookies, nil, nil)
}

// DeleteCookiesParams deletes browser cookies with matching name and url or
// domain/path pair.
type DeleteCookiesParams struct {
	Name   string `json:"name"`             // Name of the cookies to remove.
	URL    string `json:"url,omitempty"`    // If specified, deletes all the cookies with the given name where domain and path match provided URL.
	Domain string `json:"domain,omitempty"` // If specified, deletes only cookies with the exact domain.
	Path   string `json:"path,omitempty"`   // If specified, deletes only cookies with the exact path.
}

// DeleteCookies deletes browser cookies with matching name and url or
// domain/path pair.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-deleteCookies
//
// parameters:
//   name - Name of the cookies to remove.
func DeleteCookies(name string) *DeleteCookiesParams {
	return &DeleteCookiesParams{
		Name: name,
	}
}

// WithURL if specified, deletes all the cookies with the given name where
// domain and path match provided URL.
func (p DeleteCookiesParams) WithURL(url string) *DeleteCookiesParams {
	p.URL = url
	return &p
}

// WithDomain if specified, deletes only cookies with the exact domain.
func (p DeleteCookiesParams) WithDomain(domain string) *DeleteCookiesParams {
	p.Domain = domain
	return &p
}

// WithPath if specified, deletes only cookies with the exact path.
func (p DeleteCookiesParams) WithPath(path string) *DeleteCookiesParams {
	p.Path = path
	return &p
}

// Do executes Network.deleteCookies against the provided context.
func (p *DeleteCookiesParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandDeleteCookies, p, nil)
}

// DisableParams disables network tracking, prevents network events from
// being sent to the client.
type DisableParams struct{}

// Disable disables network tracking, prevents network events from being sent
// to the client.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-disable
func Disable() *DisableParams {
	return &DisableParams{}
}

// Do executes Network.disable against the provided context.
func (p *DisableParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandDisable, nil, nil)
}

// EmulateNetworkConditionsParams activates emulation of network conditions.
type EmulateNetworkConditionsParams struct {
	Offline            bool           `json:"offline"`                  // True to emulate internet disconnection.
	Latency            float64        `json:"latency"`                  // Minimum latency from request sent to response headers received (ms).
	DownloadThroughput float64        `json:"downloadThroughput"`       // Maximal aggregated download throughput (bytes/sec). -1 disables download throttling.
	UploadThroughput   float64        `json:"uploadThroughput"`         // Maximal aggregated upload throughput (bytes/sec).  -1 disables upload throttling.
	ConnectionType     ConnectionType `json:"connectionType,omitempty"` // Connection type if known.
}

// EmulateNetworkConditions activates emulation of network conditions.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-emulateNetworkConditions
//
// parameters:
//   offline - True to emulate internet disconnection.
//   latency - Minimum latency from request sent to response headers received (ms).
//   downloadThroughput - Maximal aggregated download throughput (bytes/sec). -1 disables download throttling.
//   uploadThroughput - Maximal aggregated upload throughput (bytes/sec).  -1 disables upload throttling.
func EmulateNetworkConditions(offline bool, latency float64, downloadThroughput float64, uploadThroughput float64) *EmulateNetworkConditionsParams {
	return &EmulateNetworkConditionsParams{
		Offline:            offline,
		Latency:            latency,
		DownloadThroughput: downloadThroughput,
		UploadThroughput:   uploadThroughput,
	}
}

// WithConnectionType connection type if known.
func (p EmulateNetworkConditionsParams) WithConnectionType(connectionType ConnectionType) *EmulateNetworkConditionsParams {
	p.ConnectionType = connectionType
	return &p
}

// Do executes Network.emulateNetworkConditions against the provided context.
func (p *EmulateNetworkConditionsParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandEmulateNetworkConditions, p, nil)
}

// EnableParams enables network tracking, network events will now be
// delivered to the client.
type EnableParams struct {
	MaxTotalBufferSize    int64 `json:"maxTotalBufferSize,omitempty"`    // Buffer size in bytes to use when preserving network payloads (XHRs, etc).
	MaxResourceBufferSize int64 `json:"maxResourceBufferSize,omitempty"` // Per-resource buffer size in bytes to use when preserving network payloads (XHRs, etc).
	MaxPostDataSize       int64 `json:"maxPostDataSize,omitempty"`       // Longest post body size (in bytes) that would be included in requestWillBeSent notification
}

// Enable enables network tracking, network events will now be delivered to
// the client.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-enable
//
// parameters:
func Enable() *EnableParams {
	return &EnableParams{}
}

// WithMaxTotalBufferSize buffer size in bytes to use when preserving network
// payloads (XHRs, etc).
func (p EnableParams) WithMaxTotalBufferSize(maxTotalBufferSize int64) *EnableParams {
	p.MaxTotalBufferSize = maxTotalBufferSize
	return &p
}

// WithMaxResourceBufferSize per-resource buffer size in bytes to use when
// preserving network payloads (XHRs, etc).
func (p EnableParams) WithMaxResourceBufferSize(maxResourceBufferSize int64) *EnableParams {
	p.MaxResourceBufferSize = maxResourceBufferSize
	return &p
}

// WithMaxPostDataSize longest post body size (in bytes) that would be
// included in requestWillBeSent notification.
func (p EnableParams) WithMaxPostDataSize(maxPostDataSize int64) *EnableParams {
	p.MaxPostDataSize = maxPostDataSize
	return &p
}

// Do executes Network.enable against the provided context.
func (p *EnableParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandEnable, p, nil)
}

// GetAllCookiesParams returns all browser cookies. Depending on the backend
// support, will return detailed cookie information in the cookies field.
type GetAllCookiesParams struct{}

// GetAllCookies returns all browser cookies. Depending on the backend
// support, will return detailed cookie information in the cookies field.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getAllCookies
func GetAllCookies() *GetAllCookiesParams {
	return &GetAllCookiesParams{}
}

// GetAllCookiesReturns return values.
type GetAllCookiesReturns struct {
	Cookies []*Cookie `json:"cookies,omitempty"` // Array of cookie objects.
}

// Do executes Network.getAllCookies against the provided context.
//
// returns:
//   cookies - Array of cookie objects.
func (p *GetAllCookiesParams) Do(ctx context.Context) (cookies []*Cookie, err error) {
	// execute
	var res GetAllCookiesReturns
	err = cdp.Execute(ctx, CommandGetAllCookies, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Cookies, nil
}

// GetCertificateParams returns the DER-encoded certificate.
type GetCertificateParams struct {
	Origin string `json:"origin"` // Origin to get certificate for.
}

// GetCertificate returns the DER-encoded certificate.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getCertificate
//
// parameters:
//   origin - Origin to get certificate for.
func GetCertificate(origin string) *GetCertificateParams {
	return &GetCertificateParams{
		Origin: origin,
	}
}

// GetCertificateReturns return values.
type GetCertificateReturns struct {
	TableNames []string `json:"tableNames,omitempty"`
}

// Do executes Network.getCertificate against the provided context.
//
// returns:
//   tableNames
func (p *GetCertificateParams) Do(ctx context.Context) (tableNames []string, err error) {
	// execute
	var res GetCertificateReturns
	err = cdp.Execute(ctx, CommandGetCertificate, p, &res)
	if err != nil {
		return nil, err
	}

	return res.TableNames, nil
}

// GetCookiesParams returns all browser cookies for the current URL.
// Depending on the backend support, will return detailed cookie information in
// the cookies field.
type GetCookiesParams struct {
	Urls []string `json:"urls,omitempty"` // The list of URLs for which applicable cookies will be fetched. If not specified, it's assumed to be set to the list containing the URLs of the page and all of its subframes.
}

// GetCookies returns all browser cookies for the current URL. Depending on
// the backend support, will return detailed cookie information in the cookies
// field.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getCookies
//
// parameters:
func GetCookies() *GetCookiesParams {
	return &GetCookiesParams{}
}

// WithUrls the list of URLs for which applicable cookies will be fetched. If
// not specified, it's assumed to be set to the list containing the URLs of the
// page and all of its subframes.
func (p GetCookiesParams) WithUrls(urls []string) *GetCookiesParams {
	p.Urls = urls
	return &p
}

// GetCookiesReturns return values.
type GetCookiesReturns struct {
	Cookies []*Cookie `json:"cookies,omitempty"` // Array of cookie objects.
}

// Do executes Network.getCookies against the provided context.
//
// returns:
//   cookies - Array of cookie objects.
func (p *GetCookiesParams) Do(ctx context.Context) (cookies []*Cookie, err error) {
	// execute
	var res GetCookiesReturns
	err = cdp.Execute(ctx, CommandGetCookies, p, &res)
	if err != nil {
		return nil, err
	}

	return res.Cookies, nil
}

// GetResponseBodyParams returns content served for the given request.
type GetResponseBodyParams struct {
	RequestID RequestID `json:"requestId"` // Identifier of the network request to get content for.
}

// GetResponseBody returns content served for the given request.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getResponseBody
//
// parameters:
//   requestID - Identifier of the network request to get content for.
func GetResponseBody(requestID RequestID) *GetResponseBodyParams {
	return &GetResponseBodyParams{
		RequestID: requestID,
	}
}

// GetResponseBodyReturns return values.
type GetResponseBodyReturns struct {
	Body          string `json:"body,omitempty"`          // Response body.
	Base64encoded bool   `json:"base64Encoded,omitempty"` // True, if content was sent as base64.
}

// Do executes Network.getResponseBody against the provided context.
//
// returns:
//   body - Response body.
func (p *GetResponseBodyParams) Do(ctx context.Context) (body []byte, err error) {
	// execute
	var res GetResponseBodyReturns
	err = cdp.Execute(ctx, CommandGetResponseBody, p, &res)
	if err != nil {
		return nil, err
	}

	// decode
	var dec []byte
	if res.Base64encoded {
		dec, err = base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		dec = []byte(res.Body)
	}
	return dec, nil
}

// GetRequestPostDataParams returns post data sent with the request. Returns
// an error when no data was sent with the request.
type GetRequestPostDataParams struct {
	RequestID RequestID `json:"requestId"` // Identifier of the network request to get content for.
}

// GetRequestPostData returns post data sent with the request. Returns an
// error when no data was sent with the request.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getRequestPostData
//
// parameters:
//   requestID - Identifier of the network request to get content for.
func GetRequestPostData(requestID RequestID) *GetRequestPostDataParams {
	return &GetRequestPostDataParams{
		RequestID: requestID,
	}
}

// GetRequestPostDataReturns return values.
type GetRequestPostDataReturns struct {
	PostData string `json:"postData,omitempty"` // Request body string, omitting files from multipart requests
}

// Do executes Network.getRequestPostData against the provided context.
//
// returns:
//   postData - Request body string, omitting files from multipart requests
func (p *GetRequestPostDataParams) Do(ctx context.Context) (postData string, err error) {
	// execute
	var res GetRequestPostDataReturns
	err = cdp.Execute(ctx, CommandGetRequestPostData, p, &res)
	if err != nil {
		return "", err
	}

	return res.PostData, nil
}

// GetResponseBodyForInterceptionParams returns content served for the given
// currently intercepted request.
type GetResponseBodyForInterceptionParams struct {
	InterceptionID InterceptionID `json:"interceptionId"` // Identifier for the intercepted request to get body for.
}

// GetResponseBodyForInterception returns content served for the given
// currently intercepted request.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getResponseBodyForInterception
//
// parameters:
//   interceptionID - Identifier for the intercepted request to get body for.
func GetResponseBodyForInterception(interceptionID InterceptionID) *GetResponseBodyForInterceptionParams {
	return &GetResponseBodyForInterceptionParams{
		InterceptionID: interceptionID,
	}
}

// GetResponseBodyForInterceptionReturns return values.
type GetResponseBodyForInterceptionReturns struct {
	Body          string `json:"body,omitempty"`          // Response body.
	Base64encoded bool   `json:"base64Encoded,omitempty"` // True, if content was sent as base64.
}

// Do executes Network.getResponseBodyForInterception against the provided context.
//
// returns:
//   body - Response body.
func (p *GetResponseBodyForInterceptionParams) Do(ctx context.Context) (body []byte, err error) {
	// execute
	var res GetResponseBodyForInterceptionReturns
	err = cdp.Execute(ctx, CommandGetResponseBodyForInterception, p, &res)
	if err != nil {
		return nil, err
	}

	// decode
	var dec []byte
	if res.Base64encoded {
		dec, err = base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		dec = []byte(res.Body)
	}
	return dec, nil
}

// TakeResponseBodyForInterceptionAsStreamParams returns a handle to the
// stream representing the response body. Note that after this command, the
// intercepted request can't be continued as is -- you either need to cancel it
// or to provide the response body. The stream only supports sequential read,
// IO.read will fail if the position is specified.
type TakeResponseBodyForInterceptionAsStreamParams struct {
	InterceptionID InterceptionID `json:"interceptionId"`
}

// TakeResponseBodyForInterceptionAsStream returns a handle to the stream
// representing the response body. Note that after this command, the intercepted
// request can't be continued as is -- you either need to cancel it or to
// provide the response body. The stream only supports sequential read, IO.read
// will fail if the position is specified.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-takeResponseBodyForInterceptionAsStream
//
// parameters:
//   interceptionID
func TakeResponseBodyForInterceptionAsStream(interceptionID InterceptionID) *TakeResponseBodyForInterceptionAsStreamParams {
	return &TakeResponseBodyForInterceptionAsStreamParams{
		InterceptionID: interceptionID,
	}
}

// TakeResponseBodyForInterceptionAsStreamReturns return values.
type TakeResponseBodyForInterceptionAsStreamReturns struct {
	Stream io.StreamHandle `json:"stream,omitempty"`
}

// Do executes Network.takeResponseBodyForInterceptionAsStream against the provided context.
//
// returns:
//   stream
func (p *TakeResponseBodyForInterceptionAsStreamParams) Do(ctx context.Context) (stream io.StreamHandle, err error) {
	// execute
	var res TakeResponseBodyForInterceptionAsStreamReturns
	err = cdp.Execute(ctx, CommandTakeResponseBodyForInterceptionAsStream, p, &res)
	if err != nil {
		return "", err
	}

	return res.Stream, nil
}

// ReplayXHRParams this method sends a new XMLHttpRequest which is identical
// to the original one. The following parameters should be identical: method,
// url, async, request body, extra headers, withCredentials attribute, user,
// password.
type ReplayXHRParams struct {
	RequestID RequestID `json:"requestId"` // Identifier of XHR to replay.
}

// ReplayXHR this method sends a new XMLHttpRequest which is identical to the
// original one. The following parameters should be identical: method, url,
// async, request body, extra headers, withCredentials attribute, user,
// password.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-replayXHR
//
// parameters:
//   requestID - Identifier of XHR to replay.
func ReplayXHR(requestID RequestID) *ReplayXHRParams {
	return &ReplayXHRParams{
		RequestID: requestID,
	}
}

// Do executes Network.replayXHR against the provided context.
func (p *ReplayXHRParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandReplayXHR, p, nil)
}

// SearchInResponseBodyParams searches for given string in response content.
type SearchInResponseBodyParams struct {
	RequestID     RequestID `json:"requestId"`               // Identifier of the network response to search.
	Query         string    `json:"query"`                   // String to search for.
	CaseSensitive bool      `json:"caseSensitive,omitempty"` // If true, search is case sensitive.
	IsRegex       bool      `json:"isRegex,omitempty"`       // If true, treats string parameter as regex.
}

// SearchInResponseBody searches for given string in response content.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-searchInResponseBody
//
// parameters:
//   requestID - Identifier of the network response to search.
//   query - String to search for.
func SearchInResponseBody(requestID RequestID, query string) *SearchInResponseBodyParams {
	return &SearchInResponseBodyParams{
		RequestID: requestID,
		Query:     query,
	}
}

// WithCaseSensitive if true, search is case sensitive.
func (p SearchInResponseBodyParams) WithCaseSensitive(caseSensitive bool) *SearchInResponseBodyParams {
	p.CaseSensitive = caseSensitive
	return &p
}

// WithIsRegex if true, treats string parameter as regex.
func (p SearchInResponseBodyParams) WithIsRegex(isRegex bool) *SearchInResponseBodyParams {
	p.IsRegex = isRegex
	return &p
}

// SearchInResponseBodyReturns return values.
type SearchInResponseBodyReturns struct {
	Result []*debugger.SearchMatch `json:"result,omitempty"` // List of search matches.
}

// Do executes Network.searchInResponseBody against the provided context.
//
// returns:
//   result - List of search matches.
func (p *SearchInResponseBodyParams) Do(ctx context.Context) (result []*debugger.SearchMatch, err error) {
	// execute
	var res SearchInResponseBodyReturns
	err = cdp.Execute(ctx, CommandSearchInResponseBody, p, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

// SetBlockedURLSParams blocks URLs from loading.
type SetBlockedURLSParams struct {
	Urls []string `json:"urls"` // URL patterns to block. Wildcards ('*') are allowed.
}

// SetBlockedURLS blocks URLs from loading.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setBlockedURLs
//
// parameters:
//   urls - URL patterns to block. Wildcards ('*') are allowed.
func SetBlockedURLS(urls []string) *SetBlockedURLSParams {
	return &SetBlockedURLSParams{
		Urls: urls,
	}
}

// Do executes Network.setBlockedURLs against the provided context.
func (p *SetBlockedURLSParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetBlockedURLS, p, nil)
}

// SetBypassServiceWorkerParams toggles ignoring of service worker for each
// request.
type SetBypassServiceWorkerParams struct {
	Bypass bool `json:"bypass"` // Bypass service worker and load from network.
}

// SetBypassServiceWorker toggles ignoring of service worker for each
// request.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setBypassServiceWorker
//
// parameters:
//   bypass - Bypass service worker and load from network.
func SetBypassServiceWorker(bypass bool) *SetBypassServiceWorkerParams {
	return &SetBypassServiceWorkerParams{
		Bypass: bypass,
	}
}

// Do executes Network.setBypassServiceWorker against the provided context.
func (p *SetBypassServiceWorkerParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetBypassServiceWorker, p, nil)
}

// SetCacheDisabledParams toggles ignoring cache for each request. If true,
// cache will not be used.
type SetCacheDisabledParams struct {
	CacheDisabled bool `json:"cacheDisabled"` // Cache disabled state.
}

// SetCacheDisabled toggles ignoring cache for each request. If true, cache
// will not be used.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setCacheDisabled
//
// parameters:
//   cacheDisabled - Cache disabled state.
func SetCacheDisabled(cacheDisabled bool) *SetCacheDisabledParams {
	return &SetCacheDisabledParams{
		CacheDisabled: cacheDisabled,
	}
}

// Do executes Network.setCacheDisabled against the provided context.
func (p *SetCacheDisabledParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetCacheDisabled, p, nil)
}

// SetCookieParams sets a cookie with the given cookie data; may overwrite
// equivalent cookies if they exist.
type SetCookieParams struct {
	Name         string              `json:"name"`                   // Cookie name.
	Value        string              `json:"value"`                  // Cookie value.
	URL          string              `json:"url,omitempty"`          // The request-URI to associate with the setting of the cookie. This value can affect the default domain, path, source port, and source scheme values of the created cookie.
	Domain       string              `json:"domain,omitempty"`       // Cookie domain.
	Path         string              `json:"path,omitempty"`         // Cookie path.
	Secure       bool                `json:"secure,omitempty"`       // True if cookie is secure.
	HTTPOnly     bool                `json:"httpOnly,omitempty"`     // True if cookie is http-only.
	SameSite     CookieSameSite      `json:"sameSite,omitempty"`     // Cookie SameSite type.
	Expires      *cdp.TimeSinceEpoch `json:"expires,omitempty"`      // Cookie expiration date, session cookie if not set
	Priority     CookiePriority      `json:"priority,omitempty"`     // Cookie Priority type.
	SameParty    bool                `json:"sameParty,omitempty"`    // True if cookie is SameParty.
	SourceScheme CookieSourceScheme  `json:"sourceScheme,omitempty"` // Cookie source scheme type.
	SourcePort   int64               `json:"sourcePort,omitempty"`   // Cookie source port. Valid values are {-1, [1, 65535]}, -1 indicates an unspecified port. An unspecified port value allows protocol clients to emulate legacy cookie scope for the port. This is a temporary ability and it will be removed in the future.
}

// SetCookie sets a cookie with the given cookie data; may overwrite
// equivalent cookies if they exist.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setCookie
//
// parameters:
//   name - Cookie name.
//   value - Cookie value.
func SetCookie(name string, value string) *SetCookieParams {
	return &SetCookieParams{
		Name:  name,
		Value: value,
	}
}

// WithURL the request-URI to associate with the setting of the cookie. This
// value can affect the default domain, path, source port, and source scheme
// values of the created cookie.
func (p SetCookieParams) WithURL(url string) *SetCookieParams {
	p.URL = url
	return &p
}

// WithDomain cookie domain.
func (p SetCookieParams) WithDomain(domain string) *SetCookieParams {
	p.Domain = domain
	return &p
}

// WithPath cookie path.
func (p SetCookieParams) WithPath(path string) *SetCookieParams {
	p.Path = path
	return &p
}

// WithSecure true if cookie is secure.
func (p SetCookieParams) WithSecure(secure bool) *SetCookieParams {
	p.Secure = secure
	return &p
}

// WithHTTPOnly true if cookie is http-only.
func (p SetCookieParams) WithHTTPOnly(httpOnly bool) *SetCookieParams {
	p.HTTPOnly = httpOnly
	return &p
}

// WithSameSite cookie SameSite type.
func (p SetCookieParams) WithSameSite(sameSite CookieSameSite) *SetCookieParams {
	p.SameSite = sameSite
	return &p
}

// WithExpires cookie expiration date, session cookie if not set.
func (p SetCookieParams) WithExpires(expires *cdp.TimeSinceEpoch) *SetCookieParams {
	p.Expires = expires
	return &p
}

// WithPriority cookie Priority type.
func (p SetCookieParams) WithPriority(priority CookiePriority) *SetCookieParams {
	p.Priority = priority
	return &p
}

// WithSameParty true if cookie is SameParty.
func (p SetCookieParams) WithSameParty(sameParty bool) *SetCookieParams {
	p.SameParty = sameParty
	return &p
}

// WithSourceScheme cookie source scheme type.
func (p SetCookieParams) WithSourceScheme(sourceScheme CookieSourceScheme) *SetCookieParams {
	p.SourceScheme = sourceScheme
	return &p
}

// WithSourcePort cookie source port. Valid values are {-1, [1, 65535]}, -1
// indicates an unspecified port. An unspecified port value allows protocol
// clients to emulate legacy cookie scope for the port. This is a temporary
// ability and it will be removed in the future.
func (p SetCookieParams) WithSourcePort(sourcePort int64) *SetCookieParams {
	p.SourcePort = sourcePort
	return &p
}

// Do executes Network.setCookie against the provided context.
func (p *SetCookieParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetCookie, p, nil)
}

// SetCookiesParams sets given cookies.
type SetCookiesParams struct {
	Cookies []*CookieParam `json:"cookies"` // Cookies to be set.
}

// SetCookies sets given cookies.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setCookies
//
// parameters:
//   cookies - Cookies to be set.
func SetCookies(cookies []*CookieParam) *SetCookiesParams {
	return &SetCookiesParams{
		Cookies: cookies,
	}
}

// Do executes Network.setCookies against the provided context.
func (p *SetCookiesParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetCookies, p, nil)
}

// SetExtraHTTPHeadersParams specifies whether to always send extra HTTP
// headers with the requests from this page.
type SetExtraHTTPHeadersParams struct {
	Headers Headers `json:"headers"` // Map with extra HTTP headers.
}

// SetExtraHTTPHeaders specifies whether to always send extra HTTP headers
// with the requests from this page.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setExtraHTTPHeaders
//
// parameters:
//   headers - Map with extra HTTP headers.
func SetExtraHTTPHeaders(headers Headers) *SetExtraHTTPHeadersParams {
	return &SetExtraHTTPHeadersParams{
		Headers: headers,
	}
}

// Do executes Network.setExtraHTTPHeaders against the provided context.
func (p *SetExtraHTTPHeadersParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetExtraHTTPHeaders, p, nil)
}

// SetAttachDebugStackParams specifies whether to attach a page script stack
// id in requests.
type SetAttachDebugStackParams struct {
	Enabled bool `json:"enabled"` // Whether to attach a page script stack for debugging purpose.
}

// SetAttachDebugStack specifies whether to attach a page script stack id in
// requests.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setAttachDebugStack
//
// parameters:
//   enabled - Whether to attach a page script stack for debugging purpose.
func SetAttachDebugStack(enabled bool) *SetAttachDebugStackParams {
	return &SetAttachDebugStackParams{
		Enabled: enabled,
	}
}

// Do executes Network.setAttachDebugStack against the provided context.
func (p *SetAttachDebugStackParams) Do(ctx context.Context) (err error) {
	return cdp.Execute(ctx, CommandSetAttachDebugStack, p, nil)
}

// GetSecurityIsolationStatusParams returns information about the COEP/COOP
// isolation status.
type GetSecurityIsolationStatusParams struct {
	FrameID cdp.FrameID `json:"frameId,omitempty"` // If no frameId is provided, the status of the target is provided.
}

// GetSecurityIsolationStatus returns information about the COEP/COOP
// isolation status.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-getSecurityIsolationStatus
//
// parameters:
func GetSecurityIsolationStatus() *GetSecurityIsolationStatusParams {
	return &GetSecurityIsolationStatusParams{}
}

// WithFrameID if no frameId is provided, the status of the target is
// provided.
func (p GetSecurityIsolationStatusParams) WithFrameID(frameID cdp.FrameID) *GetSecurityIsolationStatusParams {
	p.FrameID = frameID
	return &p
}

// GetSecurityIsolationStatusReturns return values.
type GetSecurityIsolationStatusReturns struct {
	Status *SecurityIsolationStatus `json:"status,omitempty"`
}

// Do executes Network.getSecurityIsolationStatus against the provided context.
//
// returns:
//   status
func (p *GetSecurityIsolationStatusParams) Do(ctx context.Context) (status *SecurityIsolationStatus, err error) {
	// execute
	var res GetSecurityIsolationStatusReturns
	err = cdp.Execute(ctx, CommandGetSecurityIsolationStatus, p, &res)
	if err != nil {
		return nil, err
	}

	return res.Status, nil
}

// LoadNetworkResourceParams fetches the resource and returns the content.
type LoadNetworkResourceParams struct {
	FrameID cdp.FrameID                 `json:"frameId"` // Frame id to get the resource for.
	URL     string                      `json:"url"`     // URL of the resource to get content for.
	Options *LoadNetworkResourceOptions `json:"options"` // Options for the request.
}

// LoadNetworkResource fetches the resource and returns the content.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-loadNetworkResource
//
// parameters:
//   frameID - Frame id to get the resource for.
//   url - URL of the resource to get content for.
//   options - Options for the request.
func LoadNetworkResource(frameID cdp.FrameID, url string, options *LoadNetworkResourceOptions) *LoadNetworkResourceParams {
	return &LoadNetworkResourceParams{
		FrameID: frameID,
		URL:     url,
		Options: options,
	}
}

// LoadNetworkResourceReturns return values.
type LoadNetworkResourceReturns struct {
	Resource *LoadNetworkResourcePageResult `json:"resource,omitempty"`
}

// Do executes Network.loadNetworkResource against the provided context.
//
// returns:
//   resource
func (p *LoadNetworkResourceParams) Do(ctx context.Context) (resource *LoadNetworkResourcePageResult, err error) {
	// execute
	var res LoadNetworkResourceReturns
	err = cdp.Execute(ctx, CommandLoadNetworkResource, p, &res)
	if err != nil {
		return nil, err
	}

	return res.Resource, nil
}

// Command names.
const (
	CommandSetAcceptedEncodings                    = "Network.setAcceptedEncodings"
	CommandClearAcceptedEncodingsOverride          = "Network.clearAcceptedEncodingsOverride"
	CommandClearBrowserCache                       = "Network.clearBrowserCache"
	CommandClearBrowserCookies                     = "Network.clearBrowserCookies"
	CommandDeleteCookies                           = "Network.deleteCookies"
	CommandDisable                                 = "Network.disable"
	CommandEmulateNetworkConditions                = "Network.emulateNetworkConditions"
	CommandEnable                                  = "Network.enable"
	CommandGetAllCookies                           = "Network.getAllCookies"
	CommandGetCertificate                          = "Network.getCertificate"
	CommandGetCookies                              = "Network.getCookies"
	CommandGetResponseBody                         = "Network.getResponseBody"
	CommandGetRequestPostData                      = "Network.getRequestPostData"
	CommandGetResponseBodyForInterception          = "Network.getResponseBodyForInterception"
	CommandTakeResponseBodyForInterceptionAsStream = "Network.takeResponseBodyForInterceptionAsStream"
	CommandReplayXHR                               = "Network.replayXHR"
	CommandSearchInResponseBody                    = "Network.searchInResponseBody"
	CommandSetBlockedURLS                          = "Network.setBlockedURLs"
	CommandSetBypassServiceWorker                  = "Network.setBypassServiceWorker"
	CommandSetCacheDisabled                        = "Network.setCacheDisabled"
	CommandSetCookie                               = "Network.setCookie"
	CommandSetCookies                              = "Network.setCookies"
	CommandSetExtraHTTPHeaders                     = "Network.setExtraHTTPHeaders"
	CommandSetAttachDebugStack                     = "Network.setAttachDebugStack"
	CommandGetSecurityIsolationStatus              = "Network.getSecurityIsolationStatus"
	CommandLoadNetworkResource                     = "Network.loadNetworkResource"
)
