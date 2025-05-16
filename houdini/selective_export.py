print("...") # just a visual cue

type_geo = hou.objNodeTypeCategory().nodeType("geo")
assert type_geo is not None, "failed to load 'geo' type"
type_subnet = hou.objNodeTypeCategory().nodeType("subnet")
assert type_subnet is not None, "failed to load 'subnet' type"
type_dth_pose_asset = hou.sopNodeTypeCategory().nodeType("DazToHuePoseAsset")
assert type_dth_pose_asset is not None, "failed to load 'DazToHuePoseAsset' type"
type_dth_export = hou.sopNodeTypeCategory().nodeType("DazToHueExport")
assert type_dth_export is not None, "failed to load 'DazToHueExport' type"

linked_templates = [s.strip() for s in hou.evalParm("linked_templates").splitlines()]
assert len(linked_templates) >= 1, "At least one linked template name must be provided"

target_nodes = []

def recursive_search(entry_node):
    for child_node in entry_node.children():
        if any(child_node.type() == x for x in [type_geo, type_subnet]):
            recursive_search(child_node)
        elif child_node.type() == type_dth_pose_asset:
            if any(child_node.evalParm("node_linking_source").endswith(t) for t in linked_templates):
                for sibling_node in child_node.parent().children():
                    if sibling_node != child_node and sibling_node.type() == type_dth_export:
                        target_nodes.append(sibling_node)
            
recursive_search(hou.node("/obj"))

for target_node in target_nodes:
    print(f"Triggering export for: {target_node.path()}")
    target_node.parm("export_trigger").pressButton()