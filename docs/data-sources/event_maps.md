---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stacuity_event_maps Data Source - stacuity"
subcategory: ""
description: |-
  Fetches the list of eventMaps.
---

# stacuity_event_maps (Data Source)

Fetches the list of eventMaps.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Attributes) (see [below for nested schema](#nestedatt--filter))

### Read-Only

- `eventmaps` (Attributes List) List of event maps. (see [below for nested schema](#nestedatt--eventmaps))

<a id="nestedatt--filter"></a>
### Nested Schema for `filter`

Optional:

- `filter` (String) Filter the results. Example 'name:TerraForm Test,monkier:tf-test'
- `limit` (Number) How many results to return
- `offset` (Number) What offset to use when querying
- `sort_by` (String) Sort by any property. Example 'asc(property),desc(property)'


<a id="nestedatt--eventmaps"></a>
### Nested Schema for `eventmaps`

Read-Only:

- `event_scope` (Attributes) (see [below for nested schema](#nestedatt--eventmaps--event_scope))
- `id` (String) Unique identifier for the event map.
- `moniker` (String) API Moniker of the event map
- `name` (String) Name of the event map
- `subscriptions` (Attributes Set) List of subscriptions attached to event map. (see [below for nested schema](#nestedatt--eventmaps--subscriptions))

<a id="nestedatt--eventmaps--event_scope"></a>
### Nested Schema for `eventmaps.event_scope`

Read-Only:

- `active` (Boolean) Active status of the event scope
- `moniker` (String) API Moniker for the type of event scope
- `name` (String) Name of the event scope


<a id="nestedatt--eventmaps--subscriptions"></a>
### Nested Schema for `eventmaps.subscriptions`

Read-Only:

- `event_endpoint` (Attributes) (see [below for nested schema](#nestedatt--eventmaps--subscriptions--event_endpoint))
- `event_map` (Attributes) (see [below for nested schema](#nestedatt--eventmaps--subscriptions--event_map))
- `event_type` (Attributes) (see [below for nested schema](#nestedatt--eventmaps--subscriptions--event_type))

<a id="nestedatt--eventmaps--subscriptions--event_endpoint"></a>
### Nested Schema for `eventmaps.subscriptions.event_endpoint`

Read-Only:

- `moniker` (String) API Moniker for the type of event scope
- `name` (String) Name of the event scope
- `summary_description` (String) Basic description
- `type` (String) Name of the event handler


<a id="nestedatt--eventmaps--subscriptions--event_map"></a>
### Nested Schema for `eventmaps.subscriptions.event_map`

Read-Only:

- `moniker` (String) API Moniker for the type of event scope
- `name` (String) Name of the event scope


<a id="nestedatt--eventmaps--subscriptions--event_type"></a>
### Nested Schema for `eventmaps.subscriptions.event_type`

Read-Only:

- `active` (Boolean) Active status of the event scope
- `moniker` (String) API Moniker for the type of event scope
- `name` (String) Name of the event scope
