import { Injectable } from '@angular/core';
import { UpdatePotableWaterSalesConfigCommand } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { CustomValidator } from 'src/models/CustomValidator';

export class AddEditWaterPurifiedFields implements Omit<UpdatePotableWaterSalesConfigCommand, '$typeName' | '$unknown'> {
    waterPrice: number = 0
    availableSupply: number = 0
    minOrderQuantity: number = 0
    organizationId: bigint = null
    waterPriceMin: number = 0
    waterPricePerKm: number = 0
    bottleCost: number = 0
    isEnabled: boolean = false
}

@Injectable({
  providedIn: 'root'
})
export class AddEditWaterPurifiedValidator extends CustomValidator<AddEditWaterPurifiedFields> {
    constructor(){
        super(new AddEditWaterPurifiedFields());

        this.ruleFor('organizationId')
            .must(oid => oid > 0).withMessage("This field is required")

        this.ruleFor('waterPriceMin')
            .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        this.ruleFor('waterPricePerKm')
            .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        this.ruleFor('waterPrice')
            .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        this.ruleFor('bottleCost')
            .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        this.ruleFor('waterPriceMin')
            .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")
    }
}