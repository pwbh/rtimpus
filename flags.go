package rtimpus

const (
	AUDIO_CODEC_SUPPORT_SND_NONE    = 0x0001 // Raw sound, no compression
	AUDIO_CODEC_SUPPORT_SND_ADPCM   = 0x0002 // ADPCM compression
	AUDIO_CODEC_SUPPORT_SND_MP3     = 0x0004 //  mp3 compression
	AUDIO_CODEC_SUPPORT_SND_INTEL   = 0x0008 // Not used
	AUDIO_CODEC_SUPPORT_SND_UNUSED  = 0x0010 // Not used
	AUDIO_CODEC_SUPPORT_SND_NELLY8  = 0x0020 // NellyMoser at 8-kHz
	AUDIO_CODEC_SUPPORT_SND_NELLY   = 0x0040 // NellyMoser compression (5, 11, 22 and 55 kHz)
	AUDIO_CODEC_SUPPORT_SND_G711A   = 0x0080 // G711A sound compression (Flash Media Server only)
	AUDIO_CODEC_SUPPORT_SND_G711U   = 0x0100 // G711U sound compression (Flash Media Server only)
	AUDIO_CODEC_SUPPORT_SND_NELLY16 = 0x0200 // NellyMouse at 16-kHz compression
	AUDIO_CODEC_SUPPORT_SND_AAC     = 0x0400 // Advanced audio coding (AAC) codec
	AUDIO_CODEC_SUPPORT_SND_SPEEX   = 0x0800 // Speex Audio
	AUDIO_CODEC_SUPPORT_SND_ALL     = 0x0FFF // All RTMP-supported audio codecs
)

const (
	VIDEO_CODEC_SUPPORT_VID_UNUSED    = 0x0001 // Obsolete value
	VIDEO_CODEC_SUPPORT_VID_JPEG      = 0x0002 // Obsolete value
	VIDEO_CODEC_SUPPORT_VID_SORENSON  = 0x0004 // Sorenson Flash video
	VIDEO_CODEC_SUPPORT_VID_HOMEBREW  = 0x0008 // V1 screen sharing
	VIDEO_CODEC_SUPPORT_VID_VP6       = 0x0010 // On2 video (Flash 8+)
	VIDEO_CODEC_SUPPORT_VID_VP6ALPHA  = 0x0020 // On2 video with alpha channel
	VIDEO_CODEC_SUPPORT_VID_HOMEBREWV = 0x0040 // Screen sharing version 2 (Flash 8+)
	VIDEO_CODEC_SUPPORT_VID_H264      = 0x0080 // H264 video
	VIDEO_CODEC_SUPPORT_VID_ALL       = 0x00FF // All RTMP-supported video codecs
)

const (
	VIDEO_FUNCTION_SUPPORT_VID_CLIENT_SEEK = 1 // Indicates that the client can perform frame-accurate seeks
)

const (
	OBJECT_ENCODING_AMF0 = 0 // AMF0 object encoding supported by Flash 6 and later
	OBJECT_ENCODING_AMF3 = 3 // AMF3 encoding from Flash 9 (AS3)
)
