# DTH-AutoTool
Automate export process from DAZ studio for DTH workflow.

:warning: For windows only :warning:

#### :bangbang: Disclaimer
This repository is still work in progress and mostly reflects my personal workflow setup. All recommendations and suggestions for improvements are welcomed. You are encouraged to play around with the code and adjust according to your personal requirements.

# Download the latest build
If you want to avoid headaches of manual compilation. Pre-compiled application with scripts can be downloaded [here](https://e.pcloud.link/publink/show?code=kZNiwlZYdusX3O8qqFgLVS94a7KU8nxj4Sk).

# Build demo example app
- Download installation script and example config folder from the _example_ repository folder.
- Open windows powershell window and enter the location on your local computer.
- Run the installation script with `.\install.ps1`, if not running powershell in administrator mode, you'll be prompted to elevate (this is required to create symbolic links by the script).
- Open _DAZ Studio Install Manager_ and install **Genesis 8 [Male|Female] Starter Essentials** (used by the example).
- Install [DTH](https://www.artstation.com/marketplace/p/BLM5K/daztohue), [Sagan Alembic Exporter Plugin](https://www.daz3d.com/forums/discussion/428856/sagan-a-daz-studio-to-blender-alembic-exporter/p1) and [DazToMaya Bridge](https://www.daz3d.com/daz-to-maya-bridge).
- Since the DTH-AutoTool app uses `Asset Manager` API, all assets including DTH must be indexed in DAZ database. The easiest way to include it is to create a new category from the _DTH_ folder and its subfolders.
  - Open the _Content Library_ pane and navigate to the location where you installed _DTH Contents_ (usually `DAZ Studio Formats/My Library/DazToHue`).
  - Right click the `DazToHue` and select `Create a Category from->Selected Folder & Sub-Folders`.
  - You should see now `DazToHue` in `Categories` section.
- In _DAZ Studio_, open the _Script IDE_ pane and load the the _main_ script located in app folder and execute it.

If everything is setup properly, script and DTH-AutoTool should automatically load and export models into the demo app's _Export_ folder.

Installation script will automatically create necessary folders and downloads the latest version of the app.
The default/demo configuration will be used with _Genesis 8.1_ characters and some basic clothing which is part of the essentials pack.

## Demo example notes
Example works with following subfolders:
- `config`: the one downloaded from repository, containts sample config files
- `app`: the location where the latest app will be downloaded
- `temp`: symbolic links to intermediate alembic cache/fbx export folders will be created
- `exports`: exported DTH objects will be stored here if default configuration was not changed

The default/demo configuration is assuming intermediate FBX (Daz to Maya) folder to be: `<user-home-directory>\Documents\DAZ 3D\Bridges\Daz To Maya\Exports`.

The alembic cache export directory is altered by the application itself, but the demo configuration sets it to be `<windows-apps-temp-dir>\DTH-AutoTool\Exports\ABC-Temp`. You can change it manually in `install.ps1` if needed.

Only symbolic links will be created inside `temp` directory.

# Install `tesseract` (optional)
DTH-AutoTool can handle popups using `key-sequencing` or `visual-lookup` methods.
The latter may optionally benefit from [Tesseract OCR engine](https://github.com/tesseract-ocr/tesseract) installed on your system as command line tool to validate located GUI elements like buttons or labels.
It is highly recommended to install `tesseract` app when using `visual-lookup` mode. You can disable this feature in `app.config.json` file by setting `screen_search_validation->tesseract_validation->enabled` to `false`.

# How it works
- Prepare your assets upfront and save them as a `*.duf` files. DTH-AutoTool app uses `Asset Manager` API to locate the assets. Your *.duf files are automatically added to DAZ database. Manually copied `*.duf` files to common DAZ directories can be browsed from UI as files but are not part of database automatically (scanning files is not supported in DTH-AutoTool at the moment). The easiest way to index them is to add assets to existing categories or for example create entirely new category (e.g. DTH category from the entire _DazToHue_ folder).
- In `project.config.json` under `templates` section, list all templates with their `character`, `clothing` and/or `morphs` subsections. None of them is mandatory for template, which means you can split different layers across templates and combine them later in `exportables` section. Every template name must be unique.
- `character`, `clothing` and/or `morphs` subsections contains _layer groups_. Every _layer group_ must have a globally unique name and contains array of layers (`*duf` files) that logically belongs together (for example, some parts of clothings may exist as separate DAZ assets, but you may want to handle it as a single group/unit - single outfit).
- `templates` can only be referenced in `exportables`. Exportables are producing real outpus. You can combine several templates together in one exportable and _override_ specific layer groups, _prepend_ or _append_ to them, or completely _replace_ entire `character`, `clothing` or `morphs` sections if needed. This provides higher granularity in setting up specific characters/figures to export.

### Character vs. Clothing vs. Morphs
- `character` layers are applied first and should define the characters basic look. If you need, you can apply also clothing `*.duf` files but the idea of these layers is their single appliance at the beginning of the pipeline.
- `clothing` layers are applied after the `character` layers but **always only one clothing** per pipeline run. After one of the clothings is applied, `morphs` layers are loaded.
- `morphs` layers are loaded after `clothing` and usually should contain timeline JCM frames or morphs to be handled by DTH workflow.
- Once layers are loaded, export is triggered. Then the process is repeated for each available `clothing`.

### Edit mode
- If you need to load only the layers without export being triggered, use `edit_mode` configuration section.
- Set `is_enabled` to `true` and specify `exportable_name` to valid value.
- If `exportable` has clothing layers, you can also specify `clothing_name` to load.

# How to compile
This project depends on [robotgo](https://github.com/go-vgo/robotgo) and [gocv](https://github.com/vcaesar/gcv) libraries.
Therefore [OpenCV](https://opencv.org/) bindings needs [opencv](https://opencv.org/) as a prerequisite.
If you are having troubles when following steps below due to version changes in _OpenCV_/_gocv_ versions or compiler versions, I'll keep tested binaries on my public [pcloud folder](https://e.pcloud.link/publink/show?code=kZNiwlZYdusX3O8qqFgLVS94a7KU8nxj4Sk).

Steps below are for windows only (the tool is not supported for other platforms at the moment).

## Prerequisites
* MinGW GCC compiler (to compile OpenCV libraries)
* CMake (to generate build files, makefile)
* Go compiler (to build the actual tool)

Check the official [gocv](https://github.com/vcaesar/gcv) README for compiler version to use.
The latest tested MinGW version was 8.1.0.

At the time of writing, _CMake_ version 4.X was already available. Do not use CMake 4+!
OpenCV libraries have not been ported yet to CMake 4, **use 3.5+ instead**.

## Steps
* Install CMake 3.5+ (not 4.X+).
* Install MinGW compiler 8.1.0 and make sure your `PATH` variable contains reference to its `bin` directory.
* Install Golang toolchain.
* Clone this repository and open command line at the root folder.
* Run `go build` command.
* If you did not previously compiled & installed _OpenCV_, previous build command will fail. However, it will still download the _gocv_ dependency onto your computer.
* Enter the location where go downloads/caches packages on your computer, usually something like `<homedir>\go\pkg\mod`.
* Find `gocv.io\x\gocv@vX.Y.Z` folder (`X.Y.Z` will be the version used by this project).
* Open a command line and run `.\win_build_opencv.cmd`.
  * The `CMD` file will automatically download required _OpenCV_ sources and triggers CMake & mingw32-make commands for you (this will take a while to compile and install).
  * By default, above `CMD` file will compile & install the _OpenCV_ into `C:\OpenCV\build\install` folder. You can tweak the `win_build_opencv.cmd` to your liking before the compilation but make sure to update paths in subsequent steps accordingly. You also need only `bin`, `include` and `x64\mingw\lib` (or `x64\mingw\bin` - see below) folders to build the _dth-autotool_, everything else can be removed.
* To compile the _dth-autotool_, you'll need to set up several environment variables before the build:
  * You should set all environment variables below just in shell prior building the project - not globally, unless you really intend to (probably for _OpenCV_ paths).
  * If building project with shared libraries reference, append following to your `PATH` environment `C:\opencv\build\install\x64\mingw\bin`.
  * If building project with static libraries reference:
    * Create new environment variable `LIBRARY_PATH` and append this location it `C:\opencv\build\install\x64\mingw\lib`.
    * (-or-) add `-LC:\opencv\build\install\x64\mingw\bin` to the `CGO_LDFLAGS`.
    * **IMPORTANT** - in my case, static libraries were called like `libopencv_xxx4110.dll.a` after the compilation & installation which is incorrect (probably a bug in CMake definitions). To rename all files, you can run following powershell command from the directory with static libraries installed: `Get-ChildItem *4110.dll.a | Rename-Item -NewName { $_.Name -Replace '4110.dll','' }`. Note that `4110` here is the _OpenCV_ version. Its removal is optional, but I like it removed, so my `-l` flags below are cleaner. If you keep it, update names after `-l` flags accordingly.
  * Set `CGO_LDFLAGS` to `-lopencv_stitching -lopencv_superres -lopencv_videostab -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dpm -lopencv_face -lopencv_photo -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_reg -lopencv_rgbd -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_surface_matching -lopencv_tracking -lopencv_datasets -lopencv_dnn -lopencv_plot -lopencv_xfeatures2d -lopencv_shape -lopencv_video -lopencv_ml -lopencv_ximgproc -lopencv_calib3d -lopencv_features2d -lopencv_highgui -lopencv_videoio -lopencv_flann -lopencv_xobjdetect -lopencv_imgcodecs -lopencv_objdetect -lopencv_xphoto -lopencv_imgproc -lopencv_core`.
  * Set `CGO_CXXFLAGS` to `--std=c++11`.
  * Set `CGO_CPPFLAGS` to `-IC:\opencv\build\install\include`.
* Run `go build -o build/dth-autotool.exe` command.

### What if compilation fails
If compilation fails on _OpenCV_ or _gocv_ dependencies, please refer to the official repositories for guidance.
Otherwise feel free to raise an issue in this repo or ask a question.

You can also download tested version of compiler, or pre-compiled 64-bit _OpenCV_ libraries I used from [my pcloud folder](https://e.pcloud.link/publink/show?code=kZNiwlZYdusX3O8qqFgLVS94a7KU8nxj4Sk).
Using my _OpenCV_ build you can skip straightly to the _dth-autotool_ build part, saving you some time or potential headaches.

#### Additional links
- [Starting with GOCV on windows](https://gocv.io/getting-started/windows/)
- [GOCV repository](https://github.com/vcaesar/gcv)
- [RobotGo repository](https://github.com/go-vgo/robotgo)
- [OpenCV official site](https://opencv.org/)