package handler

import (
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
)

type SampleHandler struct{}

func (h SampleHandler) Exec(msg types.Message, dist types.SessionManager) error {
	fmt.Printf("%s", msg.GetBody())
	return nil
}

func (h SampleHandler) GetResultQueueConfig() *types.QueueConfig {
	return nil
}

func (h SampleHandler) GetResultStorageConfig() *types.StorageConfig {
	return nil
}
