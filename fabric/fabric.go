package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"regexp"
	"strings"
)

// SubscriptionType specifies whether a route subscribes as a "queue."
// Queue subscriptions allow multiple subscriptions to the same endpoint.
// Non-queue subscriptions do not; only a single non-queue subscription
// to a given endpoint will receive any messages.
// Queue is the default.
type SubscriptionType int

const (
	Queue SubscriptionType = iota
	NotQueue
)

// Request is passed over to clients of this library. It's possible for clients
// to modify it and then pass it along a middleware chain as is done with Go's http
// routing, although we don't have any support for this at the moment.
// VERY IMPORTANT: The Body MUST be a json object. We don't support any other kind
// of body, and other types will throw an error when we marshal the message from
// the caller.
type Request struct {
	RawHeaders []string       `json:"headers"`
	Body       map[string]any `json:"body"`

	Method   string
	Endpoint string
	Headers  map[string][]string
}

// Reply is used by clients of this library.
type Reply struct {
	Status  int
	Headers []string
	Error   error
	Body    any
}

// Response is used internally in this library.
type Response struct {
	Status  int      `json:"status"`
	Headers []string `json:"headers"`
	Errors  []string `json:"errors"`
	Body    any      `json:"body"`
}

// respond
func (r *Response) respond(msg *nats.Msg) {
	j, _ := json.Marshal(r)
	if msg.Reply != "" {
		msg.Respond(j)
	}
}

// Route specifies handlers for method/endpoint pairs. Someday we will need
// to support wildcards and patterns, as with URL routing.
// The Handler takes pointers to request and reply objects, unlike HTTP routing,
// because replies are discrete messages rather than streams, as in HTTP.
type Route struct {
	Method   string
	Endpoint string
	Handler  func(*Reply, *Request)
	Type     SubscriptionType

	subject_prefix string
	subject_suffix string
}

// ServeCfg provides parameters that are optional (Verbose) or redundant with
// identify parameters. This is useful when the caller is running in a docker
// container rather than as a standard agent, and the identity parameters aren't
// obtainable in the usual way.
type ServeCfg struct {
	Tenant  string
	AgentId string
	NatsUrl string
	Verbose bool
}

// Serve
func Serve(ctx context.Context, routes []Route, cfg *ServeCfg) (err error) {
	if len(routes) == 0 {
		err = fmt.Errorf("no routes specified")
		return
	}

	// the config struct is optional.
	if cfg == nil {
		cfg = &ServeCfg{}
	}

	// look up identifiers. If present in the config, they take precedence.
	tenant := get_tenant(cfg)
	agentid := get_agentid(cfg)
	natsurl := get_natsurl(cfg)

	if tenant == "" || agentid == "" || natsurl == "" {
		err = fmt.Errorf("missing identifiers")
		return
	}

	log.Printf("%v,%v,%v", tenant, agentid, natsurl)

	// connect to NATS
	nc, err := nats.Connect(natsurl)
	if err != nil {
		return
	}
	defer nc.Drain()

	// subscribe to routes
	for _, r := range routes {
		if e := r.subscribe(nc, tenant, agentid); e == nil {
		} else {
			log.Printf("subscription error %v", e)
		}
	}

	// TODO, schedule nats reconnects. The connection doesn't self-heal.
	select {
	case <-ctx.Done():
	}

	return
}

// subscribe
func (route *Route) subscribe(nc *nats.Conn, tenant string, agent string) error {
	// validate the route verb, which becomes part of the subscription subject.
	verb := strings.ToLower(strings.TrimSpace(route.Method))
	verbs := []string{"get", "put", "post", "delete", "patch"}
	{
		ok := false
		for _, v := range verbs {
			if v == verb {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("invalid verb %s", verb)
		}
	}

	// validate the endpoint, which may be hierarchical and contain NATS wildcards (* and terminal >)
	ep := strings.ToLower(strings.TrimSpace(route.Endpoint))
	{
		re := regexp.MustCompile(`^[A-Za-z0-9\.\*_/>-]+$`)
		if !re.MatchString(ep) {
			return fmt.Errorf("invalid endpoint %s", ep)
		}
	}

	// create the subscription subject, which may contain wildcards.
	// the point of this is so we can easily strip out the actual endpoint when handling a msg.
	route.subject_prefix = fmt.Sprintf("agent.rest.%s.%s.", tenant, agent)
	route.subject_suffix = fmt.Sprintf("%s.%s", verb, ep)

	subject := route.subject_prefix + route.subject_suffix

	nats_handler := func(m *nats.Msg) {
		route.handle_msg(m)
	}

	if route.Type == Queue {
		if _, e := nc.QueueSubscribe(subject, "workers", nats_handler); e != nil {
			return e
		}
	} else {
		if _, e := nc.Subscribe(subject, nats_handler); e != nil {
			return e
		}
	}

	return nil
}

// handle_msg
func (route *Route) handle_msg(msg *nats.Msg) {
	send_error := func(e error, status int) {
		re := Response{
			Status: status,
			Errors: []string{fmt.Sprintf("%v", e)},
		}
		re.respond(msg)
	}

	var req Request
	e := json.Unmarshal(msg.Data, &req)
	if e != nil {
		send_error(e, 500)
		return
	}
	req.parseHeaders()

	subj := strings.TrimPrefix(msg.Subject, route.subject_prefix)
	parts := strings.SplitN(subj, ".", 2)
	if len(parts) == 2 {
		req.Method = parts[0]
		req.Endpoint = parts[1]
	} else {
		req.Endpoint = subj
	}

	// Call the client's handler on a goroutine because it could be a lengthy
	// operation like an inference. Capture panics in case their code isn't
	// cleanly written.
	// Remember that we convert a Reply from the client into a Response here.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				send_error(fmt.Errorf("%v", r), 500)
			}
		}()

		reply := NewReply()
		route.Handler(reply, &req)

		// convert the reply from the user code into a Response.
		resp := reply.to_response()
		// ALWAYS send a response, even if the Body is nil
		resp.respond(msg)

	}()

}

// to_response converts a Reply object filled in by client code to a Response
// object that we will use to reply to a NATS message.
// VERY IMPORTANT: if the user's Reply object has a nil Body with a 200 status,
// NO RESPONSE will be sent!
func (r *Reply) to_response() *Response {

	re := Response{
		Status: r.Status,
	}

	if r.Status == 200 && r.Error == nil {
		// should set a default header? would need a type-assertion switch based on the body type
		re.Headers = r.Headers
		re.Body = r.Body
	} else {
		re.Errors = []string{fmt.Sprintf("%v", r.Error)}
		re.Status = 500
	}

	return &re
}

// NewReply
func NewReply() *Reply {
	return &Reply{
		Status:  200,
		Headers: []string{},
	}
}

// parseHeaders reads the headers in a Request object and saves them as a map[string][]string.
// This is a convenience function.
func (r *Request) parseHeaders() {
	r.Headers = map[string][]string{}

	if r.RawHeaders == nil || len(r.RawHeaders) == 0 {
		return
	}

	for _, h := range r.RawHeaders {
		y := strings.Split(h, ":")
		if len(y) == 2 {
			y1 := strings.ToLower(strings.TrimSpace(y[0]))
			y2 := strings.TrimSpace(y[1])

			if _, ok := r.Headers[y1]; !ok {
				r.Headers[y1] = []string{}
			}

			r.Headers[y1] = append(r.Headers[y1], y2)
		}
	}
}

// AddHeader ASSUMES that the Headers field has been initialized.
func (r *Reply) AddHeader(h string) {
	r.Headers = append(r.Headers, h)
}
