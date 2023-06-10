package contentorigin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client interacts with the backend to manage media
type Client struct {
	sftpc *sftp.Client
	conn  *ssh.Client
}

// New returns a new client connected with the given parameters
func New(host, port, user, pass string) (*Client, error) {
	// TODO: validation
	hostKey := getHostKey(host)
	cconfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	hostport := fmt.Sprintf("%s:%s", host, port)
	conn, err := ssh.Dial("tcp", hostport, cconfig)
	if err != nil {
		return nil, fmt.Errorf("contentorigin.New: failed to dial remote: %w", err)
	}

	// create new SFTP client
	// TODO(wishfulthinking): make this transport agnostic
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("contentorigin.New: failed create new sftp client: %w", err)
	}

	return &Client{sftpc: client, conn: conn}, nil
}

// Put uploads bytes from src to a dst path on the remote
func (c *Client) Put(srcFile io.Reader, dst string) (int64, error) {
	// create destination file
	dstFile, err := c.sftpc.Create(dst)
	if err != nil {
		return -1, fmt.Errorf("contentorigin.Put: failed to create %s: %w", dst, err)
	}
	defer dstFile.Close()

	// copy source file to destination file
	num, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return -1, fmt.Errorf("contentorigin.Put: copy to %s failed: %w", dst, err)
	}

	return num, nil
}

// List walks the remote starting at root and returns all paths found
func (c *Client) List(root string) ([]string, error) {
	paths := []string{}
	walker := c.sftpc.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return paths, err
		}
		paths = append(paths, walker.Path())
	}
	return paths, nil
}

// Walk walks the remote starting at root and calls func on all paths found. If an error is encountered it is returned immediately.
func (c *Client) Walk(root string, fn func(p string) error) error {
	walker := c.sftpc.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}
		if err := fn(walker.Path()); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the connection to remote
func (c *Client) Close() error {
	if err := c.sftpc.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
