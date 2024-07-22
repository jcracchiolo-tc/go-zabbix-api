package zabbix_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	zapi "github.com/jcracchiolo-tc/go-zabbix.api"
)

func testCreateTemplateGroup(t *testing.T) *zapi.TemplateGroup {
	TemplateGroups := zapi.TemplateGroups{{Name: fmt.Sprintf("zabbix-testing-%d", rand.Int())}}
	err := testGetAPI(t).TemplateGroupsCreate(TemplateGroups)
	if err != nil {
		t.Fatal(err)
	}
	return &TemplateGroups[0]
}

func testDeleteTemplateGroup(TemplateGroup *zapi.TemplateGroup, t *testing.T) {
	err := testGetAPI(t).TemplateGroupsDelete(zapi.TemplateGroups{*TemplateGroup})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplateGroups(t *testing.T) {
	skipTestIfVersionLessThan(t, "6.2", "introduced support for Template Groups API")

	api := testGetAPI(t)

	groups, err := api.TemplateGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}

	TemplateGroup := testCreateTemplateGroup(t)
	if TemplateGroup.GroupID == "" || TemplateGroup.Name == "" {
		t.Errorf("Something is empty: %#v", TemplateGroup)
	}

	TemplateGroup2, err := api.TemplateGroupGetByID(TemplateGroup.GroupID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(TemplateGroup, TemplateGroup2) {
		t.Errorf("Error getting group.\nOld group: %#v\nNew group: %#v", TemplateGroup, TemplateGroup2)
	}

	groups2, err := api.TemplateGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups2) != len(groups)+1 {
		t.Errorf("Error creating group.\nOld groups: %#v\nNew groups: %#v", groups, groups2)
	}

	testDeleteTemplateGroup(TemplateGroup, t)

	groups2, err = api.TemplateGroupsGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(groups, groups2) {
		t.Errorf("Error deleting group.\nOld groups: %#v\nNew groups: %#v", groups, groups2)
	}
}
