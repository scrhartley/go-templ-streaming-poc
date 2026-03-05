# Go Templ Streaming Proof-of-Concept

Proof-of-concept of HTML streaming in Go implemented with [Templ](https://templ.guide/) templating. This implementation uses a function variant of the [Future](https://en.wikipedia.org/wiki/Futures_and_promises) pattern (e.g. `value()`) as opposed to the approach implied by Templ's documentation of using channels directly (e.g. `<-value`). This has the advantage of being able to read each value multiple times if desired and also improves the ergonomics of flushing when dealing with a group of computations.

This repo also includes a take on the [Error Boundary](https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary) error-handling mechanism. Traditional page rendering will try to collect all the data and if there's an error during this process the backend can choose to return or redirect to an error page instead. This isn't as straightforward with streaming and error boundaries offer an alternative strategy by instead rendering fallback content for a section of the page if it has failed.

**Endpoints:**  
`/async` (localhost:3000/async) - This is the primary page demonstrating the features in this repo. Due to concurrency, it may appear that most of the page completes simultaneously. The chance of this lessens when quicker computations are rendered before slower computations.  
`/sync` (localhost:3000/sync) - This is included for comparison purposes and shows how the same page behaves when concurrency is not used (in terms of ordering and speed), while still using streaming.  

**Useful commands:**  
`go tool templ generate` - Updates generated code when changes are made to the Templ template files.  
`go run main.go` - Will compile the code and run the server.  
