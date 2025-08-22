package main

import (
	"context"
	"fmt"
	"github.com/alishojaeiir/mahdaad/internal/eventbus"
	"github.com/alishojaeiir/mahdaad/internal/events"
	"github.com/alishojaeiir/mahdaad/internal/handlers"
	"github.com/alishojaeiir/mahdaad/internal/services"
	"time"
)

func main() {
	ctx := context.Background()

	// Initialize the event bus
	eventBus := eventbus.NewEventBus()

	// Set up handler channels
	emailCh := make(chan events.Event)
	dashboardCh := make(chan events.Event)
	searchCh := make(chan events.Event)

	// Subscribe handlers to the event bus
	eventBus.Subscribe("CourseCreated", emailCh)
	eventBus.Subscribe("CourseCreated", dashboardCh)
	eventBus.Subscribe("CourseCreated", searchCh)

	// Start handlers in goroutines
	go handlers.EmailHandler(emailCh)
	go handlers.DashboardHandler(dashboardCh)
	go handlers.SearchIndexerHandler(searchCh)

	// Initialize the course service
	courseService := services.NewCourseService(eventBus)

	// Simulate course creation
	err := courseService.CreateCourse(ctx, "123", "Advanced Go Programming")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Allow time for async handlers to process (in production, use proper synchronization)
	time.Sleep(1 * time.Second)
}
