package daemon

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
)

func CheckDataPoolExist(datapoolname string) (bexist bool) {
	sqlcheck := fmt.Sprintf("SELECT COUNT(1) FROM DH_DP WHERE DPNAME='%s'", datapoolname)
	row, err := g_ds.QueryRow(sqlcheck)
	if err != nil {
		fmt.Println("CheckDataPoolExist QueryRow error:", err.Error())
		return
	} else {
		var num int
		row.Scan(&num)
		if num == 0 {
			return false
		} else {
			return true
		}
	}
}

func GetDataPoolDpconn(datapoolname string) (dpconn string) {
	sqlgetdpconn := fmt.Sprintf("SELECT DPCONN FROM DH_DP WHERE DPNAME='%s'", datapoolname)
	row, err := g_ds.QueryRow(sqlgetdpconn)
	if err != nil {
		fmt.Println("GetDataPoolDpconn QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpconn)
		return dpconn
	}
}

func InsertTagToDb(dpexist bool, p ds.DsPull) (err error) {
	if dpexist == false {
		return
	}

	rpdmid, dpid := GetRepoItemId(p.Repository, p.Dataitem)
	if rpdmid == 0 {
		sqlInsertRpdm := fmt.Sprintf(`INSERT INTO DH_DP_RPDM_MAP(RPDMID ,REPOSITORY , DATAITEM, 
        	DPID  , PUBLISH ,CREATE_TIME ) VALUES (null, '%s', '%s', %d, 'N', datetime('now'))`,
			p.Repository, p.Dataitem, dpid)
		g_ds.Insert(sqlInsertRpdm)
		rpdmid, dpid = GetRepoItemId(p.Repository, p.Dataitem)
	}
	sqlInsertTag := fmt.Sprintf(`INSERT INTO DH_RPDM_TAG_MAP(TAGNAME ,RPDMID ,DETAIL,CREATE_TIME) 
		VALUES ('%s', '%d', '%s', datetime('now'))`,
		p.Tag, rpdmid, p.DestName)
	g_ds.Insert(sqlInsertTag)
	return err
}

func GetRepoItemId(repository, dataitem string) (rpdmid, dpid int) {
	sqlgetrpdmId := fmt.Sprintf("SELECT RPDMID ,DPID FROM DH_DP_RPDM_MAP WHERE REPOSITORY='%s' AND DATAITEM='%s'",
		repository, dataitem)
	row, err := g_ds.QueryRow(sqlgetrpdmId)
	if err != nil {
		fmt.Println("GetRepoItemId QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&rpdmid, &dpid)
		return
	}
}
