package msgqueue

import (
	"fmt"
	"strings"
)

type MsgWriter interface {
	WriteKV(k, v []byte) error
	Close() error
	String() string
}

// kafka:broker1,broker2,broker3
// file:path/to/file
// os:stdout
func createMsgWriter(cfg string) (MsgWriter, error) {
	if cfg == "nop" {
		return NewNopMsgWriter(), nil
	} else if strings.HasPrefix(cfg, CfgPrefixKafka) {
		brokers := strings.TrimPrefix(cfg, CfgPrefixKafka)
		return NewKafkaMsgWriter(brokers)
	} else if strings.HasPrefix(cfg, CfgPrefixFile) {
		filePath := strings.TrimPrefix(cfg, CfgPrefixFile)
		return NewFileMsgWriter(filePath)
	} else if strings.TrimPrefix(cfg, CfgPrefixOS) == "stdout" {
		return NewStdOutMsgWriter(), nil
	} else if strings.HasPrefix(cfg, CfgPrefixDir) {
		dirPath := strings.TrimPrefix(cfg, CfgPrefixDir)
		return NewDirMsgWriter(dirPath)
	} else {
		return nil, fmt.Errorf("unsupported config: %s", cfg)
	}
}
