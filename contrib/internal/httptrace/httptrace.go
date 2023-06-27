
// Package httptrace provides functionalities to trace HTTP requests that are commonly required and used across
// contrib/** integrations.
package httptrace

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"inet.af/netaddr"

	"git.proto.group/protoobp/pobp-trace-go/pobptrace"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/ext"
	"git.proto.group/protoobp/pobp-trace-go/pobptrace/tracer"
)

var (
	ipv6SpecialNetworks = []*netaddr.IPPrefix{
		ippref("fec0::/10"), // site local
	}
	defaultIPHeaders = []string{
		"x-forwarded-for",
		"x-real-ip",
		"x-client-ip",
		"x-forwarded",
		"x-cluster-client-ip",
		"forwarded-for",
		"forwarded",
		"via",
		"true-client-ip",
	}
	cfg = newConfig()
)

// StartRequestSpan starts an HTTP request span with the standard list of HTTP request span tags (http.method, http.url,
// http.useragent). Any further span start option can be added with opts.
func StartRequestSpan(r *http.Request, opts ...pobptrace.StartSpanOption) (tracer.Span, context.Context) {
	// Append our span options before the given ones so that the caller can "overwrite" them.
	// TODO(): rework span start option handling
	opts = append([]pobptrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.Tag(ext.HTTPMethod, r.Method),
		tracer.Tag(ext.HTTPURL, urlFromRequest(r)),
		tracer.Tag(ext.HTTPUserAgent, r.UserAgent()),
		tracer.Measured(),
	}, opts...)
	if r.Host != "" {
		opts = append([]pobptrace.StartSpanOption{
			tracer.Tag("http.host", r.Host),
		}, opts...)
	}
	if cfg.clientIP {
		opts = append(genClientIPSpanTags(r), opts...)
	}
	if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header)); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}
	return tracer.StartSpanFromContext(r.Context(), "http.request", opts...)
}

// FinishRequestSpan finishes the given HTTP request span and sets the expected response-related tags such as the status
// code. Any further span finish option can be added with opts.
func FinishRequestSpan(s tracer.Span, status int, opts ...tracer.FinishOption) {
	var statusStr string
	if status == 0 {
		statusStr = "200"
	} else {
		statusStr = strconv.Itoa(status)
	}
	s.SetTag(ext.HTTPCode, statusStr)
	if status >= 500 && status < 600 {
		s.SetTag(ext.Error, fmt.Errorf("%s: %s", statusStr, http.StatusText(status)))
	}
	s.Finish(opts...)
}

// ippref returns the IP network from an IP address string s. If not possible, it returns nil.
func ippref(s string) *netaddr.IPPrefix {
	if prefix, err := netaddr.ParseIPPrefix(s); err == nil {
		return &prefix
	}
	return nil
}

// genClientIPSpanTags generates the client IP related tags that need to be added to the span.
func genClientIPSpanTags(r *http.Request) []pobptrace.StartSpanOption {
	ipHeaders := defaultIPHeaders
	if len(cfg.clientIPHeader) > 0 {
		ipHeaders = []string{cfg.clientIPHeader}
	}
	var headers []string
	var ips []string
	var opts []pobptrace.StartSpanOption
	for _, hdr := range ipHeaders {
		if v := r.Header.Get(hdr); v != "" {
			headers = append(headers, hdr)
			ips = append(ips, v)
		}
	}
	if len(ips) == 0 {
		if remoteIP := parseIP(r.RemoteAddr); remoteIP.IsValid() && isGlobal(remoteIP) {
			opts = append(opts, tracer.Tag(ext.HTTPClientIP, remoteIP.String()))
		}
	} else if len(ips) == 1 {
		for _, ipstr := range strings.Split(ips[0], ",") {
			ip := parseIP(strings.TrimSpace(ipstr))
			if ip.IsValid() && isGlobal(ip) {
				opts = append(opts, tracer.Tag(ext.HTTPClientIP, ip.String()))
				break
			}
		}
	} else {
		for i := range ips {
			opts = append(opts, tracer.Tag(ext.HTTPRequestHeaders+"."+headers[i], ips[i]))
		}
		opts = append(opts, tracer.Tag(ext.MultipleIPHeaders, strings.Join(headers, ",")))
	}
	return opts
}

func parseIP(s string) netaddr.IP {
	if ip, err := netaddr.ParseIP(s); err == nil {
		return ip
	}
	if h, _, err := net.SplitHostPort(s); err == nil {
		if ip, err := netaddr.ParseIP(h); err == nil {
			return ip
		}
	}
	return netaddr.IP{}
}

func isGlobal(ip netaddr.IP) bool {
	isGlobal := !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsLinkLocalUnicast()
	if !isGlobal || !ip.Is6() {
		return isGlobal
	}
	for _, n := range ipv6SpecialNetworks {
		if n.Contains(ip) {
			return false
		}
	}
	return isGlobal
}

// urlFromRequest returns the full URL from the HTTP request. If query params are collected, they are obfuscated granted
// obfuscation is not disabled by the user (through POBP_TRACE_OBFUSCATION_QUERY_STRING_REGEXP)
func urlFromRequest(r *http.Request) string {
	// Quoting net/http comments about net.Request.URL on server requests:
	// "For most requests, fields other than Path and RawQuery will be
	// empty. (See RFC 7230, Section 5.3)"
	// This is why we don't rely on url.URL.String(), url.URL.Host, url.URL.Scheme, etc...
	var url string
	path := r.URL.EscapedPath()
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if r.Host != "" {
		url = strings.Join([]string{scheme, "://", r.Host, path}, "")
	} else {
		url = path
	}
	// Collect the query string if we are allowed to report it and obfuscate it if possible/allowed
	if cfg.queryString && r.URL.RawQuery != "" {
		query := r.URL.RawQuery
		if cfg.queryStringRegexp != nil {
			query = cfg.queryStringRegexp.ReplaceAllLiteralString(query, "<redacted>")
		}
		url = strings.Join([]string{url, query}, "?")
	}
	if frag := r.URL.EscapedFragment(); frag != "" {
		url = strings.Join([]string{url, frag}, "#")
	}
	return url
}
