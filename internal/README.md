# internal
--
    import "github.com/autom8ter/slashsub/internal"


## Usage

```go
var PROJECT_ID = os.Getenv("PROJECT_ID")
```

```go
var SLACK_SIGNING_SECRET = []byte(os.Getenv("SLACK_SIGNING_SECRET"))
```

```go
var SLASH_FUNCTION_URL = "https://us-central1-autom8ter-19.cloudfunctions.net/SlashFunction"
```

#### type SlashSub

```go
type SlashSub struct {
}
```


#### func  New

```go
func New(service string, middlewares ...driver.Middleware) (*SlashSub, error)
```

#### func (*SlashSub) Client

```go
func (s *SlashSub) Client() *driver.Client
```

#### func (*SlashSub) ListenAndServe

```go
func (s *SlashSub) ListenAndServe(addr string) error
```

#### func (*SlashSub) ServeHTTP

```go
func (s *SlashSub) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

#### func (*SlashSub) ValidateRequest

```go
func (s *SlashSub) ValidateRequest(r *http.Request) bool
```
