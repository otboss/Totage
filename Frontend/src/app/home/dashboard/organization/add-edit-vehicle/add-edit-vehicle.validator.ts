import { Injectable } from "@angular/core";
import { CreateRentalVehicleCommand, CreateRentalVehicleCommandSchema, CreateTruckCommand, CreateTruckCommandSchema, CreateVehicleCommand, CreateVehicleCommandSchema, SetJobTimeslotCommand, SetJobTimeslotCommandSchema } from "src/app/shared/services/client/modules/logistic/api/logistic_pb";
import { SetEnrollmentCommand, SetEnrollmentCommandSchema } from "src/app/shared/services/client/modules/order/api/order_pb";
// import { CreateJobTimeslotCommand, CreateJobTimeslotCommandSchema } from "src/app/shared/services/client/modules/organization/api/organization_pb";
import { ProtoValidator } from "src/models/ProtoValidator";

@Injectable({
  providedIn: 'root'
})
export class AddEditVehicleValidator extends ProtoValidator<CreateVehicleCommand> {
    constructor(){
        super(CreateVehicleCommandSchema);
    }
}

@Injectable({
  providedIn: 'root'
})
export class AddEditTruckValidator extends ProtoValidator<CreateTruckCommand> {
    constructor(){
        super(CreateTruckCommandSchema);
    }
}

@Injectable({
  providedIn: 'root'
})
export class AddEditRentalVehicleValidator extends ProtoValidator<CreateRentalVehicleCommand> {
    constructor(){
        super(CreateRentalVehicleCommandSchema);
    }
}

@Injectable({
  providedIn: 'root'
})
export class SetEnrollmentValidator extends ProtoValidator<SetEnrollmentCommand>{
    constructor(){
        super(SetEnrollmentCommandSchema);
    }
}

@Injectable({
  providedIn: 'root'
})
export class SetJobTimeslotValidator extends ProtoValidator<SetJobTimeslotCommand>{
    constructor(){
        super(SetJobTimeslotCommandSchema);
    }
}

/*
export class AddEditTruckFormFields implements Omit<Required<TruckCreate>, '$typeName' | '$unknown'> {
    truckId: bigint = BigInt(0)
    bedSize: TruckCreate_BedSize
    hasFreezer: boolean = false
    hasTipper: boolean = false
    hasRefrigerator: boolean = false
    hasCesspool: boolean = false
    hasBucket: boolean = false
    hasCrane: boolean = false
    freightPricePerKm: number = null
    cost: number = null
    licensePlateNumber: string = ""
    lengthMm: number = null
    widthMm: number = null
    heightMm?: number | undefined = null
    isPublic: boolean = false
}

export class AddEditRentalVehicleFormFields implements Omit<Required<RentalVehicle>, '$typeName' | '$unknown'> {
    driverAccountId: bigint = null;
    rentalVehicleId: bigint = null
    organizationId: bigint = null
    dailyRentalCost: number = null
    driverCost: number = null
    isDriverMandatory: boolean = false
    minDaysRental: bigint = null
    isPublic: boolean = false
}

export class AddVehicleFormFields implements Omit<Required<CreateVehicleCommand>, '$typeName' | '$unknown'> {
    driverAccountId: string;
    vehicleId: bigint;
    chassisNumber: string = ""
    licensePlateNumber: string = ""
    make: string = ""
    model: string = ""
    fuelType: FuelType = null
    organizationId: bigint = BigInt(0)
    year: bigint = BigInt(0)
    seatingCapacity: bigint = BigInt(0)
    isManualTransmission: boolean = false
}

@Injectable({
  providedIn: 'root'
})
export class AddEditVehicleValidator extends CustomValidator<AddVehicleFormFields> {
    constructor(){
        super(new AddVehicleFormFields());

        this.ruleFor("chassisNumber")
            .notEmpty().withMessage("This field is required")

        this.ruleFor("licensePlateNumber")
            .notEmpty().withMessage("This field is required")

        this.ruleFor("make")
            .notEmpty().withMessage("This field is required")

        this.ruleFor("model")
            .notEmpty().withMessage("This field is required")

        this.ruleFor("fuelType")
            .must(ft => ft != null).withMessage("This field is required")

        this.ruleFor("seatingCapacity")
            .must(sc => sc > 0).withMessage("This field is required")
            .must(sc => sc > 1).withMessage("Invalid seating capacity")

        this.ruleFor("year")
            .must(yr => yr > 0).withMessage("This field is required")
            .must(yr => yr >= 1900).withMessage("Invalid year provided")

        this.ruleFor("driverAccountId")
            .must(d => d != null).withMessage("This field is required")
    }
}

@Injectable({
  providedIn: 'root'
})
export class AddEditTruckValidator extends CustomValidator<AddEditTruckFormFields> {
    constructor(){
        super(new AddEditTruckFormFields());

        this.ruleFor('lengthMm')
            .must(l => l != null && l.toString() != "").withMessage("This field is required")
            .greaterThan(0).withMessage("Invalid length")

        this.ruleFor('widthMm')
            .must(l => l != null && l.toString() != "").withMessage("This field is required")
            .greaterThan(0).withMessage("Invalid width")

        this.ruleFor('heightMm')
            .must(l =>  l == null || l.toString() == "" ? true : l > 0).withMessage("Invalid height")

        this.ruleFor('cost')
            .must(c => c != null && c.toString() != "").withMessage("This field is required")

        this.ruleFor('freightPricePerKm')
            .must(c => c != null && c.toString() != "").withMessage("This field is required")

    }
}

@Injectable({
  providedIn: 'root'
})
export class AddEditRentalVehicleValidator extends CustomValidator<AddEditRentalVehicleFormFields> {
    constructor(){
        super(new AddEditRentalVehicleFormFields());

        // this.ruleFor('dailyRentalCost')
        //     .must(d => !!d.toString() && d.toString().trim().length > 0).withMessage("This field is required")
        //     .must(d => d >= 0).withMessage("Invalid daily rental cost")

        // this.ruleFor('driverCost')
        //     .must(d => !!d.toString() && d.toString().trim().length > 0).withMessage("This field is required")
        //     .must(d => d >= 0).withMessage("Invalid driver cost")

        this.ruleFor('minDaysRental')
            .must(l => l != null && l.toString() != "").withMessage("This field is required")
            .must(l => l.toString() != "" && l >= 1).withMessage("Invalid min days value")

        this.ruleFor('dailyRentalCost')
            .must(l => l != null && l.toString() != "").withMessage("This field is required")
    }
}**/