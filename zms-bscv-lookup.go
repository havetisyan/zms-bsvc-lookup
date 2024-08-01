package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/AthenZ/athenz/clients/go/zms"
)

func main() {
	var (
		pDomain = flag.String("d", "", "domain names")
		pKey    = flag.String("k", "", "service private key path")
		pCert   = flag.String("c", "", "service x.509 certificate path")
		pZmsUrl = flag.String("zms", "", "zms url")
	)

	flag.Parse()
	if *pZmsUrl == "" || *pDomain == "" || *pKey == "" || *pCert == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	zmsClient, err := NewZmsClient(*pZmsUrl, *pCert, *pKey)
	if err != nil {
		fmt.Printf("unable to create zms client: %v\n", err)
		os.Exit(1)
	}

	topLevelDomain := *pDomain
	idx := strings.Index(*pDomain, ".")
	if idx != -1 {
		topLevelDomain = topLevelDomain[:idx]
	}
	zmsDomain, err := zmsClient.GetDomain(zms.DomainName(topLevelDomain))
	if err != nil {
		fmt.Printf("unable to get domain details: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s,%s\n", topLevelDomain, zmsDomain.BusinessService)
}

func NewZmsClient(url string, certFile string, keyFile string) (*zms.ZMSClient, error) {
	tlsConfig, err := getTLSConfigFromFiles(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	transport := http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &zms.ZMSClient{
		URL:       url,
		Transport: &transport,
	}
	return client, err
}

func getTLSConfigFromFiles(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to formulate clientCert from key and certs bytes, error: %v", err)
	}

	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = cert

	// Set Renegotiation explicitly
	config.Renegotiation = tls.RenegotiateOnceAsClient

	return config, err
}
