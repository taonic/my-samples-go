package proxied_activities

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
)

type ProxiedActivity struct {
	conn net.Conn
	pid  int
}

func NewProxyActivity(pid int) *ProxiedActivity {
	conn, err := net.Dial("tcp", "localhost:9992")
	if err != nil {
		return nil
	}
	return &ProxiedActivity{conn: conn, pid: pid}
}

func (p *ProxiedActivity) ProxyActivity(ctx context.Context, message string) (string, error) {
	greeting := fmt.Sprintf("(%d): %s", p.pid, message)
	_, err := p.conn.Write([]byte(greeting + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	reader := bufio.NewReader(p.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return fmt.Sprintf("%s (%d)", strings.TrimSpace(response), p.pid), nil
}
