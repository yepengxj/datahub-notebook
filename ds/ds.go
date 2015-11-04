package ds

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type FormatRepoList struct {
	RepoName string `json:"repoName"`
	ItemID   string `json:"itemID"`
}

const (
	DB_DML_INSERT = "insert"
	DB_DML_DELETE = "delete"
	DB_DML_UPDATE = "update"
	DB_DML_SELECT = "select"
	DB_DDL_CREATE = "create"
	DB_DDL_DROP   = "drop"
	TABLE_ORDER   = "order_t"
	TABLE_USER    = "user"
)

type Data struct {
	Item  DataItem
	Usage DataItemUsage
}

type Repository struct {
	Repository_id   int    `json:"repository_id,omitempty"`
	Repository_name string `json:"repository_name,omitempty"`
	User_id         int    `json:"user_id,omitempty"`
	Permit_type     int    `json:"permit_type,omitempty"`
	Arrange_type    int    `json:"arrange_type,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Rank            int    `json:"rank,omitempty"`
	Status          int    `json:"status,omitempty"`
	Dataitems       int    `json:"dataitems,omitempty"`
	Tags            int    `json:"tags,omitempty"`
	Stars           int    `json:"stars,omitempty"`
	Optime          string `json:"optime,omitempty"`
}
type DataItem struct {
	Repository_id   int     `json:"repository_id,omitempty"`
	User_id         int     `json:"user_id,omitempty"`
	Dataitem_id     int     `json:"dataitem_id,omitempty"`
	Dataitem_name   string  `json:"dataitem_name,omitempty"`
	Ico_name        string  `json:"ico_name,omitempty"`
	Permit_type     int     `json:"permit_type,omitempty"`
	Key_words       string  `json:"key_words,omitempty"`
	Supply_style    int     `json:"supply_style,omitempty"`
	Priceunit_type  int     `json:"priceunit_type,omitempty"`
	Price           float32 `json:"price,omitempty"`
	Optime          string  `json:"optime,omitempty"`
	Data_format     int     `json:"data_format,omitempty"`
	Refresh_type    int     `json:"refresh_type,omitempty"`
	Refresh_num     int     `json:"refresh_num,omitempty"`
	Meta_filename   string  `json:"meta_filename,omitempty"`
	Sample_filename string  `json:"sample_filename,omitempty"`
	Comment         string  `json:"comment,omitempty"`
}

type DataItemUsage struct {
	Dataitem_id   int    `json:"-,omitempty"`
	Dataitem_name string `json:"-,omitempty"`
	Views         int    `json:"views"`
	Follows       int    `json:"follows"`
	Downloads     int    `json:"downloads"`
	Stars         int    `json:"stars"`
	Refresh_date  string `json:"refresh_date,omitempty"`
	Usability     int    `json:"usability,omitempty"`
}

type UpLoadLog struct {
	Tag      string `json:"tag,omitempty"`
	Filename string `json:"filename,omitempty"`
	Optime   string `json:"optime,omitempty"`
}

type Ds struct {
	Db *sql.DB
}

const Create_dh_dp string = `CREATE TABLE IF NOT EXISTS 
    DH_DP ( 
       DPID    INTEGER PRIMARY KEY AUTOINCREMENT, 
       DPNAME  VARCHAR(32), 
       DPTYPE  VARCHAR(32), 
       DPCONN  VARCHAR(128), 
       STATUS  CHAR(2) 
    );`

//DH_DP STATUS : 'A' valid; 'N' invalid; 'P' contain dataitem published;

const Create_dh_dp_repo_ditem_map string = `CREATE TABLE IF NOT EXISTS 
    DH_DP_RPDM_MAP ( 
    	RPDMID       INTEGER PRIMARY KEY AUTOINCREMENT, 
        REPOSITORY   VARCHAR(128), 
        DATAITEM     VARCHAR(128), 
        DPID         INTEGER, 
        PUBLISH      CHAR(2), 
        CREATE_DATE  DATE 
    );`

//DH_DP_REPO_DITEM_MAP  PUBLISH: 'Y' the dataitem is published by you,
//'N' the dataitem is pulled by you

const Create_dh_repo_ditem_tag_map string = `CREATE TABLE IF NOT EXISTS 
    DH_RPDM_TAG_MAP ( 
        TAGNAME      VARCHAR(128), 
        RPDMID       INTEGER 
    );`

type Executer interface {
	Insert(cmd string) (interface{}, error)
	Delete(cmd string) (interface{}, error)
	Update(cmd string) (interface{}, error)
	QueryRaw(cmd string) (*sql.Rows, error)
	QueryRaws(cmd string) (*sql.Rows, error)

	Create(cmd string) (interface{}, error)
	Drop(cmd string) (interface{}, error)
}

func execute(p *Ds, cmd string) (interface{}, error) {
	tx, err := p.Db.Begin()
	if err != nil {
		return nil, err
	}
	var res sql.Result
	if res, err = tx.Exec(cmd); err != nil {
		log.Printf(`Exec("%s") err %s`, cmd, err.Error())
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, nil
}

func query(p *Ds, cmd string) (*sql.Row, error) {
	return p.Db.QueryRow(cmd), nil
}
func queryRows(p *Ds, cmd string) (*sql.Rows, error) {
	return p.Db.Query(cmd)
}

func (p *Ds) Insert(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Delete(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Update(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) QueryRow(cmd string) (*sql.Row, error) {
	return query(p, cmd)
}

func (p *Ds) QueryRows(cmd string) (*sql.Rows, error) {
	return queryRows(p, cmd)
}
func (p *Ds) Create(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Drop(cmd string) (interface{}, error) {
	return execute(p, cmd)
}
