import { Injectable } from '@angular/core';
import { UpdateIrrigationWaterSalesConfigCommand, UpdateIrrigationWaterSalesConfigCommandSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { ProtoValidator } from 'src/models/ProtoValidator';
// import { CustomValidator } from 'src/models/CustomValidator';

@Injectable({
  providedIn: 'root'
})
export class AddEditWaterUnpurifiedValidator extends ProtoValidator<UpdateIrrigationWaterSalesConfigCommand> {
    constructor(){
        super(UpdateIrrigationWaterSalesConfigCommandSchema);
    }
}

// export class AddEditWaterUnpurifiedFields implements Omit<UpdateIrrigationWaterSalesConfigCommand, '$typeName' | '$unknown'> {
//     waterPrice: number = 0;
//     availableSupply: number = 0
//     minOrderQuantity: number = 0
//     organizationId: bigint = null
//     waterPriceMin: number = 0
//     waterPricePerKm: number = 0
//     barrelCost?: number = 0
//     isEnabled: boolean = false
// }

// @Injectable({
//   providedIn: 'root'
// })
// export class AddEditWaterUnpurifiedValidator extends CustomValidator<AddEditWaterUnpurifiedFields> {
//     constructor(){
//         super(new AddEditWaterUnpurifiedFields());

//         this.ruleFor('organizationId')
//             .must(oid => oid > 0).withMessage("This field is required")

//         this.ruleFor('waterPriceMin')
//             .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

//         this.ruleFor('waterPricePerKm')
//             .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

//         this.ruleFor('waterPrice')
//             .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

//         // this.ruleFor('barrelCost')
//         //     .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

//         this.ruleFor('waterPriceMin')
//             .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")
//     }
// }