import { Injectable } from "@angular/core";
import { CreateAccountCommand } from "src/app/shared/services/client/modules/organization/api/organization_pb";
import { AccountRoles } from "src/app/shared/services/client/shared/api/account_role_pb";
import { CustomValidator } from "src/models/CustomValidator";


export class AddEditAccountFormFields implements Omit<CreateAccountCommand, '$typeName' | '$unknown'> {
    email: string = ""
    organizationId: bigint = null
    roles?: AccountRoles
    isAdmin: boolean = false
    isDriver: boolean = false
    isCustomerSupport: boolean = false   
}


@Injectable({
  providedIn: 'root'
})
export class AddEditAccountValidator extends CustomValidator<AddEditAccountFormFields> {
    constructor(){
        super(new AddEditAccountFormFields());

        // TODO: Add form field rules
    }
}