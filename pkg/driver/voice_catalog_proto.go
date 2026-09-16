package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func VoiceSpeechCapabilitiesFromProto(cap *wujiv1.VoiceSpeechCapabilities) *VoiceSpeechCapabilities {
	if cap == nil {
		return nil
	}
	return &VoiceSpeechCapabilities{
		Languages: cap.GetLanguages(), Emotions: cap.GetEmotions(), Styles: cap.GetStyles(),
		SpeedMin: cap.GetSpeedMin(), SpeedMax: cap.GetSpeedMax(),
		PitchMin: int(cap.GetPitchMin()), PitchMax: int(cap.GetPitchMax()),
		EnergySupported: cap.GetEnergySupported(),
		LanguagesSupported: cap.GetLanguagesSupported(),
		EmotionsSupported:  cap.GetEmotionsSupported(),
		StylesSupported:    cap.GetStylesSupported(),
	}
}

func VoiceSpeechCapabilitiesToProto(cap *VoiceSpeechCapabilities) *wujiv1.VoiceSpeechCapabilities {
	if cap == nil {
		return nil
	}
	return &wujiv1.VoiceSpeechCapabilities{
		Languages: cap.Languages, Emotions: cap.Emotions, Styles: cap.Styles,
		SpeedMin: cap.SpeedMin, SpeedMax: cap.SpeedMax,
		PitchMin: int32(cap.PitchMin), PitchMax: int32(cap.PitchMax),
		EnergySupported: cap.EnergySupported,
		LanguagesSupported: cap.LanguagesSupported,
		EmotionsSupported:  cap.EmotionsSupported,
		StylesSupported:    cap.StylesSupported,
	}
}

func VoiceProfileFromProto(p *wujiv1.VoiceProfile) VoiceProfileInfo {
	if p == nil {
		return VoiceProfileInfo{}
	}
	return VoiceProfileInfo{
		Name: p.GetName(), VoiceID: p.GetVoiceId(), Description: p.GetDescription(),
		Source: p.GetSource(), Speech: VoiceSpeechCapabilitiesFromProto(p.GetSpeech()),
	}
}

func VoiceProfileToProto(p VoiceProfileInfo) *wujiv1.VoiceProfile {
	return &wujiv1.VoiceProfile{
		Name: p.Name, VoiceId: p.VoiceID, Description: p.Description, Source: p.Source,
		Speech: VoiceSpeechCapabilitiesToProto(p.Speech),
	}
}

func ListVoiceProfilesResponseFromProto(resp *wujiv1.ListVoiceProfilesResponse) []VoiceProfileInfo {
	if resp == nil || len(resp.GetProfiles()) == 0 {
		return nil
	}
	out := make([]VoiceProfileInfo, 0, len(resp.GetProfiles()))
	for _, profile := range resp.GetProfiles() {
		out = append(out, VoiceProfileFromProto(profile))
	}
	return out
}

func ListVoiceProfilesResponseToProto(profiles []VoiceProfileInfo) *wujiv1.ListVoiceProfilesResponse {
	out := make([]*wujiv1.VoiceProfile, 0, len(profiles))
	for _, profile := range profiles {
		out = append(out, VoiceProfileToProto(profile))
	}
	return &wujiv1.ListVoiceProfilesResponse{Profiles: out}
}
