# router.go Explanation

The `router.go` file defines the HTTP server and routing logic for the application. It sets up the routes for handling API requests, initializes the required dependencies, and provides middleware for handling errors and cross-origin requests. Below is a detailed explanation of its components:

---

## Imports
The file imports the following packages:
- **`fmt`**: A standard library package for formatted I/O, used for error messages.
- **`joi-energy-golang/api`**: A custom package for handling API responses.
- **`joi-energy-golang/endpoints/priceplans`**: A package for handling price plan-related endpoints.
- **`joi-energy-golang/endpoints/readings`**: A package for handling electricity readings-related endpoints.
- **`joi-energy-golang/repository`**: A package for accessing data repositories.
- **`log`**: A standard library package for logging messages.
- **`net/http`**: A standard library package for HTTP server and client implementations.
- **`os`**: A standard library package for interacting with the operating system.
- **`strings`**: A standard library package for string manipulation.
- **`github.com/julienschmidt/httprouter`**: A lightweight HTTP router for handling routes and URL parameters.

---

## Function: `NewServer`
```go
func NewServer() *http.Server {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("defaulting to port %s", port)
    }
    return &http.Server{Addr: "localhost:" + port, Handler: newHandler()}
}
```
This function creates and configures a new HTTP server:
1. **Port Configuration**:
   - Reads the `PORT` environment variable to determine the port number.
   - Defaults to port `8080` if the environment variable is not set.
2. **Output**: Returns an `http.Server` instance with the configured port and a handler created by `newHandler`.

---

## Function: `addRoutes`
```go
func addRoutes(r *httprouter.Router) {
    accounts := repository.NewAccounts(defaultSmartMeterToPricePlanAccounts())
    meterReadings := repository.NewMeterReadings(defaultMeterElectricityReadings())
    pricePlans := repository.NewPricePlans(defaultPricePlans(), &meterReadings)

    readingsHandler := readings.NewHandler(&meterReadings)
    pricePlanHandler := priceplans.NewHandler(priceplans.NewService(&pricePlans, &accounts))

    r.POST("/readings/store", readingsHandler.StoreReadings)
    r.GET("/readings/read/:smartMeterId", readingsHandler.GetReadings)

    r.GET("/price-plans/compare-all/:smartMeterId", pricePlanHandler.CompareAll)
    r.GET("/price-plans/recommend/:smartMeterId", pricePlanHandler.Recommend)
}
```
This function sets up the routes for the application:
1. **Dependencies**:
   - Initializes repositories (`Accounts`, `MeterReadings`, `PricePlans`) with default data.
   - Creates handlers for `readings` and `priceplans` endpoints.
2. **Routes**:
   - **Readings**:
     - `POST /readings/store`: Stores electricity readings.
     - `GET /readings/read/:smartMeterId`: Retrieves electricity readings for a specific smart meter.
   - **Price Plans**:
     - `GET /price-plans/compare-all/:smartMeterId`: Compares all price plans for a specific smart meter.
     - `GET /price-plans/recommend/:smartMeterId`: Recommends price plans for a specific smart meter.

---

## Function: `newHandler`
```go
func newHandler() http.Handler {
    r := httprouter.New()
    addRoutes(r)

    r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Access-Control-Request-Method") != "" {
            header := w.Header()
            header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
            header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
        }
        w.WriteHeader(http.StatusNoContent)
    })

    r.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
        log.Printf("panic: %+v", err)
        api.Error(w, r, fmt.Errorf("whoops! My handler has run into a panic"), http.StatusInternalServerError)
    }
    r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        api.Error(w, r, fmt.Errorf("we have OPTIONS for you but %v is not among them", r.Method), http.StatusMethodNotAllowed)
    })
    r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.Contains(r.Header.Get("Accept"), "text/html") {
            sendBrowserDoc(w, r)
            return
        }
        api.Error(w, r, fmt.Errorf("whatever route you've been looking for, it's not here"), http.StatusNotFound)
    })

    return r
}
```
This function creates a new HTTP handler with middleware:
1. **Routes**: Calls `addRoutes` to register API routes.
2. **Global OPTIONS**: Handles preflight requests for CORS (Cross-Origin Resource Sharing).
3. **Panic Handler**: Logs panics and responds with a `500 Internal Server Error`.
4. **Method Not Allowed**: Responds with a `405 Method Not Allowed` error for unsupported HTTP methods.
5. **Not Found**: Responds with a `404 Not Found` error for unknown routes. If the request accepts HTML, it serves a browser-specific document.

---

## Function: `sendBrowserDoc`
```go
func sendBrowserDoc(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusUnsupportedMediaType)
    b, err := os.ReadFile("browser.htm")
    if err != nil {
        api.Error(w, r, fmt.Errorf("read browser.htm failed: %w", err), http.StatusInternalServerError)
    }
    _, err = w.Write(b)
    if err != nil {
        api.Error(w, r, fmt.Errorf("send browser.htm failed: %w", err), http.StatusInternalServerError)
    }
}
```
This function serves an HTML document for browser-based requests:
1. **Operation**:
   - Reads the `browser.htm` file from the filesystem.
   - Writes the file's content to the response.
2. **Error Handling**: Responds with a `500 Internal Server Error` if the file cannot be read or written.

---

## Key Points
- **Routing**: The file uses `httprouter` to define and manage API routes.
- **Middleware**: Includes middleware for handling CORS, panics, unsupported methods, and unknown routes.
- **Dependency Injection**: Initializes repositories and handlers with default data, making the application modular and testable.
- **Error Handling**: Provides detailed error responses for various scenarios.

---

## Example
Here’s how the server can be started:
```go
func main() {
    server := router.NewServer()
    log.Fatal(server.ListenAndServe())
}
```

This file is the entry point for setting up the HTTP server and routing logic for the application.