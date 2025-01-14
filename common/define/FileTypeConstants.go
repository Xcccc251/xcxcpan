package define

var VIDEO_TYPE = []string{
	".mp4",
	".avi",
	".rmvb",
	".mkv",
	".mov",
}

var MUSIC_TYPE = []string{
	".mp3",
	".wav",
	".wma",
	".flac",
	".midi",
	".ra",
	".ape",
	".aac",
	".cda",
}

var IMAGE_TYPE = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".bmp",
	".webp",
	".ico",
	".svg",
	".tiff",
	".pdt",
	".psd",
	".xmp",
}

var PDF_TYPE = []string{
	".pdf",
}

var WORD_TYPE = []string{
	".doc",
	".docx",
	".xls",
	".xlsx",
	".ppt",
	".pptx",
}

var EXCEL_TYPE = []string{
	".xls",
	".xlsx",
}

var TXT_TYPE = []string{
	".txt",
}

var PROGRAME_TYPE = []string{
	".exe",
	".msi",
	".vue",
	".css",
	".scss",
	".class",
	".xml",
	".html",
	".dll",
	".jar",
	".py",
	".o",
	".sql",
	".java",
	".json",
	".c",
	".cpp",
	".h",
	".cs",
	".js",
	".ts",
	".php",
	".go",
	".rb",
	".sh",
	".bat",
	".ps1",
	".vbs",
	".bat",
	".cmd",
	".ps1",
	".vbs",
	".dll",
}

var ZIP_TYPE = []string{
	".zip",
	".rar",
	".7z",
	".tar",
	".gz",
	".bz2",
	".xz",
	".zst",
	".lzma",
	".lz",
	".lzo",
	".lz4",
	".lzh",
	".lha",
	".cab",
	".iso",
	".dmg",
	".vhd",
	".vhdx",
	".vmdk",
}

var OTHERS_TYPE = []string{
	"others",
}
var CODE_MAP = map[int][]string{
	1:  VIDEO_TYPE,
	2:  MUSIC_TYPE,
	3:  IMAGE_TYPE,
	4:  PDF_TYPE,
	5:  WORD_TYPE,
	6:  EXCEL_TYPE,
	7:  TXT_TYPE,
	8:  PROGRAME_TYPE,
	9:  ZIP_TYPE,
	10: OTHERS_TYPE,
}

var CODE_CATEGORY_MAP = map[int]string{
	1:  VIDEO,
	2:  MUSIC,
	3:  IMAGE,
	4:  DOC,
	5:  DOC,
	6:  DOC,
	7:  DOC,
	8:  OTHERS,
	9:  OTHERS,
	10: OTHERS,
}

func GetTypeCodeBySuffix(typeSuffix string) int {
	for k, v := range CODE_MAP {
		for _, vv := range v {
			if vv == typeSuffix {
				return k
			}
		}
	}
	return 10
}

func GetCategoryCodeBySuffix(typeSuffix string) int {
	typeCode := GetTypeCodeBySuffix(typeSuffix)
	category := CODE_CATEGORY_MAP[typeCode]
	return GetCategoryCodeByCategory(category)
}
