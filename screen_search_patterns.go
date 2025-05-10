package main

import "errors"

type RefPatternType uint64

const (
	RefPattern_WinHandler_DlgFilenameLbl RefPatternType = iota
	RefPattern_WinHandler_DlgOpenCancelBtnGrp
	RefPattern_WinHandler_DlgOpenBtn

	RefPattern_CommonHandler_AcceptBtn
	RefPattern_CommonHandler_RadioBtn
	RefPattern_CommonHandler_OkBtn
	RefPattern_CommonHandler_YesBtn

	RefPattern_DtmHandler_AcceptCancelBtnGrp
	RefPattern_DtmHandler_WhatIsThisBtn
	RefPattern_DtmHandler_SubdLevelsBtn
	RefPattern_DtmHandler_SubdLevelsGrpTop
	RefPattern_DtmHandler_MainExportOptionsGrp
	RefPattern_DtmHandler_AssetTypeLbl
	RefPattern_DtmHandler_AssetNameLbl

	RefPattern_SaeHandler_LoadSaveBtnGrp
	RefPattern_SaeHandler_LoadCfgBtn
	RefPattern_SaeHandler_EndFrameLbl
	RefPattern_SaeHandler_ExportBtn
	RefPattern_SaeHandler_DoneBtn

	RefPattern_TlrsHandler_YesNoBtnGrp

	RefPattern_LchHandler_ApplyToSelectedRadioBtn
)

type ScreenSearchImgReferencePattern struct {
	PatternType    RefPatternType
	ImageFilePaths []string
	ValidationText string
	Psm            TesseractPSM
	Language       TesseractLang
}

func NewSSIRP(pattern RefPatternType, imgPath string) *ScreenSearchImgReferencePattern {
	newPattern := ScreenSearchImgReferencePattern{}

	switch pattern {
	//region WinHandler
	case RefPattern_WinHandler_DlgFilenameLbl:
		newPattern.PatternType = RefPattern_WinHandler_DlgFilenameLbl
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/dlg_filename_lbl.png",
		}
	case RefPattern_WinHandler_DlgOpenCancelBtnGrp:
		newPattern.PatternType = RefPattern_WinHandler_DlgOpenCancelBtnGrp
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/win10_dlg_open_cancel_btn_grp.png",
			imgPath + "/common/win11_dlg_open_cancel_btn_grp.png",
			imgPath + "/common/win10_dlg_open_cancel_btn_grp_open_highlighted.png",
			imgPath + "/common/win11_dlg_open_cancel_btn_grp_open_highlighted.png",
			imgPath + "/common/win10_dlg_open_cancel_btn_grp_cancel_highlighted.png",
			imgPath + "/common/win11_dlg_open_cancel_btn_grp_cancel_highlighted.png",
		}
	case RefPattern_WinHandler_DlgOpenBtn:
		newPattern.PatternType = RefPattern_WinHandler_DlgOpenBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/win10_dlg_open_btn.png",
			imgPath + "/common/win11_dlg_open_btn.png",
			imgPath + "/common/win10_dlg_open_btn_highlighted.png",
			imgPath + "/common/win11_dlg_open_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Open"
		}
	//endregion DtmHandler
	//region CommonHandler
	case RefPattern_CommonHandler_AcceptBtn:
		newPattern.PatternType = RefPattern_CommonHandler_AcceptBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/accept_btn.png",
			imgPath + "/common/accept_btn_highlighted.png",
			imgPath + "/common/accept_btn_hover_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Accept"
		}
	case RefPattern_CommonHandler_RadioBtn:
		newPattern.PatternType = RefPattern_CommonHandler_RadioBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/radio_btn_off.png",
			imgPath + "/common/radio_btn_on.png",
		}
	case RefPattern_CommonHandler_OkBtn:
		newPattern.PatternType = RefPattern_CommonHandler_OkBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/ok_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "OK"
		}
	case RefPattern_CommonHandler_YesBtn:
		newPattern.PatternType = RefPattern_CommonHandler_YesBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/yes_btn.png",
			imgPath + "/common/yes_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Yes"
		}
	//endregion CommonHandler
	//region DtmHandler
	case RefPattern_DtmHandler_AcceptCancelBtnGrp:
		newPattern.PatternType = RefPattern_DtmHandler_AcceptCancelBtnGrp
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp.png",
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp_accept_highlighted.png",
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp_accept_hover_highlighted.png",
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp_cancel_highlighted.png",
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp_cancel_hover_highlighted.png",
			imgPath + "/daz_to_maya_dialog/accept_cancel_btn_grp_hover_highlighted.png",
		}
	case RefPattern_DtmHandler_WhatIsThisBtn:
		newPattern.PatternType = RefPattern_DtmHandler_WhatIsThisBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/what_is_this_btn.png",
		}
	case RefPattern_DtmHandler_SubdLevelsBtn:
		newPattern.PatternType = RefPattern_DtmHandler_SubdLevelsBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/bake_subd_levels_btn.png",
			imgPath + "/daz_to_maya_dialog/bake_subd_levels_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Bake Subdivision Levels"
		}
	case RefPattern_DtmHandler_SubdLevelsGrpTop:
		newPattern.PatternType = RefPattern_DtmHandler_SubdLevelsGrpTop
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/subd_levels_grp_top.png",
		}
	case RefPattern_DtmHandler_MainExportOptionsGrp:
		newPattern.PatternType = RefPattern_DtmHandler_MainExportOptionsGrp
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/main_export_options_grp.png",
		}
	case RefPattern_DtmHandler_AssetTypeLbl:
		newPattern.PatternType = RefPattern_DtmHandler_AssetTypeLbl
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/asset_type_lbl.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Asset Type"
		}
	case RefPattern_DtmHandler_AssetNameLbl:
		newPattern.PatternType = RefPattern_DtmHandler_AssetNameLbl
		newPattern.ImageFilePaths = []string{
			imgPath + "/daz_to_maya_dialog/asset_name_lbl.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Asset Name"
		}
	//endregion DtmHandler
	//region SaeHandler
	case RefPattern_SaeHandler_LoadSaveBtnGrp:
		newPattern.PatternType = RefPattern_SaeHandler_LoadSaveBtnGrp
		newPattern.ImageFilePaths = []string{
			imgPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp.png",
			imgPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp_load_highlighted.png",
			imgPath + "/sagan_alembic_exporter_dialog/load_save_cfg_btn_grp_save_highlighted.png",
		}
	case RefPattern_SaeHandler_LoadCfgBtn:
		newPattern.PatternType = RefPattern_SaeHandler_LoadCfgBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/sagan_alembic_exporter_dialog/load_cfg_btn.png",
			imgPath + "/sagan_alembic_exporter_dialog/load_cfg_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Load Config"
		}
	case RefPattern_SaeHandler_EndFrameLbl:
		newPattern.PatternType = RefPattern_SaeHandler_EndFrameLbl
		newPattern.ImageFilePaths = []string{
			imgPath + "/sagan_alembic_exporter_dialog/end_frame_lbl.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "End Frame"
		}
	case RefPattern_SaeHandler_ExportBtn:
		newPattern.PatternType = RefPattern_SaeHandler_ExportBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/export_btn.png",
			imgPath + "/common/export_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Export"
		}
	case RefPattern_SaeHandler_DoneBtn:
		newPattern.PatternType = RefPattern_SaeHandler_DoneBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/common/done_btn.png",
			imgPath + "/common/done_btn_highlighted.png",
		}
		if newPattern.ValidationText == "" {
			newPattern.ValidationText = "Done"
		}
	//endregion SaeHandler
	//region TlrsHandler
	case RefPattern_TlrsHandler_YesNoBtnGrp:
		newPattern.PatternType = RefPattern_TlrsHandler_YesNoBtnGrp
		newPattern.ImageFilePaths = []string{
			imgPath + "/timeline_resize_dialog/yes_no_btn_grp.png",
			imgPath + "/timeline_resize_dialog/yes_no_btn_grp_yes_highlighted.png",
			imgPath + "/timeline_resize_dialog/yes_no_btn_grp_no_highlighted.png",
		}
	//endregion TlrsHandler
	//region LchHandlerHandler
	case RefPattern_LchHandler_ApplyToSelectedRadioBtn:
		newPattern.PatternType = RefPattern_LchHandler_ApplyToSelectedRadioBtn
		newPattern.ImageFilePaths = []string{
			imgPath + "/load_character_dialog/apply_to_selected_radio_btn_off.png",
			imgPath + "/load_character_dialog/apply_to_selected_radio_btn_on.png",
		}
	//endregion LchHandlerHandler
	default:
		panic(errors.New("invalid alignment"))
	}

	return &newPattern
}
