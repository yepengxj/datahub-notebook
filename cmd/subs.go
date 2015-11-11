package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
	"os"
)

func Subs(login bool, args []string) (err error) {

	itemDetail := false
	if len(args) > 1 {
		fmt.Println("invalid argument..")
		subsUsage()
		return
	}

	uri := "/subscriptions"
	if len(args) == 1 {
		uri = "/repositories"
		uri = uri + "/" + args[0]
		itemDetail = true
	}

	resp, err := commToDaemon("GET", uri, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		subsResp(itemDetail, body)
	} else if resp.StatusCode == 401 {
		if err := Login(false, nil); err == nil {
			Subs(login, args)
		} else {
			//fmt.Println(string(body))
			//fmt.Println(resp.StatusCode ,err)
			fmt.Println(err)
		}
	} else {
		fmt.Println(string(body))
	}

	return err
}

func subsUsage() {
	fmt.Printf("usage: %s subs [[URL]/[REPO]/[ITEM]\n", os.Args[0])
}

func subsResp(detail bool, respbody []byte) {

	if detail {
		subs := ds.Data{}
		err := json.Unmarshal(respbody, &subs)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s\t%s\n", "REPOSITORY/ITEM[:TAG]", "UPDATETIME")
		printDash(n + 12)
		for _, tag := range subs.Tags {
			fmt.Printf("%s/%s:%-8s\t%s\n", subs.Item.Repository_name, subs.Item.Dataitem_name, tag.Tag, tag.Optime)
		}
	} else {
		subs := []ds.Data{}
		err := json.Unmarshal(respbody, &subs)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s/%-8s\t%s\n", "REPOSITORY", "ITEM", "TYPE")
		printDash(n + 5)
		for _, item := range subs {
			fmt.Printf("%s/%-8s\t%s\n", item.Item.Repository_name, item.Item.Dataitem_name, "File")
		}

	}

}
