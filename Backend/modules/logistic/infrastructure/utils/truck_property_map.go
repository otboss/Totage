package utils

// import (
// 	"errors"

// 	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
// 	common "waterordersystem/GOBACKEND/shared/api"
// )

// var TruckPropertyMap = &truckPropertyMap{}

// type truckPropertyMap struct{}

// func (s truckPropertyMap) ProtoTruckPropertyToDbTruckProperty(truckProperty common.FreightTruckProperty) (db_logistic.LogisticFreightTruckProperty, error) {
// 	switch truckProperty {
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_BUCKET:
// 		return db_logistic.LogisticFreightTruckPropertyBUCKET, nil
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_TIPPER:
// 		return db_logistic.LogisticFreightTruckPropertyTIPPER, nil
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_FREEZER:
// 		return db_logistic.LogisticFreightTruckPropertyFREEZER, nil
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_REFRIGERATED:
// 		return db_logistic.LogisticFreightTruckPropertyREFRIGERATED, nil
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_CESSPOOL:
// 		return db_logistic.LogisticFreightTruckPropertyCESSPOOL, nil
// 	case common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_CRANE:
// 		return db_logistic.LogisticFreightTruckPropertyCRANE, nil
// 	default:
// 		return db_logistic.LogisticFreightTruckPropertyNONE, errors.New("invalid parameter provided")
// 	}
// }

// func (s truckPropertyMap) DbTruckPropertyToProtoTruckProperty(truckProperty db_logistic.LogisticFreightTruckProperty) (common.FreightTruckProperty, error) {
// 	switch truckProperty {
// 	case db_logistic.LogisticFreightTruckPropertyBUCKET:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_BUCKET, nil
// 	case db_logistic.LogisticFreightTruckPropertyTIPPER:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_TIPPER, nil
// 	case db_logistic.LogisticFreightTruckPropertyFREEZER:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_FREEZER, nil
// 	case db_logistic.LogisticFreightTruckPropertyREFRIGERATED:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_REFRIGERATED, nil
// 	case db_logistic.LogisticFreightTruckPropertyCESSPOOL:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_CESSPOOL, nil
// 	case db_logistic.LogisticFreightTruckPropertyCRANE:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_CRANE, nil
// 	default:
// 		return common.FreightTruckProperty_FREIGHT_TRUCK_PROPERTY_UNSPECIFIED, errors.New("invalid parameter provided")
// 	}
// }
