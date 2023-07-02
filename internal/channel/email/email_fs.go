package email

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
)

type FS struct {
	path string
}

func NewFS(path string) *FS {
	return &FS{path: path}
}

func (m *FS) Send(ctx context.Context, target channel.Target, subject, body string) error {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("To: %s\n", target.Address))
	s.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	s.WriteString(fmt.Sprintf("Body: \n%s\n", body))

	p := fmt.Sprintf("%s/meetup-%s-%s.txt", m.path, target.Address, time.Now().Format("20060102150405"))
	return os.WriteFile(p, []byte(s.String()), 0644)
}
