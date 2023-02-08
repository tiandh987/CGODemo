package ptz

//type limitRepo interface {
//	Get(def bool) (*dsd.Limit, error)
//	Set(cfg *dsd.Limit) error
//	Default() error
//}
//
//type limitUseCase struct{}
//
//var _ limitRepo = (*limitUseCase)(nil)
//
////func NewLimit() limitRepo {
////	return &limitUseCase{}
////}
//
//func (l *limitUseCase) Get(def bool) (*dsd.Limit, error) {
//	limit := dsd.NewLimit()
//
//	if !def {
//		if err := config.GetConfig(limit.ConfigKey(), &limit); err != nil {
//			return nil, err
//		}
//	}
//
//	return limit, nil
//}
//
//func (l *limitUseCase) Set(cfg *dsd.Limit) error {
//	position, err := _blpInstance.Ptz.Position()
//	if err != nil {
//		return err
//	}
//
//	if cfg.CheckLeft == 1 {
//		cfg.LeftBoundary = position.Pan
//	} else if cfg.CheckRight == 1 {
//		cfg.RightBoundary = position.Pan
//	} else if cfg.CheckUp == 1 {
//		cfg.UpBoundary = position.Tilt
//	} else if cfg.CheckDown == 1 {
//		cfg.DownBoundary = position.Tilt
//	}
//
//	if cfg.LevelEnable || cfg.VerticalEnable {
//		// TODO 线性扫描、区域扫描 恢复默认配置
//
//		//mgr.LinearScans.Default()
//		//config.SetConfig(mgr.LinearScans)
//		//err := mgr.SetDefaultCfgRegionScans()
//		//if err != nil {
//		//	return false
//		//}
//	}
//
//	if err := config.SetConfig(cfg.ConfigKey(), cfg); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (l *limitUseCase) Default() error {
//	limit := dsd.NewLimit()
//
//	if err := config.SetConfig(limit.ConfigKey(), limit); err != nil {
//		return err
//	}
//
//	return nil
//}
