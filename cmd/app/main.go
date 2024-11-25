package main

import (
	"fmt"

	"github.com/azimari-toure-ikbal/translate-core/internal"
)

// func gracefulShutdown() {
// 	// Create context that listens for the interrupt signal from the OS.
// 	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	// Listen for the interrupt signal.
// 	<-ctx.Done()

// 	log.Println("shutting down gracefully, press Ctrl+C again to force")

// 	// The context is used to inform the server it has 5 seconds to finish
// 	// the request it is currently handling
// 	_, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	log.Println("Server exiting")
// }

func main() {
	err := internal.Run()

	// go gracefulShutdown()

	if err != nil {
		fmt.Printf("Something went wrong: %s\n", err)
	}
}
