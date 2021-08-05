package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	checkIPURL = "https://domains.google.com/checkip"
)

type Config struct {
	OwnerEmail string `envconfig:"EMAIL" required:"true"`
	Username   string `envconfig:"USERNAME" required:"true"`
	Password   string `envconfig:"PASSWORD" required:"true"`
	RequestURL string `envconfig:"REQUEST_URL" required:"true"`
	Hostname   string `envconfig:"HOSTNAME" required:"true"`
}

func main() {
	var env Config
	envconfig.MustProcess("", &env)

	ip, err := getIP()
	if err != nil {
		log.Fatalf("failed to get IP: %v", err)
	}

	if alreadySet(env.Hostname, ip) {
		log.Println("Hostname is already up to date.")
		return
	}

	err = setIP(env, ip)
	if err != nil {
		log.Fatalf("failed to update IP: %v", err)
	}
	log.Println("Success.")
}

func getIP() (net.IP, error) {
	resp, err := http.Get(checkIPURL)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	raw := buf.String()
	log.Printf("fetched IP %q", raw)
	ip := net.ParseIP(raw)
	if ip == nil {
		return nil, fmt.Errorf("no IP found in %q", raw)
	}
	return ip, nil
}

func setIP(env Config, ip net.IP) error {
	addr, err := url.Parse(env.RequestURL)
	if err != nil {
		log.Fatalf("failed to parse %q: %v", env.RequestURL, err)
	}
	args := make(url.Values)
	args.Add("hostname", env.Hostname)
	args.Add("myip", ip.String())
	addr.RawQuery = args.Encode()

	const agentFmt = "github.com/spencer-p/dynamic-dns %s"
	header := make(http.Header)
	header.Add("Agent", fmt.Sprintf(agentFmt, env.OwnerEmail))
	request := &http.Request{
		Method: "POST",
		Header: header,
		URL:    addr,
	}
	request.SetBasicAuth(env.Username, env.Password)
	os.Stdout.Write([]byte{'\n'})
	request.Write(os.Stdout)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("response is %q", buf.String())
		return fmt.Errorf("got status %s", resp.Status)
	}

	switch buf.String() {
	case "good " + ip.String():
		log.Println("Successfully updated IP.")
	case "nochg " + ip.String():
		log.Println("IP address not changed.")
	default:
		return fmt.Errorf("error response: %s", buf.String())
	}

	return nil
}

func alreadySet(hostname string, ip net.IP) bool {
	all, err := net.LookupIP(hostname)
	if err != nil {
		log.Println("failed to look up %s: %v", hostname, err)
		return false
	}
	for i := range all {
		if ip.Equal(all[i]) {
			return true
		}
	}
	return false
}
