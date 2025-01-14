package service

import (
	"XcxcPan/common/define"
	"os"
	"os/exec"
	"strings"
)

func CutFileForVideo(videoPath string) error {
	splitPath := strings.Split(videoPath, ".")
	dirPath := splitPath[0]
	splitDirPath := strings.Split(dirPath, "/")
	userId := splitDirPath[len(splitDirPath)-2]
	fileId := splitDirPath[len(splitDirPath)-1]

	os.MkdirAll(dirPath, os.ModePerm)
	// 输出 TS 文件模板
	outputTsTemplate := dirPath + "/" + userId + "_" + fileId + "_%03d.ts"
	// 输出 M3U8 文件名
	outputM3U8 := dirPath + "/" + define.M3U8_NAME
	// 构造 FFmpeg 命令
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath, // 输入文件
		"-c:v", "libx264", // 视频编码
		"-c:a", "aac", // 音频编码
		"-hls_time", "30", // 每个分片的时长
		"-hls_segment_filename", outputTsTemplate, // TS 文件模板
		outputM3U8, // 输出的 M3U8 文件
	)
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	}

	os.Remove(videoPath)
	return nil
}

func CreateThumbnailForVideo(videoPath string) error {
	// 输入视频文件
	splitPath := strings.Split(videoPath, ".")
	dirPath := splitPath[0]
	os.MkdirAll(dirPath, os.ModePerm)
	// 输出缩略图文件
	outputImage := dirPath + "/" + "thumbnail.jpg"
	// 构造 FFmpeg 命令
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath, // 输入文件
		"-ss", "0", // 从第 0 秒开始提取
		"-vframes", "1", // 只提取一帧
		"-q:v", "2", // 设置图像质量
		"-vf", "scale=150:-1", // 设置缩放的宽度和高度
		outputImage, // 输出文件名
	)
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func CreateThumbnailForImage(imagePath string) error {
	// 输入图片文件
	splitPath := strings.Split(imagePath, ".")
	dirPath := splitPath[0]
	os.MkdirAll(dirPath, os.ModePerm)
	// 输出缩略图文件
	outputImage := dirPath + "/" + "thumbnail.jpg"

	// FFmpeg 命令构建
	cmd := exec.Command(
		"ffmpeg",
		"-i", imagePath, // 输入文件
		"-vf", "scale=150:-1", // 缩放：宽度 100，高度按比例调整
		outputImage, // 输出文件
	)
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	}
	os.Remove(imagePath)
	return nil

}
