import { computed, effect, signal, WritableSignal } from "@angular/core";
import { ValidationErrors, Validator } from "fluentvalidation-ts";

export class CustomValidator<T> extends Validator<T> {
    readonly type: T
    readonly formValues = signal<T>(<T>{})
    readonly formErrors = signal<ValidationErrors<T>>(<ValidationErrors<T>>{})
    readonly isValid = computed(() => Object.keys(this.formErrors()).length === 0)
    private initialFormValues: T

    constructor(initialFormValues: T){
        super()
        this.initialFormValues = initialFormValues

        effect(() => {
            this.formErrors.set(this.validate(this.formValues()))
        })

        this.formValues.set(initialFormValues)
    }

    updateForm<K extends keyof T>(formField: K, newValue: T[K]) {
        this.formValues.update(prev => ({
            ...prev,
            [formField]: newValue
        }));
    }

    customFieldError(formField: keyof T, error: string) {
        this.formErrors.set({
            ...this.validate(this.formValues()),
            [formField]: error
        })
    }

    showFieldErrors(val: HTMLIonInputElement | HTMLIonSelectElement){
        setTimeout(() => {
            val.classList.add("ion-touched", "ion-invalid")
        }, 0)
    }

    reset(){
        this.formValues.set(this.initialFormValues)
    }
}