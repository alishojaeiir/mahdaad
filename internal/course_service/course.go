package course_service

import (
	"context"
	"fmt"
	"github.com/alishojaeiir/mahdaad/internal/events"
	"github.com/alishojaeiir/mahdaad/pkg/eventbus"
)

// CourseService handles course_service-related business logic in the application layer.
type CourseService struct {
	eventBus *eventbus.EventBus
}

// NewCourseService creates a new CourseService with an EventBus dependency.
func NewCourseService(eventBus *eventbus.EventBus) *CourseService {
	return &CourseService{eventBus: eventBus}
}

// CreateCourse handles the creation of a course_service and publishes a domain event.
func (s *CourseService) CreateCourse(ctx context.Context, id, name string) error {
	// Simulate domain logic (e.g., validation, persistence)
	fmt.Printf("Course created: ID=%s, Name=%s\n", id, name)

	// Publish domain event asynchronously
	event := events.CourseCreated{CourseID: id, CourseName: name}
	s.eventBus.Publish(ctx, event)

	return nil
}
