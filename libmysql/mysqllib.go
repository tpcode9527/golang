package mysqllib

import (
	"database/sql"
	"fmt"
	//"log"
	//"strconv"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

//日志打印回调函数
type LogPrint func(msg ...interface{}) (int, error)

//数据库连接实例
type MysqlInst struct {
	Inst     *sql.DB
	User     string
	Password string
	Server   string
	Port     string
	Name     string
	Charset  string
	MaxOpen  int
	MaxConn  int
	FnPrint  LogPrint
}

//事务实例
type Transaction struct {
	Inst  *MysqlInst
	Trans []string
}

//创建实例
func NewDb() *MysqlInst {
	return &MysqlInst{FnPrint: fmt.Println}
}

//创建事务对象
func (this *MysqlInst) NewTransaction() *Transaction {
	return &Transaction{Inst: this, Trans: make([]string, 0)}
}

//打开并测试连接数据库
func (this *MysqlInst) Open(user string, password string, server string, port string,
	name string, charset string, maxOpen int, maxConn int, fnPrint LogPrint) error {
	var err error = nil

	this.User = user
	this.Password = password
	this.Server = server
	this.Port = port
	this.Name = name
	this.Charset = charset
	this.MaxOpen = maxOpen
	this.MaxConn = maxConn
	if nil != fnPrint {
		this.FnPrint = fnPrint
	}

	dataSourceName := this.User + ":" + this.Password + "@tcp(" + this.Server + ":" + this.Port + ")/" + this.Name
	if len(this.Charset) > 0 {
		dataSourceName += "?charset=" + this.Charset
	}

	this.Inst, err = sql.Open("mysql", dataSourceName)
	if nil != err {
		this.FnPrint("Open mysql fail. error:", err)
		return err
	}

	if this.MaxOpen > 0 {
		this.Inst.SetMaxOpenConns(this.MaxOpen)
	}

	if this.MaxConn > 0 {
		this.Inst.SetMaxIdleConns(this.MaxConn)
	}

	//this.FnPrint("Open success.")

	err = this.Inst.Ping()
	if nil != err {
		this.FnPrint("Ping mysql fail. error:", err)
		return err
	}
	//this.FnPrint("Ping success.")

	return err
}

//查询语句
func (this *MysqlInst) Query(sql string) ([]map[string]string, error) {
	var resut []map[string]string

	rows, err := this.Inst.Query(sql)
	if nil != err {
		this.FnPrint("Query fail. sql:(", sql, "). error:", err)
		return nil, err
	}
	defer rows.Close()

	resut = make([]map[string]string, 0)

	columns, _ := rows.Columns()
	//this.FnPrint("Columns:", columns)

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	//i := 0

	for rows.Next() {

		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)

		//resut[i] = make(map[string]string)
		record := make(map[string]string)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}

		resut = append(resut, record)
	}

	//this.FnPrint("recod:", resut)

	return resut, err
}

//执行
func (this *MysqlInst) Exec(sql string) error {
	_, err := this.Inst.Exec(sql)
	if nil != err {
		this.FnPrint("Exec fail. sql(", sql, "). error:", err)
	}
	return err
}

//保存事务中的子语句
func (this *Transaction) Exec(sql string) {
	this.Trans = append(this.Trans, sql)
}

//执行事务
func (this *Transaction) Run() error {
	tx, err := this.Inst.Inst.Begin()

	for _, sql := range this.Trans {

		_, err = tx.Exec(sql)
		if nil != err {
			this.Inst.FnPrint("Transaction fail. sql(", sql, "). error:", err)
			break
		}
	}

	if nil == err {
		err = tx.Commit()
	}

	if nil != err {
		tx.Rollback()
	}

	return err
}

/*********** 单例模式访问数据库 *************/
var g_dbInst *MysqlInst = nil

func init() {
	g_dbInst = NewDb()
}

func InitMysqlDB(user string, password string, server string, port string,
	name string, charset string, maxOpen int, maxConn int, fnPrint LogPrint) error {
	return g_dbInst.Open(user, password, server, port, name, charset, maxOpen, maxConn, fnPrint)
}

func NewMysqlTransaction() *Transaction {
	return g_dbInst.NewTransaction()
}

func MysqlQuery(sql string) ([]map[string]string, error) {
	return g_dbInst.Query(sql)
}

func MysqlExec(sql string) error {
	return g_dbInst.Exec(sql)
}
