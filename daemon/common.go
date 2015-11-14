package daemon

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
)

func CheckDataPoolExist(datapoolname string) (bexist bool) {
	sqlcheck := fmt.Sprintf("SELECT COUNT(1) FROM DH_DP WHERE DPNAME='%s'", datapoolname)
	row, err := g_ds.QueryRow(sqlcheck)
	//fmt.Println(sqlcheck)
	if err != nil {
		fmt.Println("CheckDataPoolExist QueryRow error:", err.Error())
		return
	} else {
		var num int
		row.Scan(&num)
		//fmt.Println("num:", num)
		if num == 0 {
			return false
		} else {
			return true
		}
	}
}

func GetDataPoolDpconn(datapoolname string) (dpconn string) {
	sqlgetdpconn := fmt.Sprintf("SELECT DPCONN FROM DH_DP WHERE DPNAME='%s'", datapoolname)
	//fmt.Println(sqlgetdpconn)
	row, err := g_ds.QueryRow(sqlgetdpconn)
	if err != nil {
		fmt.Println("GetDataPoolDpconn QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpconn)
		return dpconn
	}
}

func GetDataPoolDpid(datapoolname string) (dpid int) {
	sqlgetdpid := fmt.Sprintf("SELECT DPID FROM DH_DP WHERE DPNAME='%s'", datapoolname)
	//fmt.Println(sqlgetdpid)
	row, err := g_ds.QueryRow(sqlgetdpid)
	if err != nil {
		fmt.Println("GetDataPoolDpid QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpid)
		return
	}
}

func InsertTagToDb(dpexist bool, p ds.DsPull) (err error) {
	if dpexist == false {
		return
	}
	DpId := GetDataPoolDpid(p.Datapool)
	if DpId == 0 {
		return
	}
	rpdmid := GetRepoItemId(p.Repository, p.Dataitem)
	//fmt.Println("GetRepoItemId1", rpdmid, DpId)
	if rpdmid == 0 {
		sqlInsertRpdm := fmt.Sprintf(`INSERT INTO DH_DP_RPDM_MAP(RPDMID ,REPOSITORY , DATAITEM, 
        	DPID  , PUBLISH ,CREATE_TIME ) VALUES (null, '%s', '%s', %d, 'N', datetime('now'))`,
			p.Repository, p.Dataitem, DpId)
		g_ds.Insert(sqlInsertRpdm)
		rpdmid = GetRepoItemId(p.Repository, p.Dataitem)
		//fmt.Println("GetRepoItemId2", rpdmid, DpId)
	}
	sqlInsertTag := fmt.Sprintf(`INSERT INTO DH_RPDM_TAG_MAP(TAGNAME ,RPDMID ,DETAIL,CREATE_TIME) 
		VALUES ('%s', '%d', '%s', datetime('now'))`,
		p.Tag, rpdmid, p.DestName)
	fmt.Println(sqlInsertTag)
	g_ds.Insert(sqlInsertTag)
	return err
}

func GetRepoItemId(repository, dataitem string) (rpdmid int) {
	sqlgetrpdmId := fmt.Sprintf("SELECT RPDMID FROM DH_DP_RPDM_MAP WHERE REPOSITORY='%s' AND DATAITEM='%s'",
		repository, dataitem)
	row, err := g_ds.QueryRow(sqlgetrpdmId)
	if err != nil {
		fmt.Println("GetRepoItemId QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&rpdmid)
		return
	}
}
