from fontTools import ttLib

tt = ttLib.TTFont("SourceCodePro-Regular.ttf")

for table in tt["cmap"].tables:
    if table.platformID == 0:
        if table.cmap[0x25CA] == "lozenge":
            print("Patching table", table.platformID, table.platEncID)
            table.cmap[0x22C4] = "lozenge"

tt.save("SourceCodePro-Regular-patched.ttf")