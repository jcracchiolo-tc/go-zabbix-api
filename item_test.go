package zabbix_test

import (
	"testing"

	zapi "github.com/jcracchiolo-tc/go-zabbix.api"
)

func testCreateItem(host *zapi.Host, t *testing.T) *zapi.Item {
	items := zapi.Items{{
		HostID: host.HostID,
		Key:    "key.lala.laa",
		Name:   "name for key",
		Type:   zapi.ZabbixTrapper,
		Delay:  "0",
	}}
	err := testGetAPI(t).ItemsCreate(items)
	if err != nil {
		t.Fatal(err)
	}
	return &items[0]
}

func testCreateItemWithApplication(app *zapi.Application, t *testing.T) *zapi.Item {
	items := zapi.Items{{
		HostID:         app.HostID,
		Key:            "key.lala.laa",
		Name:           "name for key",
		Type:           zapi.ZabbixTrapper,
		Delay:          "0",
		ApplicationIds: []string{app.ApplicationID},
	}}
	err := testGetAPI(t).ItemsCreate(items)
	if err != nil {
		t.Fatal(err)
	}
	return &items[0]
}

func testDeleteItem(item *zapi.Item, t *testing.T) {
	err := testGetAPI(t).ItemsDelete(zapi.Items{*item})
	if err != nil {
		t.Fatal(err)
	}
}

func TestItems(t *testing.T) {
	api := testGetAPI(t)

	group := testCreateHostGroup(t)
	defer testDeleteHostGroup(group, t)

	host := testCreateHost(group, t)
	defer testDeleteHost(host, t)

	item := testCreateItem(host, t)

	_, err := api.ItemGetByID(item.ItemID)
	if err != nil {
		t.Fatal(err)
	}

	item.Name = "another name"
	err = api.ItemsUpdate(zapi.Items{*item})
	if err != nil {
		t.Error(err)
	}

	testDeleteItem(item, t)
}

func TestItemsWithApplication(t *testing.T) {
	skipTestIfVersionGreaterThanOrEqual(t, "5.4", "dropped support for Application API")

	api := testGetAPI(t)

	group := testCreateHostGroup(t)
	defer testDeleteHostGroup(group, t)

	host := testCreateHost(group, t)
	defer testDeleteHost(host, t)

	app := testCreateApplication(host, t)
	defer testDeleteApplication(app, t)

	items, err := api.ItemsGetByApplicationID(app.ApplicationID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("Found items")
	}

	item := testCreateItemWithApplication(app, t)

	_, err = api.ItemGetByID(item.ItemID)
	if err != nil {
		t.Fatal(err)
	}

	item.Name = "another name"
	err = api.ItemsUpdate(zapi.Items{*item})
	if err != nil {
		t.Error(err)
	}

	testDeleteItem(item, t)
}
