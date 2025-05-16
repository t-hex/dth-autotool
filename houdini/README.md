### About Houdini helper scripts

Houdini scripts are not directly needed **dth-autotool**. These are mostly helper scripts to support Houdini part of DTH workflow with some additional automation.

#### selective_export.py
If you organize your pose assets in such a way that you utilize DTH linking feature (parent/child nodes), you can use master node names to find all child nodes referencing them and trigger the *Export* automatically.

This can be useful for example in tandem with **dth-autotool**. If you for example keep ROM presets as **dth-autotool** templates with 1:1 relationship to the DTH master nodes, once you change the ROM template and re-export assets from DAZ, you can set the list of changed/impacted master nodes in the script to automatically re-export assets from Houdini as well.

#### linked_templates_refresh.py
Given the DTH linking feature and utilizing parent/child nodes relationship with pose asset nodes there are some issues to take into account.

- Linking seems to propagate changes well only after the link is established. When new link is created to the existing master node, not all properties are properly set on the child node (especially the dynamically generated or array like parameters like JCM groups etc.).
- Reloading assets from CSV on master nodes also has some issues as well as from python script `kwargs` directly when it comes to *clear* flag. This means you can only append parameters from CSV but not replace.

While these issues might be fixed in future, as a workaround this script allows you to create a key/value map of master node name vs. CSV backup and reload pose asset parameters properly for all listed master nodes. This will also ensure correct propagation of all changes to the linked nodes as parameters are cleaned up first and re-populated from CSV again, hence forcing linked nodes to catch up.