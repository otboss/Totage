import { Injectable } from '@angular/core';
import { UpdateFreightConfigCommand, UpdatePotableWaterSalesConfigCommand } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { CustomValidator } from 'src/models/CustomValidator';

export class AddEditFreightFields implements Omit<UpdateFreightConfigCommand, '$typeName' | '$unknown'> {
    organizationId: bigint = null;
    solidWasteDumpSiteLatitude?: number = null;
    solidWasteDumpSiteLongitude?: number = null;
    liquidWasteDumpSiteLatitude?: number = null;
    liquidWasteDumpSiteLongitude?: number = null;
    isEnabled: boolean = false;
}

@Injectable({
  providedIn: 'root'
})
export class AddEditFreightValidator extends CustomValidator<AddEditFreightFields> {
    constructor(){
        super(new AddEditFreightFields());

        this.ruleFor('organizationId')
            .must(oid => oid > 0).withMessage("This field is required")

        // this.ruleFor('solidWasteDumpSiteLatitude')
        //     .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        // this.ruleFor('solidWasteDumpSiteLongitude')
        //     .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        // this.ruleFor('liquidWasteDumpSiteLatitude')
        //     .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")

        // this.ruleFor('liquidWasteDumpSiteLongitude')
        //     .must(fpm => isNaN(parseInt(fpm.toString()))).withMessage("This field is required")
    }
}