package appy

import appy_logger "github.com/nfwGytautas/appy-go/logger"

func InitializeV2() error {
	appy_logger.Logger().Info("Initializing appy-go")
	return nil
}
