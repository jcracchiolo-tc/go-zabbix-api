package zabbix_test

import (
	"fmt"
	"testing"

	dd "github.com/jcracchiolo-tc/go-zabbix.api"
	zapi "github.com/jcracchiolo-tc/go-zabbix.api"
)

func testCreateTriggerPrototype(template *dd.Template, item *dd.ItemPrototype, t *testing.T) *dd.TriggerPrototype {
	triggers := dd.TriggerPrototypes{{
		Description: "test trigger prototype",
		Expression:  fmt.Sprintf("{%s:%s.last()}=312", template.Host, item.Key),
		Priority:    zabbix.Warning,
	}}
	err := testGetAPI(t).TriggerPrototypesCreate(triggers)
	if err != nil {
		t.Fatal(err)
	}
	return &triggers[0]
}

func testDeleteTriggerPrototype(trigger *dd.TriggerPrototype, t *testing.T) {
	err := testGetAPI(t).TriggerPrototypesDelete(dd.TriggerPrototypes{*trigger})
	if err != nil {
		t.Fatal(err)
	}
}

func testTriggerPrototype(t *testing.T) {
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
	defer testDeleteItemPrototype(itemPrototype, t)

	triggerPrototype := testCreateTriggerPrototype(template, itemPrototype, t)

	getTriggerPrototype, err := api.TriggerPrototypesGet(dd.Params{"triggerids": triggerPrototype.TriggerID})
	if err != nil {
		t.Error(err)
	}
	if len(getTriggerPrototype) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getTriggerPrototype))
	}

	triggerPrototype.Description = "update_trigger_prototype_name"
	err = api.TriggerPrototypesUpdate(dd.TriggerPrototypes{*triggerPrototype})
	if err != nil {
		t.Error(err)
	}

	getByIdTriggerPrototype, err := api.TriggerPrototypeGetByID(triggerPrototype.TriggerID)
	if err != nil {
		t.Error(err)
	}
	if getByIdTriggerPrototype.Description != triggerPrototype.Description {
		t.Errorf("Trigger prototype description is %q and should be %q", getByIdTriggerPrototype.Description, triggerPrototype.Description)
	}

	testDeleteTriggerPrototype(triggerPrototype, t)
}
