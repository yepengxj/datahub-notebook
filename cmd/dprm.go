package cmd

import (
	"fmt"
)

func DpRm(needLogin bool, args []string) (err error) {
    var strdpname string
	if needLogin && !Logged {
		login(false)
	}
    if len(args) > 0 && args[0][0] != '-' {
		for _, v := range args {
		    strdpname = v
		    if v[0] != '-' {
		        str_dpurl := fmt.Sprintf("/datapool/%s", strdpname)

		        ret_body := commToDaemon(str_dpurl, "DELETE", nil)
	            fmt.Println(string(ret_body))
		    }
		}	
	}
	return nil
}
