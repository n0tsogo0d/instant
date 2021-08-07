package ffmpeg

import (
	"encoding/json"
	"os/exec"
)

type (
	Data struct {
		Streams []Stream `json:"streams"`
		Format  struct {
			Filename       string `json:"filename"`
			NBStreams      int    `json:"nb_streams"`
			NBPrograms     int    `json:"nb_programs"`
			FormatName     string `json:"format_name"`
			FormatLongName string `json:"format_long_name"`
			StartTime      string `json:"start_time"`
			Duration       string `json:"duration"`
			Size           string `json:"size"`
			BitRate        string `json:"bit_rate"`
			ProbeScore     int    `json:"probe_score"`
			Tags           struct {
				Encoder      string `json:"encoder"`
				CreationTime string `json:"creation_time"`
			} `json:"tags"`
		} `json:"format"`
	}
	Stream struct {
		Index              int    `json:"index"`
		CodecName          string `json:"codec_name"`
		CodecLongName      string `json:"codec_long_name"`
		Profile            string `json:"profile"`
		CodecType          string `json:"codec_type"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTag           string `json:"codec_tag"`
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		CodedWidth         int    `json:"coded_width"`
		CodedHeight        int    `json:"coded_height"`
		ClosedCaptions     int    `json:"closed_captions"`
		HasBFrames         int    `json:"has_b_frames"`
		SampleAspectRatio  string `json:"sample_aspect_ratio"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		PixFmt             string `json:"pix_fmt"`
		Level              int    `json:"level"`
		ColorRange         string `json:"color_range"`
		ColorSpace         string `json:"color_space"`
		ColorTransfer      string `json:"color_transfer"`
		ColorPrimaries     string `json:"color_primaries"`
		ChromaLocation     string `json:"chroma_location"`
		FieldOrder         string `json:"field_order"`
		Refs               int    `json:"refs"`
		IsAVC              string `json:"is_avc"`
		NalLengthSize      string `json:"nal_length_size"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		TimeBase           string `json:"time_base"`
		StartPts           int    `json:"start_pts"`
		StartTime          string `json:"start_time"`
		BitsPerRawSample   string `json:"bits_per_raw_sample"`
		SampleFmt          string `json:"sample_fmt"`
		SampleRate         string `json:"sample_rate"`
		Channels           int    `json:"channels"`
		ChannelLayout      string `json:"channel_layout"`
		BitsPerSample      int    `json:"bits_per_sample"`
		DmixMode           string `json:"dmix_mode"`
		LtrtCmixlev        string `json:"ltrt_cmixlev"`
		LtrtSurmixlev      string `json:"ltrt_surmixlev"`
		LoroCmixlev        string `json:"loro_cmixlev"`
		LoroSurmixlev      string `json:"loro_surmixlev"`
		DurationTs         int    `json:"duration_ts"`
		Duration           string `json:"duration"`
		Disposition        struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			Language                    string `json:"language"`
			Title                       string `json:"title"`
			Bps                         string `json:"BPS"`
			Duration                    string `json:"DURATION"`
			NumberOfFrames              string `json:"NUMBER_OF_FRAMES"`
			NumberOfBytes               string `json:"NUMBER_OF_BYTES"`
			STATISTICS_WRITING_APP      string `json:"_STATISTICS_WRITING_APP"`
			STATISTICS_WRITING_DATE_UTC string `json:"_STATISTICS_WRITING_DATE_UTC"`
			STATISTICS_TAGS             string `json:"_STATISTICS_TAGS"`
		} `json:"tags"`
	}
)

func Probe(URL string) (*Data, error) {
	var data Data

	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format",
		"json", "-show_format", "-show_streams", URL).Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &data)
	return &data, err
}
