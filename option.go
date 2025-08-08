package deimosclient

import "time"

type SetOption interface {
	applyToSet(*SetOptions)
}

type GetOption interface {
	applyToGet(*GetOptions)
}

type DeleteOption interface {
	applyToDelete(*DeleteOptions)
}

type WatchOption interface {
	applyToWatch(*WatchOptions)
}

type LockOption interface {
	applyToLock(*LockOptions)
}

type SetDeleteOption interface {
	SetOption
	DeleteOption
}

type SetLockOption interface {
	SetOption
	LockOption
}

type GetDeleteWatchOption interface {
	GetOption
	DeleteOption
	WatchOption
}

// TTL option
func WithTTL(ttl time.Duration) SetLockOption {
	return &ttlOption{ttl: ttl}
}

type ttlOption struct {
	ttl time.Duration
}

func (o *ttlOption) applyToSet(opts *SetOptions) {
	opts.ttl = o.ttl
}

func (o *ttlOption) applyToLock(opts *LockOptions) {
	opts.ttl = o.ttl
}

// Dir option
func WithDir() SetDeleteOption {
	return &dirOption{dir: true}
}

type dirOption struct {
	dir bool
}

func (o *dirOption) applyToSet(opts *SetOptions) {
	opts.dir = o.dir
}

func (o *dirOption) applyToDelete(opts *DeleteOptions) {
	opts.dir = o.dir
}

// prevExist option
func WithPrevExist(exists bool) SetOption {
	return &prevExistOption{prevExist: exists}
}

type prevExistOption struct {
	prevExist bool
}

func (o *prevExistOption) applyToSet(opts *SetOptions) {
	opts.prevExist = &o.prevExist
}

// recursive option
func WithRecursive() GetDeleteWatchOption {
	return &recursiveOption{recursive: true}
}

type recursiveOption struct {
	recursive bool
}

func (o *recursiveOption) applyToGet(opts *GetOptions) {
	opts.recursive = o.recursive
}

func (o *recursiveOption) applyToDelete(opts *DeleteOptions) {
	opts.recursive = o.recursive
}

func (o *recursiveOption) applyToWatch(opts *WatchOptions) {
	opts.recursive = o.recursive
}

// waitIndex option
func WithWaitIndex(waitIndex uint64) WatchOption {
	return &waitIndexOption{waitIndex: waitIndex}
}

type waitIndexOption struct {
	waitIndex uint64
}

func (o *waitIndexOption) applyToWatch(opts *WatchOptions) {
	opts.waitIndex = o.waitIndex
}

// WithRenewalPeriod sets the renewal period for auto-renewal
func WithRenewalPeriod(period time.Duration) LockOption {
	return &renewalPeriodOption{period: period}
}

type renewalPeriodOption struct {
	period time.Duration
}

func (o *renewalPeriodOption) applyToLock(opts *LockOptions) {
	opts.RenewalPeriod = o.period
}

// WithAutoRenewal enables or disables auto-renewal
func WithAutoRenewal(enabled bool) LockOption {
	return &autoRenewalOption{enabled: enabled}
}

type autoRenewalOption struct {
	enabled bool
}

func (o *autoRenewalOption) applyToLock(opts *LockOptions) {
	opts.AutoRenewal = o.enabled
}
