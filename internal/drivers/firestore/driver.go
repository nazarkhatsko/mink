package firestore

import (
	"context"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/nazarkhatsko/mink/pkg/driver"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Driver struct {
	mu      sync.Mutex
	clients map[string]*firestore.Client
}

func New() *Driver {
	return &Driver{clients: map[string]*firestore.Client{}}
}

func (d *Driver) Name() string { return "firestore" }

func (d *Driver) Methods() []string { return []string{"get"} }

func (d *Driver) Options(method string) []driver.Option {
	return []driver.Option{{Name: "path", Required: true}}
}

func (d *Driver) Execute(ctx context.Context, method string, options map[string]any) (driver.Output, error) {
	if method != "get" {
		return nil, driver.NewError(driver.ErrConfig, "firestore: unknown method %q", method)
	}

	projectID, _ := options["project_id"].(string)
	if projectID == "" {
		return nil, driver.NewError(driver.ErrConfig, "firestore: project_id is required")
	}

	credentialsFile, _ := options["credentials_file"].(string)
	if credentialsFile == "" {
		credentialsFile = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}
	if credentialsFile == "" {
		return nil, driver.NewError(driver.ErrConfig, "firestore: credentials_file is required")
	}

	databaseID, _ := options["database_id"].(string)
	if databaseID == "" {
		databaseID = "(default)"
	}

	path, _ := options["path"].(string)
	if path == "" {
		return nil, driver.NewError(driver.ErrConfig, "firestore: path is required")
	}
	if !isValidDocPath(path) {
		return nil, driver.NewError(driver.ErrConfig, "firestore: path %q must have an even number of non-empty segments (collection/doc/collection/doc/...)", path)
	}

	client, err := d.client(ctx, projectID, databaseID, credentialsFile)
	if err != nil {
		return nil, driver.Wrap(driver.ErrTransport, err, "firestore: create client: %v", err)
	}

	docSnap, err := client.Doc(path).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return driver.Output{"doc": nil, "exists": false}, nil
		}
		return nil, driver.Wrap(driver.ErrTransport, err, "firestore: get document: %v", err)
	}

	return driver.Output{
		"doc":    sanitize(docSnap.Data()),
		"exists": true,
	}, nil
}

func (d *Driver) client(ctx context.Context, projectID, databaseID, credentialsFile string) (*firestore.Client, error) {
	key := projectID + "|" + databaseID + "|" + credentialsFile

	d.mu.Lock()
	defer d.mu.Unlock()

	if c, ok := d.clients[key]; ok {
		return c, nil
	}

	c, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}
	d.clients[key] = c
	return c, nil
}

func isValidDocPath(path string) bool {
	segments := strings.Split(path, "/")
	if len(segments)%2 != 0 {
		return false
	}
	for _, s := range segments {
		if s == "" {
			return false
		}
	}
	return true
}

func sanitize(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			out[k] = sanitize(item)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = sanitize(item)
		}
		return out
	case *firestore.DocumentRef:
		if val == nil {
			return nil
		}
		return val.Path
	case *latlng.LatLng:
		if val == nil {
			return nil
		}
		return map[string]any{
			"latitude":  val.GetLatitude(),
			"longitude": val.GetLongitude(),
		}
	default:
		return v
	}
}
