import { CustomValidator } from 'src/models/CustomValidator';

export interface LoginFormFields {
    email: string
    password: string
}

export class LoginValidator extends CustomValidator<LoginFormFields> {
    constructor(){
        super();

        this.ruleFor('email').emailAddress().withMessage("invalid email address provided")
    }
}