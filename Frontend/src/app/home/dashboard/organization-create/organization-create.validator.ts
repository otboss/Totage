import { Injectable } from "@angular/core";
import { CreateOrganizationCommand, CreateOrganizationCommandSchema } from "src/app/shared/services/client/modules/organization/api/organization_pb";
import { ProtoValidator } from "src/models/ProtoValidator";

@Injectable({
  providedIn: 'root'
})
export class OrganizationCreateValidator extends ProtoValidator<CreateOrganizationCommand> {
    constructor(){
        super(CreateOrganizationCommandSchema);
    }
}


// import { Injectable } from '@angular/core';
// import { UnknownField } from '@bufbuild/protobuf';
// import { Validator } from 'fluentvalidation-ts'
// import { CreateOrganizationCommand, CreateOrganizationCommand_BankAccountConfig, CreateOrganizationCommand_BlockDeliveryConfig, CreateOrganizationCommand_CoconutDeliveryConfig, CreateOrganizationCommand_EggDeliveryConfig, CreateOrganizationCommand_FreightConfig, CreateOrganizationCommand_IrrigationWaterDeliveryConfig, CreateOrganizationCommand_PotableWaterDeliveryConfig, PaypalMerchantDetails } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
// import { CountryCodeAlpha3 } from 'src/app/shared/services/client/shared/api/country_pb';
// import { Currency } from 'src/app/shared/services/client/shared/api/currency_pb';
// import { CustomValidator } from 'src/models/CustomValidator';

// export class OrganizationCreateFormFields implements Omit<CreateOrganizationCommand, '$typeName' | '$unknown'> {
//     organizationId?: bigint;
//     organizationName: string = "";
//     address: string = "";
//     phone: string = "";
//     email: string = "";
//     country: CountryCodeAlpha3 = CountryCodeAlpha3.COUNTRY_CODE_UNSPECIFIED
//     currency: Currency  = Currency.CURRENCY_CODE_UNSPECIFIED;
//     dutyFreeKm: number = 12;
//     freightConfig?: CreateOrganizationCommand_FreightConfig;
//     potableWaterDeliveryConfig?: CreateOrganizationCommand_PotableWaterDeliveryConfig;
//     irrigationWaterDeliveryConfig?: CreateOrganizationCommand_IrrigationWaterDeliveryConfig;
//     blocksDeliveryConfig?: CreateOrganizationCommand_BlockDeliveryConfig;
//     eggDeliveryConfig?: CreateOrganizationCommand_EggDeliveryConfig;
//     coconutDeliveryConfig?: CreateOrganizationCommand_CoconutDeliveryConfig;
//     maxOrderQueueSize?: bigint;
//     submitForApproval?: boolean;
//     paypalMerchantDetails?: PaypalMerchantDetails;
//     bankAccountDetails: CreateOrganizationCommand_BankAccountConfig[];
// }

// @Injectable({
//   providedIn: 'root'
// })
// export class OrganizationCreateValidator extends CustomValidator<OrganizationCreateFormFields> {
//     constructor(){
//         super(new OrganizationCreateFormFields());

//         this.ruleFor('email')
//             .notEmpty().withMessage("This field is required")
//             .emailAddress().withMessage("Invalid email address provided!")
//             .maxLength(255).withMessage("Max length exceeded")
//             // .matches(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/).withMessage("Invalid email address provided!")
//         this.ruleFor('organizationName')
//             .notEmpty().withMessage("This field is required")
//             .maxLength(255).withMessage("Max length exceeded")
//         this.ruleFor('phone')
//             .notEmpty().withMessage("This field is required")
//             .maxLength(20).withMessage("Max length exceeded")
//         this.ruleFor('address')
//             .notEmpty().withMessage("This field is required")
//             .maxLength(255).withMessage("Max length exceeded")
//         this.ruleFor('currency')
//             .notEqual("" as any).withMessage("This field is required")
//             // .equal(Currency.CURRENCY_CODE_UNSPECIFIED).withMessage("This field is required")
//         this.ruleFor('country')
//             .notEqual("" as any).withMessage("This field is required")
//             // .equal(CountryCodeAlpha3.COUNTRY_CODE_UNSPECIFIED).withMessage("This field is required")
//     }
// }