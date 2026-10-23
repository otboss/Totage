import { OrderType } from "src/app/shared/services/client/shared/api/order_type_pb"
import { Assets } from "src/constants/Assets"

export const OrderTypeCards: Record<keyof typeof OrderType, {
    orderType: OrderType
    label: string
    image: string
  }> = {
    VEHICLE_RENTAL: {
        orderType: OrderType.VEHICLE_RENTAL,
        image: Assets.Images.rentalVehicle,
        label: "Vehicle Rental"
    },
    FREIGHT: {
        orderType: OrderType.FREIGHT,
        image: Assets.Images.truck,
        label: "Moving"
    },
    WATER_PURIFIED: {
        orderType: OrderType.WATER_PURIFIED,
        image: Assets.Images.water,
        label: "Drinking Water"
    },
    WATER_UNPURIFIED: {
        orderType: OrderType.WATER_UNPURIFIED,
        image: Assets.Images.waterIrrigation,
        label: "Irrigation Water"
    },
    ICE: {
        orderType: OrderType.ICE,
        image: Assets.Images.ice,
        label: "Ice"
    },
    EGG: {
        orderType: OrderType.EGG,
        image: Assets.Images.eggs,
        label: "Eggs"
    },
    COCONUT: {
        orderType: OrderType.COCONUT,
        image: Assets.Images.coconuts,
        label: "Coconuts"
    },
    COOKING_GAS: {
        orderType: OrderType.COOKING_GAS,
        image: Assets.Images.gasCylinder,
        label: "Cooking Gas"
    },
    HONEY: {
        orderType: OrderType.HONEY,
        image: Assets.Images.honeyJar,
        label: "Honey"
    },
    BLOCK: {
        orderType: OrderType.BLOCK,
        image: Assets.Images.cinderBlocks,
        label: "Blocks"
    },
    CHARCOAL: {
        orderType: OrderType.CHARCOAL,
        image: Assets.Images.charcoalPile,
        label: "Charcoal"
    },
    CONSTRUCTION: {
        orderType: OrderType.CONSTRUCTION,
        label: "",
        image: ""
    },
    CESSPOOL: {
        orderType: OrderType.CESSPOOL,
        image: Assets.Images.cespoolTruck,
        label: "Cesspool"
    },
  }