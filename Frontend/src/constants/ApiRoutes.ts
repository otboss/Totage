import { environment } from "src/environments/environment"

export namespace ApiRoutes {
    export namespace Backend {
        const LogisticService = (): string => {
            return `${environment.backendDomain}/v1.LogisticService`
        }
        const OrderService = (): string => {
            return `${environment.backendDomain}/v1.OrderService`
        }
        const UserService = (): string => {
            return `${environment.backendDomain}/v1.UserService`
        }
        const OrganizationService = (): string => {
            return `${environment.backendDomain}/v1.OrganizationService`
        }
        const LoginViaPassword = (): string => {
            return `/`
        }

    }
}