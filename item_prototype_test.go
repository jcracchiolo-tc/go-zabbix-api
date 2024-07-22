package zabbix_test

import (
	"testing"

	dd "github.com/jcracchiolo-tc/go-zabbix.api"
	zapi "github.com/jcracchiolo-tc/go-zabbix.api"
)

func testCreateItemPrototype(template *dd.Template, lldRule *dd.LLDRule, t *testing.T) *dd.ItemPrototype {
	items := dd.ItemPrototypes{{
		RuleID:      lldRule.ItemID,
		Delay:       "30",
		HostID:      template.TemplateID,
		InterfaceID: "0",
		Key:         "key.lala.lolo",
		Type:        dd.ZabbixAgent,
		Name:        "item prototype test",
		ValueType:   dd.Unsigned,
	}}
	err := testGetAPI(t).ItemPrototypesCreate(items)
	if err != nil {
		t.Fatal(err)
	}
	return &items[0]
}

func testDeleteItemPrototype(item *dd.ItemPrototype, t *testing.T) {
	err := testGetAPI(t).ItemPrototypesDelete(dd.ItemPrototypes{*item})
	if err != nil {
		t.Fatal(err)
	}
}

func testItemPrototype(t *testing.T) {
	api := testGetAPI(t)

	// Zabbix v6.2 introduced Template Groups and requires them for Templates
	var groupIds zapi.HostGroupIDs
	if compLessThan, _ := isVersionLessThan(t, "6.2"); compLessThan {
		hostGroup := testCreateHostGroup(t)
		defer testDeleteHostGroup(hostGroup, t)
		groupIds = zapi.HostGroupIDs{
			{
				GroupID: hostGroup.GroupID,
			},
		}
	} else {
		templateGroup := testCreateTemplateGroup(t)
		defer testDeleteTemplateGroup(templateGroup, t)
		groupIds = zapi.HostGroupIDs{
			{
				GroupID: templateGroup.GroupID,
			},
		}
	}

	template := testCreateTemplate(&groupIds, t)
	defer testDeleteTemplate(template, t)

	lldRule := testCreateLLDRule(template, t)
	defer testDeleteLLDRule(lldRule, t)

	itemPrototype := testCreateItemPrototype(template, lldRule, t)

	getItemPrototype, err := api.ItemPrototypesGet(dd.Params{"itemids": itemPrototype.ItemID})
	if err != nil {
		t.Error(err)
	}
	if len(getItemPrototype) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getItemPrototype))
	}

	itemPrototype.Name = "update_item_prototype_name"
	err = api.ItemPrototypesUpdate(dd.ItemPrototypes{*itemPrototype})
	if err != nil {
		t.Error(err)
	}

	getByIdItemPrototype, err := api.ItemPrototypeGetByID(itemPrototype.ItemID)
	if err != nil {
		t.Error(err)
	}
	if getByIdItemPrototype.Name != itemPrototype.Name {
		t.Errorf("Item prototype name is %q and should be %q", getByIdItemPrototype.Name, itemPrototype.Name)
	}

	testDeleteItemPrototype(itemPrototype, t)
}
