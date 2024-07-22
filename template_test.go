package zabbix_test

import (
	"testing"

	zapi "github.com/jcracchiolo-tc/go-zabbix.api"
)

func testCreateTemplate(hostGroups *zapi.HostGroupIDs, t *testing.T) *zapi.Template {
	template := zapi.Templates{zapi.Template{
		Host:   "template name",
		Groups: *hostGroups,
	}}
	err := testGetAPI(t).TemplatesCreate(template)
	if err != nil {
		t.Fatal(err)
	}
	return &template[0]
}

func testDeleteTemplate(template *zapi.Template, t *testing.T) {
	err := testGetAPI(t).TemplatesDelete(zapi.Templates{*template})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplates(t *testing.T) {
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
	if template.TemplateID == "" {
		t.Errorf("Template id is empty %#v", template)
	}

	templates, err := api.TemplatesGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) == 0 {
		t.Fatal("No templates were obtained")
	}

	_, err = api.TemplateGetByID(template.TemplateID)
	if err != nil {
		t.Error(err)
	}

	template.Name = "new template name"
	err = api.TemplatesUpdate(zapi.Templates{*template})
	if err != nil {
		t.Error(err)
	}

	testDeleteTemplate(template, t)
}
