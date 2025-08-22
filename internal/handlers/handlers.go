package handlers

import (
	"fmt"
	"github.com/alishojaeiir/mahdaad/internal/events"
)

// EmailHandler processes CourseCreated eventbus to send email notifications.
func EmailHandler(ch chan events.Event) {
	for event := range ch {
		if e, ok := event.(events.CourseCreated); ok {
			fmt.Printf("Sending email notification for course_service: %s\n", e.CourseName)
		}
	}
}

// DashboardHandler processes CourseCreated eventbus to update the admin dashboard.
func DashboardHandler(ch chan events.Event) {
	for event := range ch {
		if e, ok := event.(events.CourseCreated); ok {
			fmt.Printf("Updating admin dashboard for course_service: %s\n", e.CourseName)
		}
	}
}

// SearchIndexerHandler processes CourseCreated eventbus to index the course_service.
func SearchIndexerHandler(ch chan events.Event) {
	for event := range ch {
		if e, ok := event.(events.CourseCreated); ok {
			fmt.Printf("Indexing course_service in search system: %s\n", e.CourseName)
		}
	}
}
