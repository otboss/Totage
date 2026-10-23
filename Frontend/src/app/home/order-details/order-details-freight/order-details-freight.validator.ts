import { Injectable } from '@angular/core';
import { CustomValidator } from 'src/models/CustomValidator';

export class OrderDetailsFreightFormFields {
    pickupAddress: string = ""
}

@Injectable({
  providedIn: 'root',
})
export class OrderDetailsFreightValidator extends CustomValidator<OrderDetailsFreightFormFields> {
    constructor(){
        super(new OrderDetailsFreightFormFields());

        this.ruleFor('pickupAddress')
            .notEmpty()
            .withMessage("Pickup address is required")
            .maxLength(255)
            .withMessage("Pickup address exceeds max length")
    }
}