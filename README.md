# DTH-AutoTool
Automate export process from DAZ studio for DTH workflow.

:warning: For windows only :warning:

#### :bangbang: Disclaimer
This repository is still work in progress and mostly reflects my personal workflow setup. All recommendations and suggestions for improvements are welcomed. You are encouraged to play around with the code and adjust according to your personal requirements.

# Download template app
If you want to avoid headaches of manual compilation. Pre-compiled application with demo configuration can be downloaded [here](https://e.pcloud.link/publink/show?code=XZMtClZW3QA1kTh8z7LC0shRbvQ1kJcNc8V).

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
Therefore [OpenCV](https://opencv.org/) bindings needs [opencv](https://opencv.org/) as a prerequisite. OpenCV should be compiled as part of `robotgo` dependency automatically, but see _compilation issues_ section below.

You will need `gcc` compiler and `cmake` already pre-installed to compile OpenCV libraries. See [here](https://gocv.io/getting-started/windows/).

To link OpenCV libraries to DTH-AutoTool project, use `CGO_LDFLAGS` with following value:

```
-lade -lIlmImf -ljpeg-turbo -lopencv_aruco4100 -lopencv_bgsegm4100 -lopencv_bioinspired4100 -lopencv_calib3d4100 -lopencv_ccalib4100 -lopencv_core4100 -lopencv_datasets4100 -lopencv_dnn_objdetect4100 -lopencv_dnn_superres4100 -lopencv_dnn4100 -lopencv_dpm4100 -lopencv_face4100 -lopencv_features2d4100 -lopencv_flann4100 -lopencv_fuzzy4100 -lopencv_gapi4100 -lopencv_hfs4100 -lopencv_highgui4100 -lopencv_img_hash4100 -lopencv_imgcodecs4100 -lopencv_imgproc4100 -lopencv_intensity_transform4100 -lopencv_line_descriptor4100 -lopencv_mcc4100 -lopencv_ml4100 -lopencv_objdetect4100 -lopencv_optflow4100 -lopencv_phase_unwrapping4100 -lopencv_photo4100 -lopencv_plot4100 -lopencv_quality4100 -lopencv_rapid4100 -lopencv_reg4100 -lopencv_rgbd4100 -lopencv_shape4100 -lopencv_signal4100 -lopencv_stereo4100 -lopencv_stitching4100 -lopencv_structured_light4100 -lopencv_superres4100 -lopencv_surface_matching4100 -lopencv_text4100 -lopencv_tracking4100 -lopencv_video4100 -lopencv_videoio4100 -lopencv_videostab4100 -lopencv_wechat_qrcode4100 -lopencv_xfeatures2d4100 -lopencv_ximgproc4100 -lopencv_xobjdetect4100 -lopencv_xphoto4100 -lopenjp2 -lpng -lprotobuf -ltiff -lwebp -lzlib -lole32 -loleaut32 -luuid -lcomdlg32 -lgdi32
```

To inform the linker about where to find physical files, you can set `LIBRARY_PATH` environment variable with value:

```
C:\opencv\build\install\x64\mingw\staticlib
```

or you may want to use `-L<path>` cmd-line option (whatever you prefer).

You will need `go 1.23` compiler to compile the DTH-AutoTool project itself.
To start compilation run:

`$ go build`

command from the root directory of the project.

For more information on how to compile `robotgo` and `gocv` (OpenCV bindings for `go`) see links below:
- [Starting with GOCV on windows](https://gocv.io/getting-started/windows/)
- [GOCV repository](https://github.com/vcaesar/gcv)
- [RobotGo repository](https://github.com/go-vgo/robotgo)
- [OpenCV official site](https://opencv.org/)

#### Pre-built OpenCV binaries
You can use 64-bit binaries I already built from [this link](https://e.pcloud.link/publink/show?code=XZHtClZ2V7k51uqqikNlokbNBW0i7mMgwuV).
Note that these libraries were built with `-D_GLIBCXX_USE_CXX11_ABI=0` option so when building the DTH-AutoTool project you'll have to set `CGO_CFLAGS` with this option as well to link it properly.

#### Dealing with `robotgo` compilation issues
If you are facing issues related to `freetype` module (OpenCV) compilation, you may need to manually skip this module (it is not used in this project anyways).

If you already ran compilation of the project, you probably already have `GOCV` sources downloaded. Otherwise run `go get` command to fetch it.

Enter the `go`s staging directory, usually something like `<user-home-directory>\Go\pkg\mod\gocv.io\x\gocv@v0.38.0`. Build commands are scripted in `win_build_opencv.cmd`. Open the file in text editor and locate the `cmake` invocation. Then add this as a command line parameter `-DBUILD_opencv_freetype=OFF`.

After this change, you want only to re-generate the `makefile` using `cmake` and compile the libraries. However `win_build_opencv.cmd` contains also download steps that would overwrite your changes so you either want to comment them out and re-run the `win_build_opencv.cmd` or run `cmake` and `make` commands yourself (manually) as they are defined inside the `win_build_opencv.cmd`.

After running `make install` you should see compiled OpenCV libraries in `C:\opencv\build\install\x64\mingw\staticlib` by default.