print("...") # just a visual cue

text_map = hou.evalParm("linked_templates_preset_map");
text_map_rows = text_map.splitlines();

linked_templates_preset_map = {}
for text_map_row in text_map_rows:
    linked_template_name, linked_template_preset_path = text_map_row.split("=")
    
    linked_template_name = linked_template_name.strip()
    linked_template_preset_path = linked_template_preset_path.strip()
    
    assert linked_template_preset_path.endswith(".csv"), f"File must be CSV file: {linked_template_preset_path}"
    assert linked_templates_preset_map.get(linked_template_name) == None, f"Duplicate template name: {linked_template_name}"
    
    print(f"Searching for file: {linked_template_preset_path}")
    linked_template_preset_value = hou.findFile(linked_template_preset_path) # raises error if file does not exist
    
    linked_templates_preset_map[linked_template_name] = linked_template_preset_value

print()
print("Searching for pose asset nodes ...")
print()

type_geo = hou.objNodeTypeCategory().nodeType("geo")
assert type_geo is not None, "failed to load 'geo' type"
type_subnet = hou.objNodeTypeCategory().nodeType("subnet")
assert type_subnet is not None, "failed to load 'subnet' type"
type_dth_pose_asset = hou.sopNodeTypeCategory().nodeType("DazToHuePoseAsset")
assert type_dth_pose_asset is not None, "failed to load 'DazToHuePoseAsset' type"

linked_templates = linked_templates_preset_map.keys()
assert len(linked_templates) >= 1, "At least one linked template name must be provided"

target_nodes = []

def recursive_search(entry_node):
    for child_node in entry_node.children():
        if any(child_node.type() == x for x in [type_geo, type_subnet]):
            recursive_search(child_node)
        elif child_node.type() == type_dth_pose_asset and any(child_node.path().endswith(x) for x in linked_templates):
            target_nodes.append(child_node)
            
recursive_search(hou.node("/obj"))

for pose_asset_node in target_nodes:
    print(f"Refreshing: {pose_asset_node.path()}")
    
    for parm in pose_asset_node.parms():
        try:
            parm.revertToDefaults()
        except:
            pass # Some parameters can't be reverted (e.g., locked or hidden)
    kwargs = {
        "node": pose_asset_node,
        "filepath": linked_templates_preset_map[pose_asset_node.name()],
        "clear": False # does not work properly at the moment
    }
    pose_asset_node.hdaModule().import_from_csv(kwargs)