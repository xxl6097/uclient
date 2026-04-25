package ntfy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

const (
	DefaultNtfyHost = "https://ntfy.sh"
	DefaultTimeout  = 10 * time.Second
	maxRetryBackoff = 30 * time.Second
)

// Client 支持：同步消息、停止、重试、BasicAuth
type Client struct {
	host      string
	timeout   time.Duration
	username  string
	password  string
	mu        sync.Mutex
	callbacks map[string]chan *Message

	ctx    context.Context
	cancel context.CancelFunc
	closed bool
	topic  string
}

// Message ntfy 消息结构
type Message struct {
	ID       string            `json:"id"`
	Topic    string            `json:"topic"`
	Message  string            `json:"message"`
	Title    string            `json:"title"`
	Time     int64             `json:"time"`
	Headers  map[string]string `json:"headers"`
	Expires  int64             `json:"expires,omitempty"`
	Event    string            `json:"event,omitempty"`
	Markdown bool              `json:"markdown,omitempty"`
}

type SyncMessage struct {
	ID            string `json:"id"`
	TargetTopic   string `json:"targetTopic"`
	ResponseTopic string `json:"responseTopic"`
	Data          string `json:"data"`
}

// NewClientWithAuth 创建带 BasicAuth 认证的客户端（生产用）
func NewClientWithAuth(host, username, password string, timeout time.Duration) *Client {
	if host == "" {
		host = DefaultNtfyHost
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		host:      host,
		timeout:   timeout,
		username:  username,
		password:  password,
		callbacks: make(map[string]chan *Message),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// NewClient 兼容旧接口（无认证）
func NewClient(host string, timeout time.Duration) *Client {
	return NewClientWithAuth(host, "", "", timeout)
}

// Stop 优雅停止
func (c *Client) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.cancel()
	for id, ch := range c.callbacks {
		close(ch)
		delete(c.callbacks, id)
	}
	c.closed = true
	fmt.Println("[ntfy] 客户端已停止")
}

// SendSync 同步发送（带认证 + respTopic） reqTopic, content string
func (c *Client) SendSync(title string, data *SyncMessage) (*Message, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil, errors.New("client stopped")
	}
	if c.topic != "" && data.ResponseTopic == "" {
		data.ResponseTopic = c.topic
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// 发送：携带响应主题头
	msgID, err := c.Send(&Message{
		Title:   title,
		Topic:   data.TargetTopic,
		Message: string(body),
	})
	if err != nil {
		c.cleanupCallback(msgID)
		return nil, err
	}

	respChan := make(chan *Message, 1)
	c.callbacks[msgID] = respChan
	c.mu.Unlock()

	// 等待响应
	select {
	case <-c.ctx.Done():
		c.cleanupCallback(msgID)
		return nil, errors.New("client stopped")
	case resp := <-respChan:
		c.cleanupCallback(msgID)
		return resp, nil
	case <-time.After(c.timeout):
		c.cleanupCallback(msgID)
		return nil, errors.New("timeout")
	}
}

// ListenWithRetry 带重试的监听
func (c *Client) ListenWithRetry(topic string, handler func(*Message) string) {
	baseDelay := 1 * time.Second
	attempt := 0

	for {
		select {
		case <-c.ctx.Done():
			fmt.Printf("[ntfy][%s] 监听停止\n", topic)
			return
		default:
		}

		fmt.Printf("[ntfy][%s] 开始监听\n", topic)
		err := c.listen(topic, handler)

		if errors.Is(err, context.Canceled) {
			return
		}

		attempt++
		delay := c.exponentialBackoff(baseDelay, attempt)
		fmt.Printf("[ntfy][%s] 断开: %v, %.1fs 后重试\n", topic, err, delay.Seconds())

		select {
		case <-time.After(delay):
		case <-c.ctx.Done():
			return
		}
	}
}

// listen 单次监听（带认证）
func (c *Client) listen(topic string, handler func(*Message) string) error {
	c.topic = topic
	url := fmt.Sprintf("%s/%s/json", c.host, topic)
	req, _ := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	c.setAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("监听成功 %s\n", url)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return context.Canceled
		default:
		}

		line := scanner.Text()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var msg Message
		if err = json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		if msg.Message != "" {
			z.L().Debug("aaaa", zap.Any("msg", msg))
		}
		go func(m Message) {
			result := handler(&m)
			if result != "" {
				var syncMsg SyncMessage
				err = json.Unmarshal([]byte(m.Message), &syncMsg)
				if err == nil && syncMsg.ResponseTopic != "" {
					tempData := SyncMessage{
						ID:   m.ID,
						Data: result,
					}
					jsonBytes, e := json.Marshal(tempData)
					if e != nil {
						return
					}
					_, _ = c.Send(&Message{
						Topic:   syncMsg.ResponseTopic,
						Message: string(jsonBytes),
					})

				}
				//z.L().Warn(result)
			}
		}(msg)
	}
	return scanner.Err()
}

func (c *Client) Send(data *Message) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return c.publish(body)
}

// publish 发布消息（带认证）
func (c *Client) publish(payload []byte) (string, error) {
	url := fmt.Sprintf("%s", c.host)
	req, _ := http.NewRequestWithContext(c.ctx, "POST", url, bytes.NewReader(payload))
	c.setAuthHeader(req)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 6. 读取并解析响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ntfy err] %s\n", err.Error())
		return "", err
	}
	var result Message
	if err = json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, 响应: %s", err, string(respBody))
	}

	fmt.Printf("[ntfy publish] %s\n", result.ID)
	return result.ID, nil
}

// setAuthHeader 自动添加 BasicAuth（核心）
func (c *Client) setAuthHeader(req *http.Request) {
	if c.username != "" && c.password != "" {
		auth := c.username + ":" + c.password
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", basic)
	}
}

// DispatchResponse 分发响应（给监听协程使用）
func (c *Client) DispatchResponse(msg *Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.callbacks[msg.ID]; ok {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (c *Client) cleanupCallback(msgID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.callbacks[msgID]; ok {
		close(ch)
		delete(c.callbacks, msgID)
	}
}

func (c *Client) exponentialBackoff(base time.Duration, attempt int) time.Duration {
	d := base * (1 << attempt)
	if d > maxRetryBackoff {
		d = maxRetryBackoff
	}
	jitter := time.Duration(rand.Int63n(int64(d) / 2))
	return d/2 + jitter
}
