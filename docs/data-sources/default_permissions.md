---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tableau_default_permissions Data Source - terraform-provider-tableau"
subcategory: ""
description: |-
  Retrieve project details
---

# tableau_default_permissions (Data Source)

Retrieve project details

## Example Usage

```terraform
data "tableau_projects" "all" {
}

data "tableau_default_permissions" "default_permissions" {
    project_id  = data.tableau_projects.all.projects[0].id
    target_type = "workbooks"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) ID of the project
- `target_type` (String) Permissions for: databases,dataroles,datasources,flows,lenses,metrics,tables,virtualconnections,workbooks

### Read-Only

- `grantee_capabilities` (Attributes List) List of grantee capabilities for users and groups (see [below for nested schema](#nestedatt--grantee_capabilities))

<a id="nestedatt--grantee_capabilities"></a>
### Nested Schema for `grantee_capabilities`

Read-Only:

- `capabilities` (Attributes List) List of grantee capabilities for users and groups (see [below for nested schema](#nestedatt--grantee_capabilities--capabilities))
- `group_id` (String) ID of the group
- `user_id` (String) ID of the user

<a id="nestedatt--grantee_capabilities--capabilities"></a>
### Nested Schema for `grantee_capabilities.capabilities`

Read-Only:

- `mode` (String) Mode of the capability (Allow/Deny)
- `name` (String) Name of the capability
