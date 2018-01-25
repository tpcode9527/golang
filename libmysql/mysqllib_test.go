package mysqllib

import (
	"fmt"
	"testing"
)

func TestInst(t *testing.T) {
	inst := NewDb()
	err := inst.Open("root", "Suitang@20170601", "192.168.1.133", "3306", "tp_test", "utf8mb4", 8, 8, nil)
	fmt.Println("Open result:", err)

	result, err1 := inst.Query("select * from t_mbo_organization where did=5")
	fmt.Println(result, err1)

	result, err = inst.Query("select count(*) from t_mbo_organization")
	fmt.Println(result, err)

	err = inst.Exec("select count(*) from t_mbo_organization")

	transaction := inst.NewTransaction()
	transaction.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(8,'部门3')")
	transaction.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(7,'部门33')")
	transaction.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(9,'部门4')")
	err = transaction.Run()
}

func TestGlobal(t *testing.T) {
	err := InitMysqlDB("root", "Suitang@20170601", "192.168.1.133", "3306", "tp_test", "utf8mb4", 8, 8, nil)
	if nil != err {
		fmt.Println("InitMysqlDB fail. error", err)
	}

	rslt, err1 := MysqlQuery("select * from t_mbo_organization where did=5")
	fmt.Println(rslt, err1)

	err = MysqlExec("select count(*) from t_mbo_organization")

	trans := NewMysqlTransaction()
	trans.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(8,'部门3')")
	trans.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(7,'部门33')")
	trans.Exec("INSERT INTO t_mbo_organization(did,dname) VALUES(9,'部门4')")
	err = trans.Run()
}
