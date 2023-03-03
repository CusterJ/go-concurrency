package links

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"
)

func CreateFileWithLinks() {

}

// Generates a random slice of HTTP and HTTPS links of a given length from the domains.txt file
func GenerateLinks(count int) []string {
	rand.Seed(time.Now().UnixNano())

	var links []string
	protocol := []string{"http://", "https://"}
	domains := GetDomains()

	for i := 0; i < count; i++ {
		p := protocol[rand.Intn(len(protocol))]
		d := domains[rand.Intn(len(domains))]
		links = append(links, p+d)
	}
	return links
}

func GetDomains() []string {
	var domains []string

	file, err := os.Open("links/domains.txt")
	if err != nil {
		log.Println("GenerateLinks -> can't open domains.txt file")
		return domains
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}
	log.Println("domains len is: ", len(domains))
	return domains
}
