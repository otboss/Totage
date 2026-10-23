export namespace AppRoutes {
    
    export const Home = (): string => {
        return ""
    }
    
    export const Login = (): string => {
        return "login"
    }

    // export const Portal = (): string => {
    //     return "portal"
    // }

    export const Dashboard = (): string => {
        return `dashboard`
    }

    export const Organization = (): string => {
        return `${Dashboard()}/organization/:organizationId`
    }

    // export const OrganizationAccount = (): string => {
    //     return `${Dashboard()}/organization/:id/accounts/:accountId`
    // }

    export const OrganizationCreateAccount = (): string => {
        return `${Dashboard()}/organization/:organizationId/accounts/create`
    }

    export const OrganizationEditAccount = (): string => {
        return `${Dashboard()}/organization/:organizationId/accounts/:accountId`
    }

    export const OrderType = (): string => {
        return `${Dashboard()}/organization/:organizationId/order-type`
    }

    export const OrganizationSettings = (): string => {
        return `${Dashboard()}/organization/:id/settings`
    }

    export const OrganizationsCreate = (): string => {
        return `${Dashboard()}/organizations/create`
    }

        export const OrganizationAddEditOrderType = (): string => {
        return `${Dashboard()}/organization/:organizationId/order-type`
    }

    export const Orders = (): string => {
        return `${Home()}/orders`
    }

    export const OrderDetails = (orderId?: string): string => {
        return `${Home()}/orders/${orderId?.toString() ?? `:id`}`
    }

    export const Vehicles = (): string => {
        return `${Organization()}/vehicles`
    }

    export const VehiclesCreate = (): string => {
        return `${Organization()}/vehicles/create`
    }

    export const VehiclesEdit = (vehicleId?: BigInt): string => {
        return `${Organization()}/vehicles/${vehicleId?.toString() ?? ':vehicleId'}`
    }

    export const Register = (): string => {
        return `register`
    }

    export const TermsOfService = (): string => {
        return `terms-of-service`
    }

    export const PrivacyPolicy = (): string => {
        return `privacy-policy`
    }
    
    export const LocationSelector = (): string => {
        return "shared/location-selector"
    }
}