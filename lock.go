package deimosclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrLockNotAcquired = errors.New("failed to acquire lock")
	ErrLockNotHeld     = errors.New("lock is not held by this client")
	ErrLockExpired     = errors.New("lock has expired")
)

// DistributedLock represents a distributed lock
type DistributedLock struct {
	client    *Client
	key       string
	value     string
	ttl       time.Duration
	renewalCh chan struct{}
	stopCh    chan struct{}
	mu        sync.RWMutex
	held      bool
	lastIndex uint64
}

// LockOptions contains options for creating a distributed lock
type LockOptions struct {
	ttl           time.Duration
	RenewalPeriod time.Duration
	AutoRenewal   bool
}

// DefaultLockOptions returns default lock options
func DefaultLockOptions() *LockOptions {
	return &LockOptions{
		ttl:           30 * time.Second,
		RenewalPeriod: 10 * time.Second,
		AutoRenewal:   true,
	}
}

func newLockOptions(options []LockOption) *LockOptions {
	opts := DefaultLockOptions()
	for _, opt := range options {
		opt.applyToLock(opts)
	}
	return opts
}

// NewDistributedLock creates a new distributed lock
func (c *Client) NewDistributedLock(key, value string, opts ...LockOption) *DistributedLock {
	lockOpts := newLockOptions(opts)

	return &DistributedLock{
		client:    c,
		key:       key,
		value:     value,
		ttl:       lockOpts.ttl,
		renewalCh: make(chan struct{}),
		stopCh:    make(chan struct{}),
	}
}

// TryLock attempts to acquire the lock without blocking
func (l *DistributedLock) TryLock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.held {
		return nil // Already holding the lock
	}

	// Try to create the lock key only if it doesn't exist (atomic create)
	resp, err := l.client.Set(ctx, l.key, l.value, WithTTL(l.ttl), WithPrevExist(false))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLockNotAcquired, err)
	}

	l.held = true
	l.lastIndex = resp.Node.ModifiedIndex
	return nil
}

// Lock acquires the lock, blocking until successful or context is cancelled
func (l *DistributedLock) Lock(ctx context.Context) error {
	// First, try to acquire the lock directly.
	err := l.TryLock(ctx)
	if err == nil {
		// Lock acquired successfully.
		return nil
	}
	if !errors.Is(err, ErrLockNotAcquired) {
		// An unexpected error occurred.
		return err
	}

	// If the lock is held by someone else, watch for it to be released.
	return l.watchAndLock(ctx)
}

func (l *DistributedLock) watchAndLock(ctx context.Context) error {
	for {
		// Start watching the lock key for changes.
		watchChan := l.client.Watch(ctx, l.key)

		// Wait for a DELETE event on the lock key.
		for resp := range watchChan {
			if resp.Action == "delete" || resp.Action == "compareAndDelete" || resp.Action == "expire" {
				// The lock seems to have been released. Try to acquire it.
				err := l.TryLock(ctx)
				if err == nil {
					// Lock acquired successfully.
					return nil
				}
				if !errors.Is(err, ErrLockNotAcquired) {
					// An unexpected error occurred.
					return err
				}
				// If we fail to acquire the lock, it means someone else got it.
				// We will restart the watch loop.
				break
			}
		}

		// If the watch channel closes or we break the inner loop, check the context.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue to the next iteration of the outer loop to re-establish the watch.
		}
	}
}

// Unlock releases the lock
func (l *DistributedLock) Unlock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.held {
		return ErrLockNotHeld
	}

	// Stop auto-renewal if running
	select {
	case l.stopCh <- struct{}{}:
	default:
	}

	// Use compare-and-delete to ensure we only delete our own lock
	_, err := l.client.CompareAndDelete(ctx, l.key, WithPrevValue(l.value))
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	l.held = false
	l.lastIndex = 0
	return nil
}

// IsHeld returns true if the lock is currently held
func (l *DistributedLock) IsHeld() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.held
}

// Renew extends the lock's TTL
func (l *DistributedLock) Renew(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.held {
		return ErrLockNotHeld
	}

	// Use compare-and-swap to renew the lock
	resp, err := l.client.CompareAndSwap(ctx, l.key, l.value,
		WithPrevValue(l.value),
		WithCasTTL(l.ttl))
	if err != nil {
		l.held = false
		return fmt.Errorf("%w: %v", ErrLockExpired, err)
	}

	l.lastIndex = resp.Node.ModifiedIndex
	return nil
}

// StartAutoRenewal starts automatic lock renewal in the background
func (l *DistributedLock) StartAutoRenewal(ctx context.Context, renewalPeriod time.Duration) {
	go l.autoRenewalLoop(ctx, renewalPeriod)
}

// autoRenewalLoop runs the automatic renewal process
func (l *DistributedLock) autoRenewalLoop(ctx context.Context, renewalPeriod time.Duration) {
	ticker := time.NewTicker(renewalPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-l.stopCh:
			return
		case <-ticker.C:
			if l.IsHeld() {
				if err := l.Renew(ctx); err != nil {
					// Log error but continue trying
					// In a real implementation, you might want to use a proper logger
					fmt.Printf("Failed to renew lock %s: %v\n", l.key, err)
				}
			}
		}
	}
}

// WithLock executes a function while holding the lock
func (l *DistributedLock) WithLock(ctx context.Context, fn func() error) error {
	if err := l.Lock(ctx); err != nil {
		return err
	}
	defer func() {
		if unlockErr := l.Unlock(ctx); unlockErr != nil {
			fmt.Printf("Failed to unlock: %v\n", unlockErr)
		}
	}()

	return fn()
}

// LockInfo returns information about the current lock state
type LockInfo struct {
	Key       string
	Value     string
	Held      bool
	LastIndex uint64
	TTL       time.Duration
}

// Info returns information about the lock
func (l *DistributedLock) Info() LockInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return LockInfo{
		Key:       l.key,
		Value:     l.value,
		Held:      l.held,
		LastIndex: l.lastIndex,
		TTL:       l.ttl,
	}
}
