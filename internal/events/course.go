package events

// CourseCreated is a domain event raised when a services is created.
type CourseCreated struct {
	CourseID   string
	CourseName string
}

// Type returns the event type.
func (e CourseCreated) Type() string {
	return "CourseCreated"
}
