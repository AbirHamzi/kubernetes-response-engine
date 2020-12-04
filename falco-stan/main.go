// Original file: https://github.com/falcosecurity/kubernetes-response-engine/blob/master/falco-nats/main.go

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"regexp"
	"strings"
)

var slugRegularExpression = regexp.MustCompile("[^a-z0-9]+")

func main() {

	var urls = flag.String("s", "nats://nats.nats.svc:4222", "The stan server URLs (separated by comma)")
	var pipePath = flag.String("f", "/tmp/shared-pipe/nats", "The named pipe path")
	var clientID = flag.String("p", "clientID", "The name of the pod that will be used as a unique clientID")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	// Stan is a client to a NATS server, Uncomment the following to reuse
	// a NATS connection instead of creating a new one.
	/*	nc, err := nats.Connect(*urls)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
	*/

	sc, err := stan.Connect("stan", *clientID, stan.NatsURL(*urls))

	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, *urls)
	}
	defer sc.Close()

	pipe, err := os.OpenFile(*pipePath, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Opened pipe %s", *pipePath)

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	log.Printf("Scanning %s", *pipePath)

	for scanner.Scan() {
		msg := []byte(scanner.Text())

		subj, err := subjectAndRuleSlug(msg)
		if err != nil {
			log.Fatal(err)
		}

		err = sc.Publish(subj, msg)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Published [%s] : '%s'\n", subj, msg)
		}
	}
}

func usage() {
	log.Fatalf("Usage: falco-stan [-s server (%s)] <subject> <msg> \n", stan.DefaultNatsURL)
}

type parsedAlert struct {
	Priority string `json:"priority"`
	Rule     string `json:"rule"`
}

func subjectAndRuleSlug(alert []byte) (string, error) {
	var result parsedAlert
	err := json.Unmarshal(alert, &result)

	if err != nil {
		return "", err
	}

	subject := "falco." + result.Priority + "." + slugify(result.Rule)
	subject = strings.ToLower(subject)

	return subject, nil
}

func slugify(input string) string {
	return strings.Trim(slugRegularExpression.ReplaceAllString(strings.ToLower(input), "_"), "_")
}
