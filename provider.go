package blubber

import (
	"net"
	"regexp"
	"strconv"
	"io/ioutil"
	"path/filepath"
    "gopkg.in/yaml.v2"
)

type providerYaml struct {
	Public		  bool
	Responses     map[string]struct{
		Description string
		Score string
	}
}

type Score struct {
	OpenRelay int
	OpenProxy int
	Compromised int
	Spam int
	Reputation int
	Other int
}

type Response struct {
	Addr net.IP
	Description string
	Score Score
}

type Provider struct {
	Hostname	string
	Public		bool
	ResponseList []Response
}

func strToScore(str string) Score {
	re := regexp.MustCompile("(?P<key>[A-Za-z]*):(?P<value>-?\\d*)")

	s := Score{}
	for _,e := range re.FindAllStringSubmatch(str, -1) {
		key := e[1]
		value,_ := strconv.Atoi(e[2])

		switch key {
		case "OpenRelay":
			s.OpenRelay += value
		case "OpenProxy":
			s.OpenProxy += value
		case "Compromised":
			s.Compromised += value
		case "Spam":
			s.Spam += value
		case "Reputation":
			s.Reputation += value
		default:
			s.Other += value
		}
	}

	return s
}

func yamlToProvider(hostname string, yaml providerYaml) Provider {
	p := Provider{
		Hostname: hostname,
		Public: yaml.Public,
	}
	for k,v := range yaml.Responses {
		r := Response{
			Addr: net.ParseIP(k),
			Description: v.Description,
			Score: strToScore(v.Score),
		}
		p.ResponseList = append(p.ResponseList, r)
	}

	return p
}

func (bl *Blubber)LoadProvidersFromFile(file string) error {
	var providerYamlMap map[string]providerYaml

	filename, _ := filepath.Abs(file)

	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &providerYamlMap)
	if err != nil {
		return err
	}

	for k,v := range providerYamlMap {
		bl.ProviderList = append(bl.ProviderList, yamlToProvider(k,v))
	}

	return nil
}