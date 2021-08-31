package appnotice

type nopNotifier struct{}

// NewNopNotifier returns a notifier that does not notify anything.
func NewNopNotifier() Notifier {
	return &nopNotifier{}
}

// Error does nothing.
func (*nopNotifier) Error(err error) error {
	return nil
}

// Critical does nothing.
func (*nopNotifier) Critical(err error) error {
	return nil
}
