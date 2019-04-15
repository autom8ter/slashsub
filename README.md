# slashsub
--
    import "github.com/autom8ter/slashsub"


## Usage

```go
var PROJECT_ID = os.Getenv("PROJECT_ID")
```

```go
var SERVICE = os.Getenv("SERVICE")
```

```go
var SLACK_SIGNING_SECRET = []byte(os.Getenv("SLACK_SIGNING_SECRET"))
```

```go
var SLASH_FUNCTION_URL = "https://us-central1-autom8ter-19.cloudfunctions.net/SlashFunction"
```

```go
var TOPIC = os.Getenv("TOPIC")
```

#### func  SlashFunction

```go
func SlashFunction(w http.ResponseWriter, r *http.Request)
```

#### type HandlerFunc

```go
type HandlerFunc func(ctx context.Context, msg *proto.Message, _ *api.Msg) error
```


#### type SlashSub

```go
type SlashSub struct {
}
```


#### func  New

```go
func New(middlewares ...driver.Middleware) (*SlashSub, error)
```

#### func (*SlashSub) ListenAndServe

```go
func (s *SlashSub) ListenAndServe(addr string) error
```

#### func (*SlashSub) ServeHTTP

```go
func (s *SlashSub) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

#### func (*SlashSub) Subscribe

```go
func (s *SlashSub) Subscribe(jSON bool, handler HandlerFunc)
```

#### func (*SlashSub) ValidateRequest

```go
func (s *SlashSub) ValidateRequest(r *http.Request) bool
```
