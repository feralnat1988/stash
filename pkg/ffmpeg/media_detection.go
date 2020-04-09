package ffmpeg

import (
	"bytes"
	"github.com/stashapp/stash/pkg/logger"
	"os"
)

// detect file format from magic file number
// https://github.com/lex-r/filetype/blob/73c10ad714e3b8ecf5cd1564c882ed6d440d5c2d/matchers/video.go

func mkv(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3 &&
		containsMatroskaSignature(buf, []byte{'m', 'a', 't', 'r', 'o', 's', 'k', 'a'})
}

func webm(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3 &&
		containsMatroskaSignature(buf, []byte{'w', 'e', 'b', 'm'})
}

func containsMatroskaSignature(buf, subType []byte) bool {
	limit := 4096
	if len(buf) < limit {
		limit = len(buf)
	}

	index := bytes.Index(buf[:limit], subType)
	if index < 3 {
		return false
	}

	return buf[index-3] == 0x42 && buf[index-2] == 0x82
}

//returns container as string ("" on error or no match)
//implements only mkv or webm as ffprobe can't distinguish between them
//and not all browsers support mkv
func MagicContainer(file_path string) Container {
	file, err := os.Open(file_path)
	if err != nil {
		logger.Errorf("[magicfile] %v", err)
		return ""
	}

	defer file.Close()

	buf := make([]byte, 4096)
	_, err = file.Read(buf)
	if err != nil {
		logger.Errorf("[magicfile] %v", err)
		return ""
	}

	if webm(buf) {
		return Webm
	}
	if mkv(buf) {
		return Matroska
	}
	return ""
}
