// Package appnotice defines an interface for notifications.
// It also defines an implementation for each notification destination.
package appnotice

// Notifier represents a notification interface
// that sends a message to each destination.
type Notifier interface {
	Error(err error) error
	Critical(err error) error
}
