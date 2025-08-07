package image

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"os"
)

// CompressJPEG 压缩JPEG（输入输出为[]byte，不改变尺寸和格式）
// quality: 1-100，值越低压缩率越高
func CompressJPEG(input []byte, quality int) ([]byte, error) {
	// 解码字节流为图像对象
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// 编码为JPEG（控制质量）
	var outputBuf bytes.Buffer
	err = jpeg.Encode(&outputBuf, img, &jpeg.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, err
	}

	return outputBuf.Bytes(), nil
}

// CompressPNG 压缩PNG（输入输出为[]byte，不改变尺寸和格式）
// 使用imaging库重新编码，通过优化参数减小体积
func CompressPNG(input []byte) ([]byte, error) {
	// 解码PNG字节流
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// 创建缓冲区
	var outputBuf bytes.Buffer

	// 重新编码PNG，通过设置压缩级别优化体积
	// 注意：PNG是无损压缩，压缩级别影响速度而非画质
	err = imaging.Encode(&outputBuf, img, imaging.PNG, imaging.PNGCompressionLevel(6))
	if err != nil {
		return nil, err
	}

	return outputBuf.Bytes(), nil
}

// CompressImage 根据格式自动选择压缩方法
func CompressImage(input []byte, format string, quality int) ([]byte, error) {
	switch format {
	case "jpeg", "jpg":
		return CompressJPEG(input, quality)
		// 不知道为什么png的图片用这个库压缩后反而更大了，所以png图片就不压缩了
	case "png":
		//return CompressPNG(input)
		return input, nil
	default:
		return nil, os.ErrInvalid
	}
}
