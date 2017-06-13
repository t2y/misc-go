package main

import (
	"flag"
	"log"
	"os"

	"github.com/ncw/swift"
)

var user string
var key string
var url string
var domain string
var tenant string

func init() {
	flag.StringVar(&user, "user", "", "")
	flag.StringVar(&key, "key", "", "")
	flag.StringVar(&url, "url", "", "")
	flag.StringVar(&domain, "domain", "", "")
	flag.StringVar(&tenant, "tenant", "", "")
}

func createConn(user, key, url, domain, tenant string) (c *swift.Connection, err error) {
	c = &swift.Connection{
		UserName: user,
		ApiKey:   key,
		AuthUrl:  url,
		Domain:   domain,
		Tenant:   tenant,
	}

	if err = c.Authenticate(); err != nil {
		return
	}

	return
}

func ListContainers(c *swift.Connection, limit int, prefix string) (err error) {
	opts := &swift.ContainersOpts{
		Limit:  limit,
		Prefix: prefix,
	}

	var list []swift.Container
	if list, err = c.Containers(opts); err != nil {
		return
	}

	for _, v := range list {
		log.Printf("%s,%d,%d\n", v.Name, v.Count, v.Bytes)
	}
	return
}

func ListObjects(c *swift.Connection, container string, limit int, prefix string) (err error) {
	opts := &swift.ObjectsOpts{
		Limit:  limit,
		Prefix: prefix,
	}

	var list []swift.Object
	if list, err = c.Objects(container, opts); err != nil {
		return
	}

	for _, v := range list {
		log.Printf("%s,%s,%d,%d\n", v.Name, v.Hash, v.Bytes, v.LastModified.Unix())
	}
	return
}

func main() {
	flag.Parse()

	c, err := createConn(user, key, url, domain, tenant)
	if err != nil {
		log.Fatal(err)
	}

	if err := ListContainers(c, 3, ""); err != nil {
		log.Fatal(err)
	}

	containerName := "mycontainer1"
	if err := c.ContainerCreate(containerName, nil); err != nil {
		log.Fatal(err)
	}

	if err := ListContainers(c, 3, ""); err != nil {
		log.Fatal(err)
	}

	testTxt := "test.txt"
	f, err := os.Open(testTxt)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	headerPut, err := c.ObjectPut(containerName, testTxt, f, false, "", "text", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(headerPut)

	if err := ListObjects(c, containerName, 3, ""); err != nil {
		log.Fatal(err)
	}

	f = os.Stdout
	headerGet, err := c.ObjectGet(containerName, testTxt, f, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(headerGet)
}
