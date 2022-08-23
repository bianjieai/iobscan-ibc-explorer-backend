package task

import "testing"

func Test_IbcTxMigrate(t *testing.T) {
	new(IbcTxMigrateTask).Run()
}

func Test_MigrateSetting(t *testing.T) {
	if err := new(IbcTxMigrateTask).migrateSetting(); err != nil {
		t.Fatal(err)
	}
}

func Test_MigrateNormal(t *testing.T) {
	if err := new(IbcTxMigrateTask).migrateNormal(); err != nil {
		t.Fatal(err)
	}
}
