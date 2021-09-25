package datelogzap

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(fileName string) (*zap.SugaredLogger, error) {
	now := time.Now()
	path := fmt.Sprintf("runtime/logs/%s/", now.Format("2006-01-02"))

	if !IfNotExistMKDir(path) {
		err := errors.New("create dir error , " + path)
		return nil, err
	}

	encoder := getEncoder()
	file, err := os.OpenFile(path+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		err := errors.New("create log file error  " + path + fileName)
		return nil, err
	}
	writer := zapcore.AddSync(file)
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}

//目录不存在就创建目录
func IfNotExistMKDir(path string) bool {
	if !DirCheck(path) {
		err := os.MkdirAll(path, os.ModePerm)
		return err == nil
	}
	return true
}

//检测目录是否存在
func DirCheck(path string) bool {

	info, err := os.Stat(path)

	if err != nil {
		return false
	}
	return info.IsDir()
}

//获取日志的编码
func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
