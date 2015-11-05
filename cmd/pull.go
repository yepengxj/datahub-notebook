package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Pull(login bool, args []string) (err error) {

	if len(args) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
		return
	}
	u, err := url.Parse(args[0])
	if err != nil {
		panic(err)
	}
	source := u.Path
	var uri string
	if u.Path[0] == '/' {
		source = u.Path[1:]
	}

	if url := strings.Split(source, "/"); len(url) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
		return
	} else {
		target := strings.Split(url[1], ":")
		if len(target) == 1 {
			target = append(target, "latest")
		} else if len(target[1]) == 0 {
			target[1] = "latest"
		}
		uri = fmt.Sprintf("%s/%s:%s", url[0], target[0], target[1])
	}

	fmt.Println(uri)

	return dl(uri)
	//return nil
}

func pullUsage() {
	fmt.Printf("usage: %s pull [[URL]/[REPO]/[ITEM][:TAG]] [DATAPOOL]\n", os.Args[0])
}

func dl(uri string) error {
	ip := os.Getenv("DAEMON_IP_PEER")
	if len(ip) == 0 {
		ip = "http://127.0.0.1:35800"
	}

	target := ip + "/pull?target=" + uri
	fmt.Println(target)
	n, err := download(target)
	if err != nil {
		fmt.Printf("[%d bytes returned.]\n", n)
		fmt.Println(err)
	}
	return err
}

/*download routine, supports resuming broken downloads.*/
func download(url string) (int64, error) {
	fmt.Printf("we are going to download %s\n", url)
	out, err := os.OpenFile("/var/lib/datahub/test", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	stat, err := out.Stat()
	if err != nil {
		return 0, err
	}
	out.Seek(stat.Size(), 0)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "go-downloader")
	/* Set download starting position with 'Range' in HTTP header*/
	req.Header.Set("Range", "bytes="+strconv.FormatInt(stat.Size(), 10)+"-")
	fmt.Printf("%v bytes had already been downloaded.\n", stat.Size())

	resp, err := http.DefaultClient.Do(req)

	/*Save response body to file only when HTTP 2xx received. TODO*/
	if err != nil || (resp != nil && resp.StatusCode/100 != 2) {
		if resp != nil {
			fmt.Println("http status code:", resp.StatusCode, err)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
		}
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}
	fmt.Printf("%d bytes downloaded.", n)
	return n, nil
}
