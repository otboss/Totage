import { computed, effect, signal } from "@angular/core"
import { clone, create, Message } from "@bufbuild/protobuf"
import { GenMessage } from "@bufbuild/protobuf/codegenv2"
import { createValidator } from "@bufbuild/protovalidate"

export class ProtoValidator<T extends Message> {
    readonly type: T
    readonly formValues = signal<T>(<T>{})
    readonly formErrors = signal<Partial<DeepStringify<T>>>({})
    readonly isValid = computed(() => Object.keys(this.formErrors()).length === 0)
    private protoValidator = createValidator()
    private initialFormValues: T
    private schema: GenMessage<T>

    constructor(schema: GenMessage<T>){
        this.schema = schema
        this.initialFormValues = this.createDeepDefault(schema)
        this.reset()
        effect(() => {
            this.validate(this.formValues())
        })
    }

    private createDeepDefault(schema: any): any {
        const instance = create(schema) as any

        const fields = schema.fields && typeof schema.fields.list === "function" 
            ? schema.fields.list() 
            : (schema.fields || [])

        for (const field of fields) {
            const isMessage = field.kind === "message" || field.fieldKind === "message"

            if (isMessage) {
                if (field.repeated) {
                    instance[field.localName] = []
                } else if (field.message) {
                    instance[field.localName] = this.createDeepDefault(field.message)
                }
            }
        }

        return instance
    }

    updateForm<P extends NestedPaths<T>>(path: P, newValue: PathType<T, P>) {
        this.formValues.update(prev => {
            // 2. Deep clone the protobuf message correctly
            const cloned = clone(this.schema, prev)
            const segments = path.split('.')
            let current: any = cloned

            // Safely traverse down to the parent of the leaf node
            for (let i = 0; i < segments.length - 1; i++) {
                const segment = segments[i]
                current = current[segment]
            }

            // Assign the new value directly to the cloned reference
            current[segments[segments.length - 1]] = newValue
            return cloned
        })
    }

    showFieldErrors(val: HTMLIonInputElement | HTMLIonSelectElement){
        setTimeout(() => {
            val.classList.add("ion-touched", "ion-invalid")
        }, 0)
    }

    customFieldError<P extends NestedPaths<T>>(path: P, error: string) {
        // 1. Run validation first to establish the latest baseline automated errors
        this.validate(this.formValues())

        // 2. Inject the custom error down the nested path structure
        this.formErrors.update(prev => {
            // Deep clone the error object to avoid direct state mutation side-effects
            const cloned = JSON.parse(JSON.stringify(prev));
            const segments = path.split('.')
            let current: any = cloned;

            for (let i = 0; i < segments.length - 1; i++) {
            const segment = segments[i];
            const nextSegment = segments[i + 1];
            const isNextSegmentIndex = !isNaN(Number(nextSegment));

            // Dynamically initialize missing paths as either arrays or objects
            if (!current[segment]) {
                current[segment] = isNextSegmentIndex ? [] : {};
            }
            current = current[segment]
            }

            // Assign the error string to the target leaf node
            current[segments[segments.length - 1]] = error
            return cloned
        })
    }

    reset(): void {
        this.formValues.set(create(this.schema, this.initialFormValues))
    }

    private validate(values: T): void {
        const errors: any = {}

        try {
            const result = this.protoValidator.validate(this.schema, values)

            if (result.kind === "invalid") {
                for (const violation of result.violations) {
                    let pathSegments: string[] = []
                    
                    if ((violation as any).fieldPath) {
                        const rawPath = (violation as any).fieldPath;
                        pathSegments = rawPath
                            .replace(/\[(\d+)\]/g, '.$1')
                            .split('.')
                            .filter(Boolean)
                    } else if (violation.field && Array.isArray(violation.field)) {
                        for (const f of violation.field) {
                            const part = f as any; 
                            if (part.name) pathSegments.push(part.name)
                            if (part.kind === "list_sub" && typeof part.index === "number") {
                                pathSegments.push(String(part.index))
                            } else if (part.index !== undefined && part.index !== null) {
                                pathSegments.push(String(part.index))
                            }
                        }
                    }

                    if (pathSegments.length > 0) {
                        const message = violation.message || "Invalid field content"
                        
                        let current = errors
                        for (let i = 0; i < pathSegments.length; i++) {
                            const segment = pathSegments[i]
                            
                            if (i === pathSegments.length - 1) {
                                current[segment] = message
                            } else {
                                const nextSegment = pathSegments[i + 1]
                                const isNextSegmentIndex = !isNaN(Number(nextSegment))

                                if (!current[segment]) {
                                    current[segment] = isNextSegmentIndex ? [] : {}
                                }
                                current = current[segment]
                            }
                        }
                    }
                }
            } else if (result.kind === "error") {
                console.error("Compilation or runtime error inside evaluation:", result.message)
            }
        } catch (err) {
            console.error("Critical validation failure:", err)
        }

        // --- NEW FILL LOGIC ---
        // Safely retrieve the schema fields for the 'jobTimeslots' array items
        const timeslotsField = (this.schema.fields as any)?.find?.((f: any) => f.localName === 'jobTimeslots');
        const timeslotMessageSchema = timeslotsField?.message;

        if (timeslotMessageSchema) {
            // Safely pull the properties of a single JobTimeslot schema
            const fieldsList = typeof timeslotMessageSchema.fields?.list === 'function'
                ? timeslotMessageSchema.fields.list()
                : (timeslotMessageSchema.fields || []);

            const schemaKeys = fieldsList.map((f: any) => f.localName).filter(Boolean);

            // Get the active number of elements in the form template right now
            const currentValues = this.formValues() as any;
            const totalSlotsCount = Array.isArray(currentValues?.jobTimeslots) ? currentValues.jobTimeslots.length : 0;

            if (totalSlotsCount > 0) {
                // Ensure the errors object has an array allocated for jobTimeslots
                if (!errors.jobTimeslots) {
                    errors.jobTimeslots = [];
                }

                for (let i = 0; i < totalSlotsCount; i++) {
                    // If this slot doesn't exist yet in the error schema, create it
                    if (!errors.jobTimeslots[i]) {
                        errors.jobTimeslots[i] = {};
                    }

                    // Fill missing error keys explicitly with null
                    for (const key of schemaKeys) {
                        if (errors.jobTimeslots[i][key] === undefined) {
                            errors.jobTimeslots[i][key] = null;
                        }
                    }
                }
            }
        }
        // ----------------------

        this.formErrors.set(errors as Partial<DeepStringify<T>>)
    }
}

type DeepStringify<T> = {
  [K in keyof T]: T[K] extends Array<infer U>
    ? Array<DeepStringify<U>> // Map arrays to arrays of stringified errors
    : T[K] extends object
      ? T[K] extends (...args: any[]) => any
        ? T[K]
        : DeepStringify<T[K]>
      : string
}

type NestedPaths<T> = T extends object
  ? T extends Array<infer U>
    ? `${number}` | `${number}.${NestedPaths<U>}`
    : {
        [K in keyof T & string]: K extends `$${string}` // Completely filters out typeName, unknown, etc.
          ? never
          : T[K] extends object
          ? K | `${K}.${NestedPaths<T[K]>}`
          : K
      }[keyof T & string]
  : never

type PathType<T, P extends string> = P extends `${infer Head}.${infer Tail}`
  ? T extends Array<infer U>
    ? PathType<U, Tail>
    : Head extends keyof T
    ? PathType<T[Head], Tail>
    : never
  : T extends Array<infer U>
  ? U
  : P extends keyof T
  ? T[P]
  : never