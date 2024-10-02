# Espia-Go 🕵️
`espia-go` is a Go package that provides a convenient wrapper around the
`analytics.proper.ai` API for tracking user interaction events. It ensures
that all requests are well-formed and allows for sensible defaults like
automatic session generation.

## Installation
To install espia-go, use the following command:

```bash
go get github.com/propertechnologies/espia-go
```
Setting up the Espia instance:

```go
import "github.com/propertechnologies/espia-go/espia"

func main() {
    // Initialize Espia
    espia.Espia(espia.EspiaSetup{
        Source:      "your-project-name",
        AutoSession: true,
        Enabled:     true,
    })
}
```

Tracking events anywhere in your project:

```go
import "github.com/propertechnologies/espia-go/espia"

func main() {
    // Track an event with optional metadata
    err := espia.Track("category_label", espia.Metadata{
        "foo": "optional object",
    }, nil)

    if err != nil {
        fmt.Println("Error tracking event:", err)
    }
}
```

## Features

- Automatic Session Generation: Automatically creates a session ID for each user, unless provided.
- Metadata Handling: Supports permanent metadata that is included with every event, along with event-specific metadata.
- Flexible Configuration: Allows enabling or disabling tracking, custom session management, and more.

