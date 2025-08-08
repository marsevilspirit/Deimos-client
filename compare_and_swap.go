package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CompareAndSwapOptions struct {
	ttl       time.Duration
	prevValue string
	prevIndex uint64
}

type CompareAndSwapOption interface {
	applyToCompareAndSwap(*CompareAndSwapOptions)
}

func newCompareAndSwapOptions(options []CompareAndSwapOption) *CompareAndSwapOptions {
	opts := CompareAndSwapOptions{}
	for _, opt := range options {
		opt.applyToCompareAndSwap(&opts)
	}
	return &opts
}

// CompareAndSwap performs an atomic compare-and-swap operation
func (c *Client) CompareAndSwap(ctx context.Context, key, value string, opts ...CompareAndSwapOption) (*Response, error) {
	casOpts := newCompareAndSwapOptions(opts)

	URL := c.buildURL(key)
	query := url.Values{}
	query.Set("value", value)

	if casOpts.ttl > 0 {
		query.Set("ttl", fmt.Sprintf("%d", int64(casOpts.ttl.Seconds())))
	}

	if casOpts.prevValue != "" {
		query.Set("prevValue", casOpts.prevValue)
	}

	if casOpts.prevIndex > 0 {
		query.Set("prevIndex", strconv.FormatUint(casOpts.prevIndex, 10))
	}

	body := strings.NewReader(query.Encode())
	req, err := http.NewRequestWithContext(ctx, "PUT", URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

// CompareAndDelete performs an atomic compare-and-delete operation
func (c *Client) CompareAndDelete(ctx context.Context, key string, opts ...CompareAndDeleteOption) (*Response, error) {
	cadOpts := newCompareAndDeleteOptions(opts)

	URL := c.buildURL(key)
	query := url.Values{}

	if cadOpts.prevValue != "" {
		query.Set("prevValue", cadOpts.prevValue)
	}

	if cadOpts.prevIndex > 0 {
		query.Set("prevIndex", strconv.FormatUint(cadOpts.prevIndex, 10))
	}

	body := strings.NewReader(query.Encode())
	req, err := http.NewRequestWithContext(ctx, "DELETE", URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

type CompareAndDeleteOptions struct {
	prevValue string
	prevIndex uint64
}

type CompareAndDeleteOption interface {
	applyToCompareAndDelete(*CompareAndDeleteOptions)
}

func newCompareAndDeleteOptions(options []CompareAndDeleteOption) *CompareAndDeleteOptions {
	opts := CompareAndDeleteOptions{}
	for _, opt := range options {
		opt.applyToCompareAndDelete(&opts)
	}
	return &opts
}

// Options for CompareAndSwap
func WithPrevValue(prevValue string) interface {
	CompareAndSwapOption
	CompareAndDeleteOption
} {
	return &prevValueOption{prevValue: prevValue}
}

func WithPrevIndex(prevIndex uint64) interface {
	CompareAndSwapOption
	CompareAndDeleteOption
} {
	return &prevIndexOption{prevIndex: prevIndex}
}

func WithCasTTL(ttl time.Duration) CompareAndSwapOption {
	return &casTTLOption{ttl: ttl}
}

type prevValueOption struct {
	prevValue string
}

func (o *prevValueOption) applyToCompareAndSwap(opts *CompareAndSwapOptions) {
	opts.prevValue = o.prevValue
}

func (o *prevValueOption) applyToCompareAndDelete(opts *CompareAndDeleteOptions) {
	opts.prevValue = o.prevValue
}

type prevIndexOption struct {
	prevIndex uint64
}

func (o *prevIndexOption) applyToCompareAndSwap(opts *CompareAndSwapOptions) {
	opts.prevIndex = o.prevIndex
}

func (o *prevIndexOption) applyToCompareAndDelete(opts *CompareAndDeleteOptions) {
	opts.prevIndex = o.prevIndex
}

type casTTLOption struct {
	ttl time.Duration
}

func (o *casTTLOption) applyToCompareAndSwap(opts *CompareAndSwapOptions) {
	opts.ttl = o.ttl
}
