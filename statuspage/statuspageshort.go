package statuspage

type IdccStatus struct {
	FbParametersKnown     bool `json:"fb_parameters_known"`
	AtcPublished          bool `json:"atc_published"`
	AtcValidation         bool `json:"atc_validation_detected"`
	AtcValidationFallback bool `json:"atc_validation_fallback_detected"`
	EmptyDomain           bool `json:"empty_domain_detected"`
	Fallback              bool `json:"fallback_detected"`
	Iva                   bool `json:"iva_detected"`
	IvaPresolved          bool `json:"iva_presolved_detected"`
	IvaFallback           bool `json:"iva_fallback_detected"`
	ReturnedBranches      bool `json:"returned_branches_detected"`
}

type StatusJson struct {
	Flowbased struct {
		DACC struct {
			Dfp               bool `json:"dfp_detected"`
			Spanning          bool `json:"spanning_detected"`
			FbParametersKnown bool `json:"fb_parameters_known"`
			IvaFallback       bool `json:"iva_fallback_detected"`
			IvaPresolved      bool `json:"iva_presolved_detected"`
			ReturnedBranches  bool `json:"returned_branches_detected"`
		} `json:"DACC"`
		IDCCb IdccStatus `json:"IDCC(b)"`
		IDCCc IdccStatus `json:"IDCC(c)"`
		IDCCd IdccStatus `json:"IDCC(d)"`
	} `json:"flowbased"`
}

func GetTableDataShort(data *StatusJson) [][]string {
	fbstatus := data.Flowbased
	light := func(green bool, yellow bool, red bool) string {
		if red {
			return "🔴"
		}
		if yellow {
			return "🟡"
		}
		if green {
			return "🟢"
		}

		return "⚪"
	}
	table := [][]string{
		{"", "DA", "(b)", "(c)", "(d)"},
		{
			"R",
			light(fbstatus.DACC.FbParametersKnown, false, false),
			light(fbstatus.IDCCb.FbParametersKnown, false, false),
			light(fbstatus.IDCCc.FbParametersKnown, false, false),
			light(fbstatus.IDCCd.FbParametersKnown, false, false),
		},
		{
			"D",
			light(fbstatus.DACC.FbParametersKnown, fbstatus.DACC.Spanning, fbstatus.DACC.Dfp),
			light(fbstatus.IDCCb.FbParametersKnown, false, fbstatus.IDCCb.Fallback),
			light(fbstatus.IDCCc.FbParametersKnown, false, fbstatus.IDCCc.Fallback),
			light(fbstatus.IDCCd.FbParametersKnown, false, fbstatus.IDCCd.Fallback),
		},
		{
			"ED",
			"",
			light(fbstatus.IDCCb.FbParametersKnown, fbstatus.IDCCb.EmptyDomain, false),
			light(fbstatus.IDCCc.FbParametersKnown, fbstatus.IDCCc.EmptyDomain, false),
			light(fbstatus.IDCCd.FbParametersKnown, fbstatus.IDCCd.EmptyDomain, false),
		},
		{
			"V",
			light(fbstatus.DACC.FbParametersKnown, fbstatus.DACC.IvaPresolved, fbstatus.DACC.IvaFallback),
			light(fbstatus.IDCCb.FbParametersKnown, fbstatus.IDCCb.Iva, fbstatus.IDCCb.IvaFallback),
			light(fbstatus.IDCCc.FbParametersKnown, fbstatus.IDCCc.Iva, fbstatus.IDCCc.IvaFallback),
			light(fbstatus.IDCCd.FbParametersKnown, fbstatus.IDCCd.Iva, fbstatus.IDCCd.IvaFallback),
		},
		{
			"R B",
			light(fbstatus.DACC.FbParametersKnown, fbstatus.DACC.ReturnedBranches, false),
			light(fbstatus.IDCCb.FbParametersKnown, fbstatus.IDCCb.ReturnedBranches, false),
			light(fbstatus.IDCCc.FbParametersKnown, fbstatus.IDCCc.ReturnedBranches, false),
			light(fbstatus.IDCCd.FbParametersKnown, fbstatus.IDCCd.ReturnedBranches, false),
		},
		{
			"A V",
			"",
			light(fbstatus.IDCCb.FbParametersKnown, fbstatus.IDCCb.AtcValidation, fbstatus.IDCCb.AtcValidationFallback),
			light(fbstatus.IDCCc.FbParametersKnown, fbstatus.IDCCc.AtcValidation, fbstatus.IDCCc.AtcValidationFallback),
			light(fbstatus.IDCCd.FbParametersKnown, fbstatus.IDCCd.AtcValidation, fbstatus.IDCCd.AtcValidationFallback),
		},
		{
			"A",
			"",
			light(fbstatus.IDCCb.AtcPublished, false, fbstatus.IDCCb.Fallback),
			light(fbstatus.IDCCc.AtcPublished, false, fbstatus.IDCCc.Fallback),
			light(fbstatus.IDCCd.AtcPublished, false, fbstatus.IDCCd.Fallback),
		},
	}
	return table
}
