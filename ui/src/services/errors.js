export const KindNotFound = 'ERR_NOT_FOUND'
export const KindUnexpected = 'ERR_UNEXPECTED'
export const KindConflict = 'ERR_CONFLICT'
export const KindUnauthorized = 'ERR_UNAUTHORIZED'
export const KindBadRequest = 'ERR_BADREQUEST'
export const KindUnsupported = 'ERR_UNSUPPORTED'
export const KindIrrelevant = 'ERR_IRRELEVANT'

export class ApiError extends Error {
  constructor (kind, message) {
    super()
    this.kind = kind
    this.message = message
  }

  toString () {
    return this.message
  }
}
