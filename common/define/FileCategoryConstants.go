package define

var VIDEO = "video"
var MUSIC = "music"
var IMAGE = "image"
var DOC = "doc"
var OTHERS = "others"

var VIDEO_CATEGORY = map[string]int{
	VIDEO:  1,
	MUSIC:  2,
	IMAGE:  3,
	DOC:    4,
	OTHERS: 5,
}

func ExistsCategory(category string) bool {
	_, ok := VIDEO_CATEGORY[category]
	return ok
}
