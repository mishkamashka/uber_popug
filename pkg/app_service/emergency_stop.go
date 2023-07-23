package app_service

// Returns `true` or not nil `err` to stop `launcher`
type Command func(l *launcher) (bool, error)

// Stops the launcher
func EmergencyStop(err error) Command {
	return func(_ *launcher) (bool, error) {
		return true, err
	}
}
