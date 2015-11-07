package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
)

func dpPostOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	rw.WriteHeader(http.StatusOK)

	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		reqJson := cmd.FormatDpCreate{}
		err := json.Unmarshal(result, &reqJson)
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
			fmt.Println(rw, "invalid argument.")
		}
		if len(reqJson.Name) == 0 {
			fmt.Fprintln(rw, "invalid argument.")
		} else {
			msg := &ds.MsgResp{}
			var sdpDirName string
			if len(reqJson.Conn) == 0 {
				sdpDirName = g_strDpPath + reqJson.Name
				reqJson.Conn = sdpDirName

			} else if reqJson.Conn[0] != '/' {
				sdpDirName = g_strDpPath + reqJson.Conn
				reqJson.Conn = sdpDirName
				sdpDirName = sdpDirName + "/" + reqJson.Name
			} else {
				sdpDirName = reqJson.Conn + "/" + reqJson.Name
			}

			if err := os.MkdirAll(sdpDirName, 0755); err != nil {
				msg.Msg = err.Error()
			} else {
				msg.Msg = fmt.Sprintf("OK. dp:%s total path:%s", reqJson.Name, sdpDirName)
				sql_dp_insert := fmt.Sprintf(`insert into DH_DP (DPID, DPNAME, DPTYPE, DPCONN, STATUS)
					values (null, '%s', '%s', '%s', 'A')`, reqJson.Name, reqJson.Type, reqJson.Conn)
				if _, err := g_ds.Insert(sql_dp_insert); err != nil {
					os.Remove(sdpDirName)
					msg.Msg = err.Error()
				}
			}
			resp, _ := json.Marshal(msg)
			respStr := string(resp)
			fmt.Fprintln(rw, respStr)
		}

	}

}

func dpGetAllHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	result, _ := ioutil.ReadAll(r.Body)
	reqJson := cmd.FormatDp{}
	err := json.Unmarshal(result, &reqJson)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		fmt.Println(rw, "invalid argument.")
	}

	msg := &ds.MsgResp{}

	msg.Msg = "OK."
	wrtocli := cmd.FormatDp{}
	sql_dp := fmt.Sprintf(`SELECT DPNAME, DPTYPE FROM DH_DP WHERE STATUS = 'A'`)
	rows, err := g_ds.QueryRows(sql_dp)
	if err != nil {
		msg.Msg = err.Error()
	}
	defer rows.Close()
	bresultflag := false
	for rows.Next() {
		bresultflag = true
		rows.Scan(&wrtocli.Name, &wrtocli.Type)
		resp, _ := json.Marshal(wrtocli)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}
	if bresultflag == false {
		msg.Msg = "There isn't any datapool."
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}

}

func dpGetOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)
	dpname := ps.ByName("dpname")

	//In future, we need to get dptype in Json to surpport FILE\ DB\ SDK\ API datapool
	result, _ := ioutil.ReadAll(r.Body)
	reqJson := cmd.FormatDp{}
	err := json.Unmarshal(result, &reqJson)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		fmt.Println(rw, "invalid argument.")
	}

	msg := &ds.MsgResp{}
	msg.Msg = "OK."

	sql_total := fmt.Sprintf(`SELECT COUNT(*) FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, string(dpname))
	row, err := g_ds.QueryRow(sql_total)
	if err != nil {
		msg.Msg = err.Error()
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
		return
	}
	var total int
	row.Scan(&total)
	if total == 0 {
		msg.Msg = fmt.Sprintf("Datapool %v not found.", dpname)
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
		return
	}

	sql_dp := fmt.Sprintf(`SELECT DPID, DPNAME, DPTYPE, DPCONN FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, string(dpname))
	rows, err := g_ds.QueryRows(sql_dp)
	if err != nil {
		msg.Msg = err.Error()
	}
	defer rows.Close()
	var dpid int
	//I use queryrows because  there will be distinct dataitems in the same datapool
	for rows.Next() {
		onedp := cmd.FormatDp_dpname{}
		onedp.Items = make([]cmd.Item, 0, 16)
		rows.Scan(&dpid, &onedp.Name, &onedp.Type, &onedp.Conn)
		if dpid > 0 {
			//Use "left out join" to get repository/dataitem records, whether it has tags or not.
			sql_tag := fmt.Sprintf(`SELECT A.REPOSITORY, A.DATAITEM, B.TAGNAME, strftime(B.CREATE_TIME), A.PUBLISH 
				FROM DH_DP_RPDM_MAP A LEFT JOIN DH_RPDM_TAG_MAP B
				ON (A.RPDMID = B.RPDMID)
				WHERE A.DPID = %v`, dpid)
			tagrows, err := g_ds.QueryRows(sql_tag)
			if err != nil {
				msg.Msg = err.Error()
			}

			fmt.Println(msg)
			defer tagrows.Close()
			for tagrows.Next() {
				fmt.Println(tagrows)
				item := cmd.Item{}
				tagrows.Scan(&item.Repository, &item.DataItem, &item.Tag, &item.Time, &item.Publish)
				onedp.Items = append(onedp.Items, item)
			}
		}

		resp, _ := json.Marshal(onedp)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}

}

func dpDeleteOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	dpname := ps.ByName("dpname")
	msg := &ds.MsgResp{}

	sql_dp_rm := fmt.Sprintf(`SELECT DPID, DPTYPE FROM DH_DP WHERE DPNAME ='%s'`, dpname)
	dprows, err := g_ds.QueryRows(sql_dp_rm)
	if err != nil {
		msg.Msg = err.Error()
	}

	bresultflag := false

	dpid_type := make([]strc_dp, 0, 8)
	strcone := strc_dp{}
	for dprows.Next() {
		dprows.Scan(&strcone.Dpid, &strcone.Dptype)
		dpid_type = append(dpid_type, strcone)
	}
	dprows.Close()

	for _, v := range dpid_type {
		var dpid = v.Dpid
		var dptype = v.Dptype
		bresultflag = true
		//dprow.Scan(&dpid, &dptype)
		sql_dp_item := fmt.Sprintf("SELECT PUBLISH FROM DH_DP_RPDM_MAP WHERE DPID = %v ", dpid)
		row, err := g_ds.QueryRow(sql_dp_item)
		if err != nil {
			msg.Msg = err.Error()
		}
		//time.Sleep(60*time.Second)
		var sPublish string
		row.Scan(&sPublish)
		if sPublish == "Y" {
			msg.Msg = fmt.Sprintf(`Datapool %s with type:%s can't be removed , it contains published DataItem !`, dpname, dptype)
		} else {
			sql_update := fmt.Sprintf("UPDATE DH_DP SET STATUS = 'N' WHERE DPID = %v", dpid)
			_, err := g_ds.Update(sql_update)
			if err != nil {
				msg.Msg = err.Error()
			} else {
				msg.Msg = fmt.Sprintf("Datapool %s with type:%s removed successfully!", dpname, dptype)
			}
		}
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}
	if bresultflag == false {
		msg.Msg = fmt.Sprintf("Datapool %s not found.\n", dpname)
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}
}
