export const ERR_LOGIN = 'ERR_LOGIN';
export const ERR_CREDENTIALS = 'ERR_CREDENTIALS';
export const ERR_USERNAMETAKEN = 'ERR_USERNAMETAKEN';
export const ERR_EMAILTAKEN = 'ERR_EMAILTAKEN';
export const ERR_RESOURCE_NOTFOUND = 'ERR_RESOURCE_NOTFOUND';

export default class ApiError extends Error {
    kind: string;
    status: number;

    constructor(kind: string, status: number) {
        super(kind);

        this.kind = kind;
        this.status = status;
    }

    isA(kind: string): boolean {
        return this.kind == kind;
    }
}
